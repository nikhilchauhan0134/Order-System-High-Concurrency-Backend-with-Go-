package db

import (
	"database/sql"
	"time"
)

// NewDB creates a database connection with pooling
func NewDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("sqlserver", dsn)
	if err != nil {
		return nil, err
	}

	// Connection pool settings (VERY IMPORTANT)
	db.SetMaxOpenConns(50)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(30 * time.Minute)

	// Validate connection
	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
