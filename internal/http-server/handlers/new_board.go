package handlers

import (
	"encoding/json"
	"log/slog"
	"main/internal/http-server/middleware"
	"main/internal/http-server/models"
	"main/internal/http-server/repository"
	"net/http"
)

func NewBoardHandler(storage repository.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := middleware.GetLoggerFromContext(r.Context())

		player, ok := middleware.GetPlayerFromContext(r.Context())
		if !ok {
			log.Error("missing 'player' value in context, expected included")

			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		idGame, err := storage.SaveNewGame(player)
		if err != nil {
			log.Error("error save new game to storage and get ID of game")

			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		newBoard := make([][]int, 3)
		for i := range newBoard {
			newBoard[i] = make([]int, 3)
		}
		if player == -1 {
			newBoard[1][1] = 1 // Ход нейросети

			jsonBoard, err := json.Marshal(newBoard)
			if err != nil {
				log.Error("error marshal new board to json format")

				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			err = storage.SaveNewMove(idGame, string(jsonBoard))
			if err != nil {
				log.Error("error save new move to storage")

				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
		}

		log.Info("New board was generated.", slog.Int("player", player), slog.Int64("idGame", idGame))

		response := models.NewBoardResponse{
			Board:    newBoard,
			IdGame:   idGame,
			GameOver: false,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}
