package internal

import (
	"database/sql"
	"fmt"

	_ "github.com/glebarez/go-sqlite" // get sqlite driver
)

type SqliteStorage struct {
	db *sql.DB
}

func SqliteNew(storagePath string) (*SqliteStorage, error) {
	fn := "internal.sqlite.New"

	db, err := sql.Open("sqlite", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s, %w", fn, err)
	}

	stmt, err := db.Prepare(`
		CREATE TABLE IF NOT EXISTS url(
			id INTEGER PRIMARY KEY,
			alias TEXT NOT NULL UNIQUE,
			url TEXT NOT NULL
		);
		CREATE INDEX IF NOT EXISTS idx_alias ON url(alias);
	`)
	if err != nil {
		return nil, fmt.Errorf("%s, %w", fn, err)
	}

	if _, err := stmt.Exec(); err != nil {
		return nil, fmt.Errorf("%s, %w", fn, err)
	}

	return &SqliteStorage{db: db}, nil
}
