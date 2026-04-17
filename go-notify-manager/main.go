package main

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"go-notify-manager/internal/database"
	service "go-notify-manager/internal/services"
	"go-notify-manager/internal/ui"

	"github.com/diamondburned/gotk4/pkg/gtk/v4"
	"github.com/godbus/dbus/v5"
	_ "github.com/mattn/go-sqlite3"
)

const (
	stylePath = ".config/go-notify/manager.css"
	dbPath    = ".local/share/go-notify/history.db"
)

func main() {

	app := gtk.NewApplication("com.github.june.notif-manager", 0)
	app.ConnectActivate(func() {
		home, _ := os.UserHomeDir()
		dbUrl := filepath.Join(home, dbPath)
		cssData, err := os.ReadFile(filepath.Join(home, stylePath))

		db, err := database.NewManagerDB(dbUrl)
		if err != nil {
			log.Fatal(err)
		}

		win := ui.NewWindow(app, cssData)

		kbService := service.NewKeyboardService(win)
		kbService.BindShortcuts()

		// ListHeader
		win.ListBox.SetHeaderFunc(func(row *gtk.ListBoxRow, before *gtk.ListBoxRow) {
			currentFullDate := ""
			if box, ok := row.Child().(*gtk.Box); ok {
				currentFullDate = box.Name()
			}

			if len(currentFullDate) < 10 {
				return
			}
			currentDate := currentFullDate[:10]

			if before == nil {
				row.SetHeader(ui.CreateDateHeader(currentDate, win.ListBox))
				return
			}

			beforeFullDate := ""
			if box, ok := before.Child().(*gtk.Box); ok {
				beforeFullDate = box.Name()
			}

			if len(beforeFullDate) < 10 {
				return
			}
			beforeDate := beforeFullDate[:10]

			if currentDate != beforeDate {
				row.SetHeader(ui.CreateDateHeader(currentDate, win.ListBox))
			} else {
				row.SetHeader(nil)
			}
		})

		// Filter
		win.ListBox.SetFilterFunc(func(row *gtk.ListBoxRow) bool {
			query := strings.ToLower(win.SearchEntry.Text())
			if query == "" {
				return true
			}
			target := strings.ToLower(row.Name())
			return strings.Contains(target, query)
		})

		win.SearchEntry.ConnectChanged(func() {
			win.ListBox.InvalidateFilter()
		})

		// SQL Data
		var loadData func()
		loadData = func() {
			for child := win.ListBox.FirstChild(); child != nil; child = win.ListBox.FirstChild() {
				win.ListBox.Remove(child)
			}

			rows, _ := db.LoadHistory(nil)
			defer rows.Close()

			for rows.Next() {
				var appName, sum, body, createdAt string
				var urgency int
				var id int

				rows.Scan(&id, &appName, &sum, &body, &urgency, &createdAt)
				currentID := id
				row := ui.NewHistoryRow(appName, sum, body, urgency, createdAt, func() {
					err := db.DeleteOne(currentID)
					if err != nil {
						log.Println(err)
					}
					loadData()
				})
				win.ListBox.Append(row)
			}
		}

		// Clear
		win.ClearButton.ConnectClicked(func() {
			db.ClearHistory()
			loadData()
		})

		// Switch Notify Status
		var isLocked bool
		conn, _ := dbus.SessionBus()
		obj := conn.Object("org.freedesktop.Notifications", "/org/freedesktop/Notifications")
		obj.Call("org.freedesktop.Notifications.GetFocusMode", 0).Store(&isLocked)
		win.SwitchButton.ConnectStateSet(func(state bool) bool {
			obj.Call("org.freedesktop.Notifications.SetFocusMode", 0, !state)
			if state {
				log.Println("已開啟")
			} else {
				log.Println("已關閉")
			}
			return false
		})

		loadData()
		win.Window.Present()
	})

	app.Run(os.Args)
}
