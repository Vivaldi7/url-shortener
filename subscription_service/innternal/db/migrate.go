package db

import (
	"database/sql"
	"fmt"
	"path/filepath"

	"github.com/pressly/goose/v3"
	"github.com/vivaldi7/golang_code/subscription_service/logger"
)

func RunMigrations(db *sql.DB) error {
	// filepath.Join автоматически использует правильный разделитель для ОС
	// На Windows: migrations\
	// На Linux:   migrations/
	// Указываем директорию с миграциями
	migrationsDir := filepath.Join(".", "innternal", "db", "migrations")

	// Проверяем статус
	if err := goose.Status(db, migrationsDir); err != nil {
		logger.Log.Warn("Failed to get migration status", "error", err)
		return err
	}

	// Применяем все ожидающие миграции
	if err := goose.Up(db, migrationsDir); err != nil {
		return fmt.Errorf("failed to apply migrations: %w", err)
	}

	logger.Log.Info("Migrations applied successfully")
	return nil
}
