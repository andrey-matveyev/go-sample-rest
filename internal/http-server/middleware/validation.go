package middleware

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"main/internal/http-server/models"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type ctxPlayerKeyType string

const ctxPlayerKey ctxPlayerKeyType = "ctxPlayerKey"

func PlayerValidationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		playerStr := chi.URLParam(r, "player")
		player, err := strconv.Atoi(playerStr)

		if err != nil || (player != 1 && player != -1) {
			rLog(r).Error("invalid value 'player'. Expected 1 or -1.",
				slog.String("playerStr", playerStr))

			http.Error(w, "Invalid value 'player'.", http.StatusBadRequest)
			return
		}

		ctx := context.WithValue(r.Context(), ctxPlayerKey, player)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetPlayerFromContext(ctx context.Context) (int, bool) {
	item, ok := ctx.Value(ctxPlayerKey).(int)
	return item, ok
}

type ctxBoardKeyType string

const ctxBoardKey ctxBoardKeyType = "ctxBoardKey"

func BoardValidationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		buf, err := io.ReadAll(r.Body)
		if err != nil {
			rLog(r).Error("error reading request body.",
				"error", slog.StringValue(err.Error()))

			http.Error(w, "Invalid body in request.", http.StatusBadRequest)
			return
		}

		var requestBody models.MakeMoveRequest
		err = json.Unmarshal(buf, &requestBody)
		if err != nil {
			rLog(r).Error("error parsing request body as JSON.",
				slog.String("error", err.Error()),
				slog.String("body", string(buf)),
			)
			http.Error(w, "Invalid body in request.", http.StatusBadRequest)
			return
		}

		if len(requestBody.Board) != 3 {
			rLog(r).Error("error count of row in the 'board'. Expected 3",
				slog.Int("rows", len(requestBody.Board)),
			)
			http.Error(w, "Invalid body in request.", http.StatusBadRequest)
			return
		}

		for i := range 2 {
			if len(requestBody.Board[i]) != 3 {
				rLog(r).Error("error count of column in the 'board'. Expected 3",
					slog.Int("columns", len(requestBody.Board[i])),
				)
				http.Error(w, "Invalid body in request.", http.StatusBadRequest)
				return
			}
		}

		for i := range 2 {
			for j := range 2 {
				if requestBody.Board[i][j] != 1 && requestBody.Board[i][j] != 0 && requestBody.Board[i][j] != -1 {
					rLog(r).Error("invalid value in the 'board'. Expected 1, 0 or -1.",
						slog.Int("value", requestBody.Board[i][j]))

					http.Error(w, "Invalid body in request.", http.StatusBadRequest)
					return
				}
			}
		}

		ctx := context.WithValue(r.Context(), ctxBoardKey, requestBody)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetBoardFromContext(ctx context.Context) (models.MakeMoveRequest, bool) {
	item, ok := ctx.Value(ctxBoardKey).(models.MakeMoveRequest)
	return item, ok
}
