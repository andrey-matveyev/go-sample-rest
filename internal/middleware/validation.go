package middleware

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func PlayerValidationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//const op = "ValidationMiddleware"
		//log := logger.CtxLogger(r.Context())

		playerStr := chi.URLParam(r, "player")
		//log.Info("ValidationMiddleware", slog.String("playerStr", playerStr))
		fmt.Printf("playerStr from ValidationMiddleware: %s\n", playerStr)
		fmt.Printf("request_ptr: %p\n", r)
		fmt.Printf("context_ptr: %p\n", r.Context)

		next.ServeHTTP(w, r)
	})
}

func BoardValidationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		/*const op = "ValidationMiddleware"
		log := logger.CtxLogger(r.Context())

		// player
		playerStr := chi.URLParam(r, "player")
		log.Info("ValidationMiddleware", slog.String("playerStr", playerStr))
		/*player, err := strconv.Atoi(playerStr)
		if err != nil || (player != 1 && player != -1) {
			log.Warn("Invalid value 'player'. Expected 1 or -1.",
				slog.String("playerStr", playerStr),
				slog.Int("player", player),
				slog.String("op", op))
			//http.Error(w, "Invalid value 'player'.", http.StatusBadRequest)
			//return
		}

		ctx := r.Context()*/
		next.ServeHTTP(w, r)
	})
}
