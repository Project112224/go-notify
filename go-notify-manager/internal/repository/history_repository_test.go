package repository

import (
	"database/sql"
	"path/filepath"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func setupTestDB(t *testing.T) (*SQLiteHistoryRepository, string) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test_manager.db")

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		t.Fatalf("Failed to open sqlite db: %v", err)
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
    );`
	if _, err := db.Exec(query); err != nil {
		t.Fatalf("Failed to init table: %v", err)
	}

	nowStr := time.Now().Format(time.RFC3339)
	insertQuery := `INSERT INTO notifications (app_name, summary, body, urgency, icon_path, created_at)
                    VALUES (?, ?, ?, ?, ?, ?)`
	db.Exec(insertQuery, "Firefox", "Download done", "file.zip", 1, "firefox", nowStr)
	db.Exec(insertQuery, "Discord", "Message", "Hello", 2, "discord", nowStr)
	db.Close()

	repo, err := NewSQLiteHistoryRepository(dbPath)
	if err != nil {
		t.Fatalf("Failed to create repo: %v", err)
	}

	return repo, dbPath
}

func TestSQLiteHistoryRepository_LoadAllAndDelete(t *testing.T) {
	repo, _ := setupTestDB(t)
	defer repo.Close()

	items, err := repo.LoadAll()
	if err != nil {
		t.Fatalf("LoadAll failed: %v", err)
	}
	if len(items) != 2 {
		t.Errorf("Expected 2 items, got %d", len(items))
	}

	// Test DeleteOne
	if err := repo.DeleteOne(items[0].ID); err != nil {
		t.Errorf("DeleteOne failed: %v", err)
	}

	remaining, _ := repo.LoadAll()
	if len(remaining) != 1 {
		t.Errorf("Expected 1 item remaining, got %d", len(remaining))
	}

	// Test ClearAll
	if err := repo.ClearAll(); err != nil {
		t.Errorf("ClearAll failed: %v", err)
	}

	afterClear, _ := repo.LoadAll()
	if len(afterClear) != 0 {
		t.Errorf("Expected 0 items after ClearAll, got %d", len(afterClear))
	}
}
