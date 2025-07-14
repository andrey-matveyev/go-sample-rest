package handlers

import (
	"encoding/json"
	"log"
	"main/internal/models"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func NewBoardHandler(w http.ResponseWriter, r *http.Request) {
	playerStr := chi.URLParam(r, "player") // Извлекаем параметр из URL с помощью Chi
	player, err := strconv.Atoi(playerStr)
	if err != nil || (player != 1 && player != -1) {
		http.Error(w, "Некорректное значение 'player' в URL-пути. Ожидается 1 или -1.", http.StatusBadRequest)
		return
	}

	newBoard := make([][]int, 3)
	for i := range newBoard {
		newBoard[i] = make([]int, 3)
	}

	// Если игрок -1, делаем первый ход (например, в центр)
	if player == -1 {
		newBoard[1][1] = 1 // Например, сервер делает первый ход 1
		log.Printf("Сгенерирована новая доска для игрока %d с первым ходом сервера.", player)
	} else {
		log.Printf("Сгенерирована новая пустая доска для игрока %d.", player)
	}

	response := models.NewBoardResponse{
		Board:    newBoard,
		IdGame:   rand.Intn(1000), // Просто случайный ID игры для примера
		GameOver: false,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
