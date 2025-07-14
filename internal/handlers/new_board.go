package handlers

import (
	"encoding/json"
	"log/slog"
	"main/internal/middleware"
	"main/internal/models"
	"math/rand"
	"net/http"
)

func NewBoardHandler(w http.ResponseWriter, r *http.Request) {
	log := middleware.GetLoggerFromContext(r.Context())

	player, ok := middleware.GetPlayerFromContext(r.Context())
	if !ok {
		log.Error("Missing 'player' value in context. Expected included.")

		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	newBoard := make([][]int, 3)
	for i := range newBoard {
		newBoard[i] = make([]int, 3)
	}

	if player == -1 {
		newBoard[1][1] = 1
	}
	log.Info("New board was generated.", slog.Int("player", player))

	response := models.NewBoardResponse{
		Board:    newBoard,
		IdGame:   rand.Intn(1000), // Просто случайный ID игры для примера
		GameOver: false,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
