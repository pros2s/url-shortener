package sqlite

import (
	"database/sql"
	"fmt"

	"url-shortener/internal/storage"

	"github.com/glebarez/go-sqlite"
	sqlite3 "modernc.org/sqlite/lib"
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

func (s *SqliteStorage) SaveToUrl(url string, alias string) (int64, error) {
	fn := "internal.sqlite.SaveToUrl"

	stmt, err := s.db.Prepare("INSERT INTO url(url, alias) VALUES(?, ?)")
	if err != nil {
		return 0, fmt.Errorf("%s, %w", fn, err)
	}

	res, err := stmt.Exec(url, alias)
	if err != nil {
		if sqliteErr, ok := err.(*sqlite.Error); ok && sqliteErr.Code() == sqlite3.SQLITE_CONSTRAINT_UNIQUE {
			return 0, fmt.Errorf("%s, %w", fn, storage.ErrUrlExists)
		}

		return 0, fmt.Errorf("%s, %w", fn, err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s: failed to get last insert id: %w", fn, err)
	}

	return id, nil
}
