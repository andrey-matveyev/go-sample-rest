package middleware

import (
	"log/slog"
	"main/internal/logger"
	"net/http"
)

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctxLogger := slog.Default().With(
			slog.String("RequestURI", r.RequestURI),
			slog.String("Method", r.Method),
			slog.String("URL_Path", r.URL.Path),
		)
		ctx := logger.ContextWithLogger(r.Context(), ctxLogger)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
