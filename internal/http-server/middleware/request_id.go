package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

type ctxRequestIDKeyType string

const ctxRequestIDKey ctxRequestIDKeyType = "ctxRequestIDKey"

func RequestIDMiddleware(prefix string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			newUUID := uuid.New().String()
			shortUUID := strings.ReplaceAll(newUUID, "-", "")
			requestID := fmt.Sprintf("%s-%s", prefix, shortUUID)

			ctx := context.WithValue(r.Context(), ctxRequestIDKey, requestID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetRequestIDFromContext(ctx context.Context) (string, bool) {
	id, ok := ctx.Value(ctxRequestIDKey).(string)
	return id, ok
}
