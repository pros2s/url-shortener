package main

import (
	"log/slog"
	"os"

	"url-shortener/internal/config"
	"url-shortener/internal/lib/sl"
	"url-shortener/internal/storage/sqlite"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)
	log.Info("Log info", slog.String("env", cfg.Env))

	// create storage
	storage, err := sqlite.SqliteNew(cfg.StoragePath)
	if err != nil {
		log.Error("Failed to init storage", sl.AttrByErr(err))
		os.Exit(1)
	}

	// save
	id, err := storage.SaveToUrl("https://google.com", "google")
	if err != nil {
		log.Error("Failed to save url", sl.AttrByErr(err))
		os.Exit(1)
	}

	log.Info("Save url with id: ", slog.Attr{Key: "id", Value: slog.Int64Value(id)})

	// remove
	if err := storage.RemoveUrl(id); err != nil {
		log.Error("Failed to delete url", sl.AttrByErr(err))
		os.Exit(1)
	}

	log.Info("Delete url with id: ", slog.Attr{Key: "id", Value: slog.Int64Value(id)})
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
