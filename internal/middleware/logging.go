package middleware

import (
	"log/slog"
	"main/internal/logger"
	"net/http"

	chi_mw "github.com/go-chi/chi/v5/middleware"
)

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctxLogger := slog.Default().With(
			slog.String("Method", r.Method),
			slog.String("RequestURI", r.RequestURI),
			slog.String("RequestID", chi_mw.GetReqID(r.Context())), // depends on chi-middleware
		)
		//ctx := context.WithValue(ctx, contextKey, logger)
		ctx := r.Context()
		ctx = logger.ContextWithLogger(ctx, ctxLogger)

		//next.ServeHTTP(w, r.WithContext(ctx)) // logger in context
		next.ServeHTTP(w, r.WithContext(ctx)) // logger in context
	})
}
