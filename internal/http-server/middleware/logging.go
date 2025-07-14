package middleware

import (
	"context"
	"log/slog"
	"net/http"
)

type ctxLoggerKeyType string

const ctxLoggerKey ctxLoggerKeyType = "ctxLoggerKey"

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID, ok := GetRequestIDFromContext(r.Context())
		if !ok {
			requestID = "no-request-id" // default value
		}
		ctxLogger := slog.Default().With(
			slog.String("Method", r.Method),
			slog.String("RequestURI", r.RequestURI),
			slog.String("RequestID", requestID),
		)

		ctx := context.WithValue(r.Context(), ctxLoggerKey, ctxLogger)
		next.ServeHTTP(w, r.WithContext(ctx)) // logger in context
	})
}

// Can be used (if needed) in a current package in a middleware call chain
// Instead of GetLoggerFromContext
func rLog(r *http.Request) *slog.Logger {
	return GetLoggerFromContext(r.Context())
}

// Returns logger from context.
func GetLoggerFromContext(ctx context.Context) *slog.Logger {
	if logger, ok := ctx.Value(ctxLoggerKey).(*slog.Logger); ok {
		return logger
	}

	return slog.Default()
}
