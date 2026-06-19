package logger

import (
	"log/slog"
	"os"
)

// Создаем глобальную переменню
var Log *slog.Logger

func InitLogger(level string) {
	var logLevel slog.Level
	switch level {
	case "debug":
		logLevel = slog.LevelDebug
	case "info":
		logLevel = slog.LevelInfo
	case "warn":
		logLevel = slog.LevelWarn
	case "error":
		logLevel = slog.LevelError
	default:
		logLevel = slog.LevelInfo
	}

	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{ //здесь обработчик формирует логи в json и выводит в стандартный вывод
		Level: logLevel,
	})

	Log = slog.New(handler) //создаем экземпляр логера
	slog.SetDefault(Log)    //делаем его дефолтным
}
