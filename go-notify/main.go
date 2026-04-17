package main

import (
	"log"
	"os"
	"path/filepath"

	"go-notify/internal/database"
	dbus "go-notify/internal/dbus"
	ui "go-notify/internal/ui"

	"go-notify/internal/models"

	_ "github.com/mattn/go-sqlite3"

	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

const (
	stylePath = ".config/go-notify/style.css"
	dbPath    = ".local/share/go-notify/history.db"
)

func main() {
	notifChan := make(chan *models.Notification, 50)
	home, _ := os.UserHomeDir()

	dbUrl := filepath.Join(home, dbPath)
	os.MkdirAll(dbUrl, 0755)
	dbSvc, err := database.NewDBService(dbUrl)
	if err != nil {
		log.Fatal(err)
	}
	defer dbSvc.Close()

	app := gtk.NewApplication("com.github.june.notif-center", 0)
	app.ConnectActivate(func() {
		conn, err := dbus.StartServer(notifChan, dbSvc)
		if err != nil {
			log.Fatal(err)
		}

		data, err := os.ReadFile(filepath.Join(home, stylePath))
		if err != nil {
			log.Print("[Go-Notify] %s", err)
		}
		win := ui.NewNotifWindow(app, data)
		win.Conn = conn
		win.Listen(notifChan)
	})

	app.Run(os.Args)
}
