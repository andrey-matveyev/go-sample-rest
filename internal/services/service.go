package services

import (
	"log"
	"math/rand"
)

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
func ProcessGameMove(currentBoard [][]int, playerMakingMove int) (nextBoard [][]int, gameOverStatus int, winningPlayer int) {
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
