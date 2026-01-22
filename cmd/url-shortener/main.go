package main

import (
	"log/slog"
	"os"

	"url-shortener/internal"
	"url-shortener/internal/lib/sl"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := internal.MustLoad()

	log := setupLogger(cfg.Env)
	log.Info("Log info", slog.String("env", cfg.Env))

	storage, err := internal.SqliteNew(cfg.StoragePath)
	if err != nil {
		log.Error("Failed to init storage", sl.AttrByErr(err))
		os.Exit(1)
	}

	_ = storage
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envDev:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return log
}
