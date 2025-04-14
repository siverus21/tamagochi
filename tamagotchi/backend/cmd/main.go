package main

import (
	"fmt"
	"log/slog"
	"os"
	"tamagochi-team/siverus/tamagochi/internal/config"
	"tamagochi-team/siverus/tamagochi/internal/config/lib/logger/sl"
	"tamagochi-team/siverus/tamagochi/internal/storage/postgres"
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

    storage, err := postgres.New(connStr)
    if err != nil {
        log.Error("filed init storage", sl.Err(err))
        os.Exit(1)
    }

    users, err := storage.GetAllUsers()
    if err != nil {
        fmt.Println("Ошибка:", err)
        return
    }

    for _, u := range users {
        fmt.Printf(u.UserName)
    }


	// TODO: init router -> chi

	// TODO: run server
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