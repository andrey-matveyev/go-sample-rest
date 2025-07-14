package middleware

import (
	"log/slog"
	"net/http"
)

func CorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)

			rLog(r).Info(
				"Preflight request received",
				slog.String("from_origin", r.Header.Get("Origin")))
			return
		}
		next.ServeHTTP(w, r)
	})
}
