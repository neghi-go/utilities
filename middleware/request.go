package middleware

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/neghi-go/utilities"
)

type statusResponseWriter struct {
	http.ResponseWriter // Embed a http.ResponseWriter
	statusCode          int
	headerWritten       bool
}

func newstatusResponseWriter(w http.ResponseWriter) *statusResponseWriter {
	return &statusResponseWriter{
		ResponseWriter: w,
		statusCode:     http.StatusOK,
	}
}

func (mw *statusResponseWriter) WriteHeader(statusCode int) {
	mw.ResponseWriter.WriteHeader(statusCode)

	if !mw.headerWritten {
		mw.statusCode = statusCode
		mw.headerWritten = true
	}
}

func (mw *statusResponseWriter) Write(b []byte) (int, error) {
	mw.headerWritten = true
	return mw.ResponseWriter.Write(b)
}

func (mw *statusResponseWriter) Unwrap() http.ResponseWriter {
	return mw.ResponseWriter
}

func Request(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		now := time.Now()
		req_id := utilities.Generate(12)

		req := r.WithContext(context.WithValue(r.Context(), "request_id", req_id))
		res := newstatusResponseWriter(w)

		res.Header().Add("X-Request-ID", req_id)
		res.Header().Set("Cache-Control", "no-store")
		next.ServeHTTP(res, req)
		slog.Info(
			r.Method,
			slog.String("request_path", r.URL.Path),
			slog.Int64("statusCode", int64(res.statusCode)),
			slog.String("request_id", req_id),
			slog.Duration("tte", time.Since(now)),
			slog.String("user-agent", r.UserAgent()),
		)
	})
}
