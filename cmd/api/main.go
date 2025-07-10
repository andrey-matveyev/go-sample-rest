package main

import (
	"encoding/json"
	"fmt"
	"log"
	"main/internal/config"
	"math/rand" // Для rand.Seed
	"net/http"
	"strconv" // Добавляем импорт для конвертации строки в число

	"github.com/go-chi/chi/v5"            // Импортируем Chi роутер
	"github.com/go-chi/chi/v5/middleware" // Опционально: встроенный middleware Chi
)

// --- ОБНОВЛЕННЫЕ СТРУКТУРЫ ДАННЫХ В СООТВЕТСТВИИ С OPENAPI ---
// Для запроса нового поля /api/v1/new-board/{player}
type NewBoardResponse struct {
	Board    [][]int `json:"board"`
	IdGame   int     `json:"idGame"`
	GameOver bool    `json:"gameOver"`
}

// Для запроса хода /api/v1/make-move/{player}
type MakeMoveRequest struct {
	Board  [][]int `json:"board"`
	IdGame int     `json:"idGame"`
}

// Для ответа на ход /api/v1/make-move/{player}
type MakeMoveResponse struct {
	Board     [][]int `json:"board"`
	IdGame    int     `json:"idGame"`
	GameOver  bool    `json:"gameOver"`
	WinPlayer int     `json:"winPlayer"` // 0 - ничья, 1 или -1 - победитель
}

// --- КОНЕЦ ОБНОВЛЕННЫХ СТРУКТУР ---

// --- ВАША БИЗНЕС-ЛОГИКА (PLACEHOLDER) ---
// Это ваша функция бизнес-логики, которую вы затем замените реальной.
// Пока она будет просто возвращать тестовые данные.
// board: текущая доска
// player: игрок, который только что сделал ход (1 или -1)
// Возвращает:
//
//	nextBoard: доска после хода компьютера/ответа сервера
//	gameOverStatus: 0 - игра продолжается, 1 - победа, -1 - ничья
//	winningPlayer: 0 - никто не победил/ничья, 1 или -1 - победивший игрок
func processGameMove(currentBoard [][]int, playerMakingMove int) (nextBoard [][]int, gameOverStatus int, winningPlayer int) {
	// --- ЗАМЕНИТЕ ЭТО ВАШЕЙ РЕАЛЬНОЙ ЛОГИКОЙ ---
	log.Printf("Бизнес-логика: Обрабатывается ход игрока %d. Текущая доска: %+v", playerMakingMove, currentBoard)

	// Пример: Просто делаем "ход" в случайную пустую ячейку (для тестирования)
	nextBoard = make([][]int, 3)
	for i := range currentBoard {
		nextBoard[i] = make([]int, 3)
		copy(nextBoard[i], currentBoard[i])
	}

	// Найдем случайную пустую ячейку для тестового хода
	emptyCells := [][]int{}
	for r := 0; r < 3; r++ {
		for c := 0; c < 3; c++ {
			if nextBoard[r][c] == 0 {
				emptyCells = append(emptyCells, []int{r, c})
			}
		}
	}

	if len(emptyCells) > 0 {
		randomIndex := rand.Intn(len(emptyCells))
		randomCell := emptyCells[randomIndex]
		// Отвечает игрок, противоположный тому, кто прислал запрос
		opponentPlayer := -playerMakingMove
		nextBoard[randomCell[0]][randomCell[1]] = opponentPlayer
	}

	// Пример: Игра продолжается, пока не заполнены все ячейки
	isBoardFull := true
	for r := 0; r < 3; r++ {
		for c := 0; c < 3; c++ {
			if nextBoard[r][c] == 0 {
				isBoardFull = false
				break
			}
		}
		if !isBoardFull {
			break
		}
	}

	if isBoardFull {
		log.Println("Бизнес-логика: Доска заполнена, ничья.")
		return nextBoard, 0, 0 // Ничья
	}

	log.Println("Бизнес-логика: Игра продолжается.")
	return nextBoard, -1, 0 // Игра продолжается
	// --- КОНЕЦ PLACEHOLDER ЛОГИКИ ---
}

// --- КОНЕЦ ВАШЕЙ БИЗНЕС-ЛОГИКИ ---

// corsMiddleware - применяем его ко всем маршрутам через r.Use()
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// --- НОВЫЙ ХЭНДЛЕР: POST /api/v1/new-board/{player} ---
func newBoardHandler(w http.ResponseWriter, r *http.Request) {
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

	response := NewBoardResponse{
		Board:    newBoard,
		IdGame:   rand.Intn(1000), // Просто случайный ID игры для примера
		GameOver: false,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// --- НОВЫЙ ХЭНДЛЕР: POST /api/v1/make-move/{player} ---
func makeMoveHandler(w http.ResponseWriter, r *http.Request) {
	playerStr := chi.URLParam(r, "player")
	player, err := strconv.Atoi(playerStr)
	if err != nil || (player != 1 && player != -1) {
		http.Error(w, "Некорректное значение 'player' в URL-пути. Ожидается 1 или -1.", http.StatusBadRequest)
		return
	}

	var requestBody MakeMoveRequest
	err = json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		http.Error(w, fmt.Sprintf("Некорректный JSON-формат запроса: %v", err), http.StatusBadRequest)
		return
	}

	log.Printf("Получен ход от игрока %d. IdGame: %d, Доска: %+v", player, requestBody.IdGame, requestBody.Board)

	// --- ВЫЗЫВАЕМ ВАШУ БИЗНЕС-ЛОГИКУ ---
	// Здесь вы подставите вызов вашей функции, которая обработает ход.
	// Она должна вернуть новую доску, статус игры и победителя.
	nextBoard, gameOverStatus, winningPlayer := processGameMove(requestBody.Board, player)
	// --- КОНЕЦ ВЫЗОВА БИЗНЕС-ЛОГИКИ ---

	gameOver := false
	if gameOverStatus != -1 { // -1 означает, что игра продолжается
		gameOver = true
	}

	response := MakeMoveResponse{
		Board:     nextBoard,
		IdGame:    requestBody.IdGame, // ID игры остается тем же
		GameOver:  gameOver,
		WinPlayer: winningPlayer,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

const configFile string = "./cfg/config.yaml"

func main() {
	cfg := config.ReadConfig(configFile)
	log.Println("Config: ", cfg)

	r := chi.NewRouter() // Создаем новый Chi роутер

	// Middleware для всех маршрутов
	r.Use(middleware.Logger) // Логирование запросов
	r.Use(corsMiddleware)    // Наш CORS middleware

	// --- РЕГИСТРАЦИЯ НОВЫХ МАРШРУТОВ И ХЭНДЛЕРОВ ---
	r.Post("/api/v1/new-board/{player}", newBoardHandler)
	r.Options("/api/v1/new-board/{player}", newBoardHandler) // Для CORS preflight

	r.Post("/api/v1/make-move/{player}", makeMoveHandler)
	r.Options("/api/v1/make-move/{player}", makeMoveHandler) // Для CORS preflight
	// --- КОНЕЦ РЕГИСТРАЦИИ МАРШРУТОВ ---

	// Запуск сервера (используем HTTP для простоты локальной разработки,
	// потом можно легко переключить на HTTPS, как обсуждали ранее)
	log.Println("Сервер запущен на HTTP по порту :8080")
	log.Fatal(http.ListenAndServe(":8080", r)) // Передаем роутер 'r'
	// Если хотите HTTPS, замените на:
	// log.Fatal(http.ListenAndServeTLS(":8443", certFile, keyFile, r))
}
