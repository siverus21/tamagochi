package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"tamagochi-team/siverus/tamagochi/internal/config"
	"tamagochi-team/siverus/tamagochi/internal/config/lib/logger/sl"
	"tamagochi-team/siverus/tamagochi/internal/storage/postgres"

	mwLogger "tamagochi-team/siverus/tamagochi/internal/http-server/middleware/logger"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	// TODO: config -> cleanenv
	cfg := config.MustLoad()

	// TODO: init logger -> slog
    log := setupLogger(cfg.Env)

    log.Info("start app", slog.String("env", cfg.Env))
    
	// TODO: init storage -> postgresql

	// Формируем строку подключения для PostgreSQL
	connStr := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=%s",
		cfg.Postgres.User, cfg.Postgres.Password, cfg.Postgres.DBName, cfg.Postgres.SSLMode)

    _, err := postgres.New(connStr)
    if err != nil {
        log.Error("filed init storage", sl.Err(err))
        os.Exit(1)
    }

	// TODO: init router -> chi

	router := chi.NewRouter()

	// middleware
	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(mwLogger.New(log))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	// TODO: run server
	log.Info("strating server", slog.String("address", cfg.Address))

	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Error("failed to start server")
	}

	log.Error("server stopped")
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(
            slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
        )
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	default: // If env config is invalid, set prod settings by default due to security
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}