package main

import (
	"log"
	"os"
	"path/filepath"

	"go-notify-manager/internal/logger"
	"go-notify-manager/internal/repository"
	"go-notify-manager/internal/services"
	"go-notify-manager/internal/ui"
	"go-notify-manager/internal/viewmodel"

	"github.com/diamondburned/gotk4/pkg/gtk/v4"
	_ "github.com/mattn/go-sqlite3"
)

const (
	stylePath = ".config/go-notify/manager.css"
	dbPath    = ".local/share/go-notify/history.db"
)

func main() {
	if logFile, err := logger.InitLogger("go-notify-manager.log"); err == nil {
		defer logFile.Close()
	}

	app := gtk.NewApplication("com.github.june.notif-manager", 0)
	app.ConnectActivate(func() {
		home, _ := os.UserHomeDir()
		dbUrl := filepath.Join(home, dbPath)
		cssData, _ := os.ReadFile(filepath.Join(home, stylePath))

		repo, err := repository.NewSQLiteHistoryRepository(dbUrl)
		if err != nil {
			log.Fatal(err)
		}

		focusSvc := services.NewDBusFocusModeService()
		vm := viewmodel.NewHistoryViewModel(repo, focusSvc)

		win := ui.NewWindow(app, cssData, vm)

		kbService := ui.NewKeyboardService(win)
		kbService.BindShortcuts()

		vm.InitFocusMode()
		vm.LoadHistory()

		win.Window.Present()
	})

	app.Run(os.Args)
}
