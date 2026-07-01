package repository

import (
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

type SQLite struct {
	db *sql.DB
}

func NewSQLite(path string) (*SQLite, error) {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("create data directory: %w", err)
	}

	db, err := sql.Open("sqlite3", path+"?_journal_mode=WAL&_foreign_keys=on")
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("ping database: %w", err)
	}

	slog.Info("database connected", "path", path)
	return &SQLite{db: db}, nil
}

func (s *SQLite) Close() error {
	return s.db.Close()
}

func (s *SQLite) DB() *sql.DB {
	return s.db
}

func (s *SQLite) Migrate() error {
	schema := `
	CREATE TABLE IF NOT EXISTS uploads (
		id TEXT PRIMARY KEY,
		filename TEXT NOT NULL,
		original_name TEXT NOT NULL,
		mime_type TEXT NOT NULL,
		size INTEGER NOT NULL,
		path TEXT NOT NULL,
		title TEXT DEFAULT '',
		description TEXT DEFAULT '',
		notes TEXT DEFAULT '',
		order_index INTEGER DEFAULT 0,
		report_id TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS reports (
		id TEXT PRIMARY KEY,
		title TEXT NOT NULL,
		project TEXT DEFAULT '',
		company TEXT DEFAULT '',
		author TEXT DEFAULT '',
		version TEXT DEFAULT '',
		status TEXT DEFAULT 'draft',
		pdf_path TEXT DEFAULT '',
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	`

	if _, err := s.db.Exec(schema); err != nil {
		return fmt.Errorf("run migrations: %w", err)
	}

	slog.Info("database migrations completed")
	return nil
}
