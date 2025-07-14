package models

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
