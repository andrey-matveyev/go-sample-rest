package handlers

import (
	"encoding/json"
	"main/internal/http-server/middleware"
	"main/internal/http-server/models"
	"main/internal/http-server/repository"
	"main/internal/services"
	"net/http"
)

func MakeMoveHandler(storage repository.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := middleware.GetLoggerFromContext(r.Context())

		player, ok := middleware.GetPlayerFromContext(r.Context())
		if !ok {
			log.Error("missing 'player' value in context. Expected included.")

			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		body, ok := middleware.GetBodyFromContext(r.Context())
		if !ok {
			log.Error("missing 'board' in context. Expected included.")

			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		jsonBoard, err := json.Marshal(body.Board)
		if err != nil {
			log.Error("error marshal board to json format")

			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		err = storage.SaveNewMove(body.IdGame, string(jsonBoard))
		if err != nil {
			log.Error("error save new move to storage")

			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		//log.Printf("Получен ход от игрока %d. IdGame: %d, Доска: %+v", player, requestBody.IdGame, requestBody.Board)

		// --- ВЫЗЫВАЕМ ВАШУ БИЗНЕС-ЛОГИКУ ---
		// Здесь вы подставите вызов вашей функции, которая обработает ход.
		// Она должна вернуть новую доску, статус игры и победителя.
		nextBoard, gameOverStatus, winningPlayer := services.ProcessGameMove(body.Board, player)
		// --- КОНЕЦ ВЫЗОВА БИЗНЕС-ЛОГИКИ ---
		jsonBoard, err = json.Marshal(nextBoard)
		if err != nil {
			log.Error("error marshal board to json format")

			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		err = storage.SaveNewMove(body.IdGame, string(jsonBoard))
		if err != nil {
			log.Error("error save new move to storage")

			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		gameOver := false
		if gameOverStatus != -1 { // -1 означает, что игра продолжается
			gameOver = true
		}

		response := models.MakeMoveResponse{
			Board:     nextBoard,
			IdGame:    body.IdGame, // ID игры остается тем же
			GameOver:  gameOver,
			WinPlayer: winningPlayer,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}
