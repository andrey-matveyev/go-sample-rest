package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// RequestBody определяет структуру ожидаемого JSON-тела запроса
type RequestBody struct {
	Matrix [][]int `json:"matrix"`
	Value  int     `json:"value"`
}

func tableDataHandler(w http.ResponseWriter, r *http.Request) {
	// Проверяем, что это POST-запрос
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не разрешен. Ожидается POST.", http.StatusMethodNotAllowed)
		return
	}

	// Декодируем JSON-тело запроса
	var requestData RequestBody
	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		http.Error(w, fmt.Sprintf("Некорректный JSON-формат запроса: %v", err), http.StatusBadRequest)
		return
	}

	// Дополнительная валидация (опционально):
	// Проверяем, что матрица имеет размер 3x3
	if len(requestData.Matrix) != 3 {
		http.Error(w, "Матрица должна быть 3x3", http.StatusBadRequest)
		return
	}
	for _, row := range requestData.Matrix {
		if len(row) != 3 {
			http.Error(w, "Каждая строка матрицы должна содержать 3 элемента", http.StatusBadRequest)
			return
		}
	}

	log.Printf("Получена матрица: %+v", requestData.Matrix)
	log.Printf("Получен параметр (целое число): %d", requestData.Value)

	// Здесь можно обрабатывать полученные данные
	// Например, сохранить их в базу данных, выполнить какую-то логику и т.д.

	// Отправляем ответ клиенту
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Данные успешно получены на сервере!"})
}

func main() {
	// Устанавливаем CORS-заголовки
	corsHandler := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")              // Разрешить запросы с любого домена
			w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS") // Разрешаем POST и OPTIONS
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")  // Разрешаем заголовок Content-Type
			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusOK)
				return
			}
			next.ServeHTTP(w, r)
		})
	}

	http.Handle("/api/table-data", corsHandler(http.HandlerFunc(tableDataHandler)))

	log.Println("Сервер запущен на :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
