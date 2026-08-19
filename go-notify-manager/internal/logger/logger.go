package logger

import (
	"io"
	"log"
	"os"
	"path/filepath"
)

func InitLogger(logFileName string) (*os.File, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	logDir := filepath.Join(home, ".local/share/go-notify")
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil, err
	}

	logPath := filepath.Join(logDir, logFileName)
	logFile, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, err
	}

	multiWriter := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(multiWriter)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	log.Printf("=== [%s] 日誌系統初始化完成 === (路徑: %s)", logFileName, logPath)
	return logFile, nil
}
