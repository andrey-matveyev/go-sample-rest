package main

import (
	"context"
	"fmt"
	"log/slog"
	"main/internal/config"
	"main/internal/http-server/handlers"
	"main/internal/http-server/middleware"
	"main/internal/http-server/repository"
	"main/internal/logger"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
)

const configFile string = "./cfg/config.yaml"
const requestIDprefix string = "req"

func main() {
	cfg := config.ReadConfig(configFile)
	log := logger.SetupDefaultLogger(cfg)

	storage, err := repository.NewStorage(cfg.RepositoryFile)
	if err != nil {
		log.Error("storage not created.", slog.String("error", err.Error()))
		return
	}
	log.Info("Storage created successfully.")

	defer func() {
		if err := storage.Shutdown(); err != nil {
			log.Error("error shutting down storage.", slog.String("error", err.Error()))
			return
		}
		log.Info("Storage shutdown gracefully.")
	}()

	router := chi.NewRouter()

	router.Use(middleware.RequestIDMiddleware(requestIDprefix))
	router.Use(middleware.LoggingMiddleware)
	router.Use(middleware.CorsMiddleware)
	router.Route("/api/v1", func(r chi.Router) {
		r.Route("/new-board/{player}", func(r chi.Router) {
			r.Use(middleware.PlayerValidationMiddleware)
			r.Post("/", handlers.NewBoardHandler(storage))
			r.Options("/", handlers.NewBoardHandler(storage))
		})
		r.Route("/make-move/{player}", func(r chi.Router) {
			r.Use(middleware.PlayerValidationMiddleware)
			r.Use(middleware.BoardValidationMiddleware)
			r.Post("/", handlers.MakeMoveHandler(storage))
			r.Options("/", handlers.MakeMoveHandler(storage))
		})
	})

	server := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	serverErrors := make(chan error, 1)
	go func() {
		defer close(serverErrors)

		log.Info("Starting http-server...")

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			serverErrors <- fmt.Errorf("http-server startup error: %w", err)
		}
	}()

	select {
	case err := <-serverErrors:
		log.Error("error starting http-server", slog.String("error", err.Error()))
		return
	case <-time.After(3 * time.Second):
		log.Info("Http-server started successfully.", slog.String("address", server.Addr))

		osSignals := make(chan os.Signal, 1)
		defer close(osSignals)

		signal.Notify(osSignals, syscall.SIGINT, syscall.SIGTERM)

		sig := <-osSignals

		log.Info("Received signal.", slog.String("signal", sig.String()))
		log.Info("Http-server shutting down...")

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			log.Error("http-server shutdown error.", slog.String("error", err.Error()))
			return
		}
		log.Info("Http-server stopped gracefully.")
	}
}
