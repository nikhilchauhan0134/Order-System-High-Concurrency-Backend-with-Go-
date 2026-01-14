package db

import (
	"database/sql"
	"fmt"

	"OrderSystemHighConcurrency/order-processor/internal/config"

	_ "github.com/denisenkom/go-mssqldb"
)

func NewDB(cfg *config.Config) (*sql.DB, error) {
	db, err := sql.Open("sqlserver", cfg.DBDSN)
	if err != nil {
		return nil, fmt.Errorf("sql open error: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("db ping error: %w", err)
	}

	return db, nil
}
