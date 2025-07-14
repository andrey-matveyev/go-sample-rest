package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"main/internal/http-server/models"
	"main/internal/services"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func MakeMoveHandler(w http.ResponseWriter, r *http.Request) {
	playerStr := chi.URLParam(r, "player")
	player, err := strconv.Atoi(playerStr)
	if err != nil || (player != 1 && player != -1) {
		http.Error(w, "Некорректное значение 'player' в URL-пути. Ожидается 1 или -1.", http.StatusBadRequest)
		return
	}

	var requestBody models.MakeMoveRequest
	err = json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		http.Error(w, fmt.Sprintf("Некорректный JSON-формат запроса: %v", err), http.StatusBadRequest)
		return
	}

	log.Printf("Получен ход от игрока %d. IdGame: %d, Доска: %+v", player, requestBody.IdGame, requestBody.Board)

	// --- ВЫЗЫВАЕМ ВАШУ БИЗНЕС-ЛОГИКУ ---
	// Здесь вы подставите вызов вашей функции, которая обработает ход.
	// Она должна вернуть новую доску, статус игры и победителя.
	nextBoard, gameOverStatus, winningPlayer := services.ProcessGameMove(requestBody.Board, player)
	// --- КОНЕЦ ВЫЗОВА БИЗНЕС-ЛОГИКИ ---

	gameOver := false
	if gameOverStatus != -1 { // -1 означает, что игра продолжается
		gameOver = true
	}

	response := models.MakeMoveResponse{
		Board:     nextBoard,
		IdGame:    requestBody.IdGame, // ID игры остается тем же
		GameOver:  gameOver,
		WinPlayer: winningPlayer,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
