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
		category TEXT DEFAULT '',
		priority TEXT DEFAULT '',
		severity TEXT DEFAULT '',
		status TEXT DEFAULT '',
		recommendation TEXT DEFAULT '',
		order_index INTEGER DEFAULT 0,
		report_id TEXT REFERENCES reports(id) ON DELETE CASCADE,
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

	CREATE INDEX IF NOT EXISTS idx_uploads_report_id ON uploads(report_id);
	CREATE INDEX IF NOT EXISTS idx_uploads_created_at ON uploads(created_at);
	CREATE INDEX IF NOT EXISTS idx_reports_created_at ON reports(created_at);

	CREATE TABLE IF NOT EXISTS blocks (
		id TEXT PRIMARY KEY,
		report_id TEXT NOT NULL REFERENCES reports(id) ON DELETE CASCADE,
		type TEXT NOT NULL,
		position INTEGER NOT NULL DEFAULT 0,
		content TEXT NOT NULL DEFAULT '{}',
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE INDEX IF NOT EXISTS idx_blocks_report_id ON blocks(report_id);
	CREATE INDEX IF NOT EXISTS idx_blocks_position ON blocks(report_id, position);
	`

	if _, err := s.db.Exec(schema); err != nil {
		return fmt.Errorf("run migrations: %w", err)
	}

	migrations := []string{
		`ALTER TABLE uploads ADD COLUMN category TEXT DEFAULT ''`,
		`ALTER TABLE uploads ADD COLUMN priority TEXT DEFAULT ''`,
		`ALTER TABLE uploads ADD COLUMN severity TEXT DEFAULT ''`,
		`ALTER TABLE uploads ADD COLUMN status TEXT DEFAULT ''`,
		`ALTER TABLE uploads ADD COLUMN recommendation TEXT DEFAULT ''`,
		`ALTER TABLE reports ADD COLUMN template_id TEXT DEFAULT ''`,
	}
	for _, m := range migrations {
		s.db.Exec(m)
	}

	slog.Info("database migrations completed")
	return nil
}
