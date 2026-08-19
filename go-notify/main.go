package main

import (
	"log"
	"os"
	"path/filepath"

	dbus "go-notify/internal/dbus"
	"go-notify/internal/models"
	"go-notify/internal/repository"
	"go-notify/internal/service"
	ui "go-notify/internal/ui"
	"go-notify/internal/util"
	"go-notify/internal/viewmodel"

	_ "github.com/mattn/go-sqlite3"

	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

const (
	stylePath = ".config/go-notify/style.css"
	dbPath    = ".local/share/go-notify/history.db"
)

func main() {
	if logFile, err := util.InitLogger("go-notify.log"); err == nil {
		defer logFile.Close()
	}

	notifChan := make(chan *models.Notification, 50)
	home, _ := os.UserHomeDir()

	dbUrl := filepath.Join(home, dbPath)
	os.MkdirAll(filepath.Dir(dbUrl), 0755)

	repo, err := repository.NewSQLiteHistoryRepository(dbUrl)
	if err != nil {
		log.Fatal(err)
	}
	defer repo.Close()

	launcherSvc := service.NewDesktopLauncherService()
	vm := viewmodel.NewNotificationViewModel(launcherSvc)

	app := gtk.NewApplication("com.github.june.notif-center", 0)
	app.ConnectActivate(func() {
		conn, err := dbus.StartServer(notifChan, repo)
		if err != nil {
			log.Fatal(err)
		}
		vm.Conn = conn

		data, err := os.ReadFile(filepath.Join(home, stylePath))
		if err != nil {
			log.Printf("[Go-Notify] %s", err)
		}
		win := ui.NewNotifWindow(app, data, vm)
		win.Listen(notifChan)
	})

	app.Run(os.Args)
}
