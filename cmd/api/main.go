package main

import (
	"log"
	"main/internal/config"
	"main/internal/handlers"
	"main/internal/logger"
	"main/internal/middleware"
	"net/http"

	"github.com/go-chi/chi/v5" // Импортируем Chi роутер
)

const configFile string = "./cfg/config.yaml"
const requestIDprefix string = "req"

func main() {
	cfg := config.ReadConfig(configFile)
	logger.SetupDefaultLogger(cfg)

	r := chi.NewRouter() // Создаем новый Chi роутер

	// Middleware для всех маршрутов
	r.Use(middleware.RequestIDMiddleware(requestIDprefix)) //
	r.Use(middleware.LoggingMiddleware)                    // Логирование запросов
	r.Use(middleware.CorsMiddleware)                       // Наш CORS middleware
	r.Route("/api/v1", func(r chi.Router) {
		r.Route("/new-board/{player}", func(r chi.Router) {
			r.Use(middleware.PlayerValidationMiddleware)
			r.Post("/", handlers.NewBoardHandler)
			r.Options("/", handlers.NewBoardHandler)
		})
		r.Route("/make-move/{player}", func(r chi.Router) {
			r.Use(middleware.PlayerValidationMiddleware)
			r.Use(middleware.BoardValidationMiddleware)
			r.Post("/", handlers.MakeMoveHandler)
			r.Options("/", handlers.MakeMoveHandler)
		})
	})

	// Запуск сервера (используем HTTP для простоты локальной разработки,
	// потом можно легко переключить на HTTPS, как обсуждали ранее)
	log.Println("Сервер запущен на HTTP по порту :8080")
	log.Fatal(http.ListenAndServe(":8080", r)) // Передаем роутер 'r'
	// Если хотите HTTPS, замените на:
	// log.Fatal(http.ListenAndServeTLS(":8443", certFile, keyFile, r))
}
