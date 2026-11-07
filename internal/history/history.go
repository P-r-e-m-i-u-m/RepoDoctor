package history

import (
	"database/sql"
	_ "modernc.org/sqlite"
	"os"
	"path/filepath"
	"time"
)

type Store struct{ db *sql.DB }

func DefaultPath() (string, error) {
	base, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	dir := filepath.Join(base, "repodoctor")
	if err := os.MkdirAll(dir, 0700); err != nil {
		return "", err
	}
	return filepath.Join(dir, "history.db"), nil
}
func Open(path string) (*Store, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS accepted_fixes (id INTEGER PRIMARY KEY AUTOINCREMENT, title TEXT NOT NULL, command TEXT NOT NULL, accepted_at TEXT NOT NULL)")
	if err != nil {
		db.Close()
		return nil, err
	}
	return &Store{db: db}, nil
}
func (s *Store) Record(title, command string, at time.Time) error {
	_, err := s.db.Exec("INSERT INTO accepted_fixes(title,command,accepted_at) VALUES(?,?,?)", title, command, at.UTC().Format(time.RFC3339))
	return err
}
func (s *Store) Close() error { return s.db.Close() }
