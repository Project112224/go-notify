package database

import (
	"database/sql"
	"time"

	_ "github.com/mattn/go-sqlite3"

	"go-notify/internal/models"
)

type DBService struct {
	db *sql.DB
}

func NewDBService(dbPath string) (*DBService, error) {
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

	return &DBService{db}, nil
}

func (s *DBService) Save(n *models.HistoryNotif) error {
	query := `INSERT INTO notifications (app_name, summary, body, urgency, icon_path, created_at)
              VALUES (?, ?, ?, ?, ?, ?)`
	_, err := s.db.Exec(query, n.AppName, n.Summary, n.Body, n.Urgency, n.Icon, time.Now())
	return err
}

func (s *DBService) Close() {
	s.db.Close()
}
