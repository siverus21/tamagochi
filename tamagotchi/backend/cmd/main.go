package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"tamagochi-team/siverus/tamagochi/internal/config"
	"tamagochi-team/siverus/tamagochi/internal/lib/logger/sl"
	"tamagochi-team/siverus/tamagochi/internal/storage/postgres"

	handlers "tamagochi-team/siverus/tamagochi/internal/http-server/handlers/auth"
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
	// Загружаем конфигурацию
	cfg := config.MustLoad()

	// Настраиваем логгер
	log := setupLogger(cfg.Env)
	log.Info("start app", slog.String("env", cfg.Env))

	// Формируем строку подключения для PostgreSQL
	connStr := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=%s",
		cfg.Postgres.User, cfg.Postgres.Password, cfg.Postgres.DBName, cfg.Postgres.SSLMode)

	// Инициализируем хранилище и сохраняем экземпляр в переменную
	storageInstance, err := postgres.New(connStr)
	if err != nil {
		log.Error("failed init storage", sl.Err(err))
		os.Exit(1)
	}

	// Инициализируем маршрутизатор chi
	router := chi.NewRouter()

	// Подключаем middleware
	router.Use(middleware.Logger)
	router.Use(mwLogger.New(log))
	router.Use(middleware.Recoverer)

	// Регистрируем публичный маршрут /login, передавая storageInstance в обработчик
	router.Post("/login", func(w http.ResponseWriter, r *http.Request) {
		handlers.LoginHandler(storageInstance, w, r)
	})

	// Запускаем сервер
	log.Info("starting server", slog.String("address", cfg.Address))
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
	default: // Если значение env неверное, используем настройки для prod в целях безопасности
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}
