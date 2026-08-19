package repository

import (
	"database/sql"
	"time"

	_ "github.com/mattn/go-sqlite3"

	"go-notify/internal/models"
)

type HistoryRepository interface {
	Save(n *models.HistoryNotif) error
	Close() error
}

type SQLiteHistoryRepository struct {
	db *sql.DB
}

func NewSQLiteHistoryRepository(dbPath string) (*SQLiteHistoryRepository, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	query := `
    CREATE TABLE IF NOT EXISTS notifications (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        app_name TEXT,
        summary TEXT,
        body TEXT,
        urgency INTEGER,
        icon_path TEXT,
        created_at DATETIME
    );
    PRAGMA journal_mode=WAL;`

	if _, err := db.Exec(query); err != nil {
		return nil, err
	}

	return &SQLiteHistoryRepository{db: db}, nil
}

func (r *SQLiteHistoryRepository) Save(n *models.HistoryNotif) error {
	query := `INSERT INTO notifications (app_name, summary, body, urgency, icon_path, created_at)
              VALUES (?, ?, ?, ?, ?, ?)`
	_, err := r.db.Exec(query, n.AppName, n.Summary, n.Body, n.Urgency, n.Icon, time.Now())
	return err
}

func (r *SQLiteHistoryRepository) Close() error {
	return r.db.Close()
}
