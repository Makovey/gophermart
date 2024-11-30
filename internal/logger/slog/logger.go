package slog

import (
	"log/slog"
	"os"

	"github.com/Makovey/gophermart/internal/logger"
)

type EnvString string

const (
	Prod  EnvString = "prod"
	Dev   EnvString = "dev"
	Local EnvString = "local"
)

func NewLogger(env EnvString) logger.Logger {
	var log *slog.Logger

	switch env {
	case Prod:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	case Dev:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	default:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	}

	return log
}
