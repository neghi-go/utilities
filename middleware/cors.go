package middleware

import "net/http"

func Cors(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		w.Header().Add("Access-Control-Allow-Methods", "GET, POST, PATCH, DELETE, HEAD")
		w.Header().Add("Access-Control-Allow-Headers", "Authorization, Content-Type, X-Refresh")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
		}

		h.ServeHTTP(w, r)
	})
}
