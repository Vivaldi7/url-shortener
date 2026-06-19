package db

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/vivaldi7/golang_code/subscription_service/innternal/config"
	"github.com/vivaldi7/golang_code/subscription_service/logger"
)

func Connect(cfg *config.Config) (*sql.DB, error) {
	connStr := fmt.Sprintf("host = %s port = %s user = %s password = %s dbname = %s sslmode = %s",
		cfg.PGHost, cfg.PGPort, cfg.PGUser, cfg.PGPassword, cfg.PGName, cfg.SSLMode)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		logger.Log.Error("Failed to open database connection", "error", err)
		return nil, err
	}

	if err = db.Ping(); err != nil {
		logger.Log.Error("Failed to ping database", "error", err)
		return nil, err
	}

	db.SetMaxOpenConns(25)                 // лимит
	db.SetMaxIdleConns(10)                 // кэш готовых соединений
	db.SetConnMaxLifetime(5 * time.Minute) // время жизни

	logger.Log.Info("Database connected successfully")
	return db, nil
}
