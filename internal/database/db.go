package database

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

func DbConnect(dbPath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("Fail on open to database: %w", err)
	}

	if err := db.Ping(); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("Fail on connect to database: %w", err)
	}

	db.SetMaxOpenConns(1)

	return db, nil
}
