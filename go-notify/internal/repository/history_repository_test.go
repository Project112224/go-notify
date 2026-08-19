package repository

import (
	"path/filepath"
	"testing"
	"time"

	"go-notify/internal/models"
)

func TestSQLiteHistoryRepository_SaveAndClose(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test_history.db")

	repo, err := NewSQLiteHistoryRepository(dbPath)
	if err != nil {
		t.Fatalf("Failed to create repository: %v", err)
	}

	notif := &models.HistoryNotif{
		AppName:    "TestApp",
		ReplacesId: 0,
		Summary:    "Test Summary",
		Body:       "Test Body",
		Urgency:    1,
		Time:       time.Now(),
		Icon:       "dialog-information",
	}

	if err := repo.Save(notif); err != nil {
		t.Errorf("Failed to save notification: %v", err)
	}

	if err := repo.Close(); err != nil {
		t.Errorf("Failed to close repository: %v", err)
	}
}
