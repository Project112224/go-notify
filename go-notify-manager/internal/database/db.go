package database

import (
	"database/sql"
	"fmt"
)

type ManagerDB struct {
	Conn *sql.DB
}

const (
	tableName = "notifications"
)

func NewManagerDB(dbPath string) (*ManagerDB, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}
	return &ManagerDB{Conn: db}, nil
}

func (m *ManagerDB) LoadHistory(limit int) (*sql.Rows, error) {
	query := fmt.Sprintf("SELECT id, app_name, summary, body, urgency, created_at FROM %s ORDER BY created_at DESC LIMIT ?", tableName)
	return m.Conn.Query(query, limit)
}

func (m *ManagerDB) ClearHistory() error {
	query := fmt.Sprintf("DELETE FROM %s", tableName)
	_, err := m.Conn.Exec(query)
	return err
}

func (m *ManagerDB) DeleteOne(id int) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id = ?", tableName)
	res, err := m.Conn.Exec(query, id)
	if err != nil {
		return err
	}

	count, _ := res.RowsAffected()
	if count == 0 {
		fmt.Println("警告：找不到 ID 為", id, "的資料，請檢查 Table Schema")
	}
	return nil
}
