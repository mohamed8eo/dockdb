package logger

import (
	"log/slog"
	"os"

	"github.com/pterm/pterm"
)

var log *slog.Logger

func Init() {
	pterm.DefaultLogger.Level = pterm.LogLevelDebug

	handler := pterm.NewSlogHandler(&pterm.DefaultLogger)
	log = slog.New(handler)
}

func Debug(msg string, args ...any) {
	log.Debug(msg, args...)
}

func Info(msg string, args ...any) {
	log.Info(msg, args...)
}

func Warn(msg string, args ...any) {
	log.Warn(msg, args...)
}

func Error(msg string, args ...any) {
	log.Error(msg, args...)
}

func Fatal(msg string, args ...any) {
	log.Error(msg, args...)
	os.Exit(1)
}
