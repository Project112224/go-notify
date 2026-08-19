package repository

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"

	"go-notify-manager/internal/models"
)

type HistoryRepository interface {
	LoadAll() ([]models.HistoryItem, error)
	MarkAllAsRead() error
	ClearAll() error
	DeleteOne(id int) error
	DeleteByDate(dateStr string) error
	Close() error
}

type SQLiteHistoryRepository struct {
	db *sql.DB
}

const tableName = "notifications"

func NewSQLiteHistoryRepository(dbPath string) (*SQLiteHistoryRepository, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	// Migration: ensure is_read column exists
	_, _ = db.Exec("ALTER TABLE notifications ADD COLUMN is_read INTEGER DEFAULT 0;")

	return &SQLiteHistoryRepository{db: db}, nil
}

func (r *SQLiteHistoryRepository) LoadAll() ([]models.HistoryItem, error) {
	query := fmt.Sprintf("SELECT id, app_name, summary, body, urgency, icon_path, created_at, COALESCE(is_read, 0) FROM %s ORDER BY created_at DESC", tableName)
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []models.HistoryItem
	for rows.Next() {
		var item models.HistoryItem
		var isReadInt int
		if err := rows.Scan(&item.ID, &item.AppName, &item.Summary, &item.Body, &item.Urgency, &item.IconPath, &item.CreatedAt, &isReadInt); err != nil {
			log.Printf("[HistoryRepository] Scan 失敗: %v", err)
			continue
		}
		item.IsRead = (isReadInt == 1)
		items = append(items, item)
	}
	return items, nil
}

func (r *SQLiteHistoryRepository) MarkAllAsRead() error {
	query := fmt.Sprintf("UPDATE %s SET is_read = 1 WHERE is_read = 0 OR is_read IS NULL", tableName)
	_, err := r.db.Exec(query)
	if err != nil {
		log.Printf("[HistoryRepository] 標記全部已讀失敗: %v", err)
		return err
	}
	log.Println("[HistoryRepository] 所有未讀通知已成功標記為已讀")
	return nil
}

func (r *SQLiteHistoryRepository) ClearAll() error {
	query := fmt.Sprintf("DELETE FROM %s", tableName)
	_, err := r.db.Exec(query)
	return err
}

func (r *SQLiteHistoryRepository) DeleteOne(id int) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id = ?", tableName)
	res, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	count, _ := res.RowsAffected()
	if count == 0 {
		log.Printf("[HistoryRepository] 警告：找不到 ID 為 %d 的資料", id)
	}
	return nil
}

func (r *SQLiteHistoryRepository) DeleteByDate(dateStr string) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE created_at LIKE ?", tableName)
	_, err := r.db.Exec(query, dateStr+"%")
	return err
}

func (r *SQLiteHistoryRepository) Close() error {
	return r.db.Close()
}
