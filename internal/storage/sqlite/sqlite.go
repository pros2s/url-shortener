package sqlite

import (
	"database/sql"
	"errors"
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
	defer stmt.Close()

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
	defer stmt.Close()

	res, err := stmt.Exec(url, alias)

	var sqliteErr *sqlite.Error
	if errors.As(err, &sqliteErr); sqliteErr.Code() == sqlite3.SQLITE_CONSTRAINT_UNIQUE {
		return 0, fmt.Errorf("%s, %w", fn, storage.ErrUrlExists)
	}

	if err != nil {
		return 0, fmt.Errorf("%s, %w", fn, err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s: failed to get last insert id: %w", fn, err)
	}

	return id, nil
}

func (s *SqliteStorage) RemoveUrl(id int64) error {
	fn := "internal.sqlite.RemoveUrl"

	stmt, err := s.db.Prepare("DELETE FROM url WHERE id = ?")
	if err != nil {
		return fmt.Errorf("%s, %w", fn, err)
	}
	defer stmt.Close()

	res, err := stmt.Exec(id)
	if err != nil {
		return fmt.Errorf("%s, %w", fn, err)
	}

	rowsCount, err := res.RowsAffected()
	fmt.Println(rowsCount)
	if err != nil {
		return fmt.Errorf("%s: get rows affected: %w", fn, err)
	}

	if rowsCount == 0 {
		return fmt.Errorf("%s, %w", fn, storage.ErrUrlNotFound)
	}

	return nil
}

func (s *SqliteStorage) GetUrl(alias string) (string, error) {
	fn := "internal.sqlite.GetUrl"

	stmt, err := s.db.Prepare("SELECT url FROM url WHERE alias = ?")
	if err != nil {
		return "", fmt.Errorf("%s, %w", fn, err)
	}
	defer stmt.Close()

	var resURL string
	err = stmt.QueryRow(alias).Scan(&resURL)

	if errors.Is(err, sql.ErrNoRows) {
		return "", fmt.Errorf("%s, %w", fn, storage.ErrUrlNotFound)
	}

	if err != nil {
		return "", fmt.Errorf("%s, %w", fn, err)
	}

	return resURL, nil
}
