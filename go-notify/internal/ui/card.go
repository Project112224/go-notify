package ui

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/diamondburned/gotk4/pkg/core/glib"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
	"github.com/diamondburned/gotk4/pkg/pango"

	"go-notify/internal/models"
)

const (
	timeoutDuration = 5000
)

func (ctrl *NotifWindow) NewNotifCard(notif *models.Notification) *gtk.Box {
	card := gtk.NewBox(gtk.OrientationHorizontal, 12)
	isHovering := false
	setUrgencyClass(card, notif.Urgency)

	if notif.Icon != "" {
		card.Append(createIconWidget(notif.Icon))
	}
	textContainer := createTextContainer(notif.Summary, notif.Body)
	card.Append(textContainer)

	ctrl.setCardGestures(card, notif, &isHovering)
	ctrl.setCardRemoveTimer(&isHovering, card)
	return card
}

func setUrgencyClass(card *gtk.Box, urgency byte) {
	switch urgency {
	// Low
	case 0:
		card.AddCSSClass("urgency-low")
	// Critical
	case 2:
		card.AddCSSClass("urgency-critical")
	default:
		card.AddCSSClass("notif-card")
	}
}

func createIconWidget(iconPath string) *gtk.Image {
	var icon *gtk.Image

	if strings.HasPrefix(iconPath, "/") {
		icon = gtk.NewImageFromFile(iconPath)
	} else {
		icon = gtk.NewImageFromIconName(iconPath)
	}
	icon.SetPixelSize(40)
	icon.SetVAlign(gtk.AlignStart)
	return icon
}

func createTextContainer(summary string, content string) *gtk.Box {
	container := gtk.NewBox(gtk.OrientationVertical, 2)

	title := gtk.NewLabel("")
	title.SetMarkup(fmt.Sprintf("<b>%s</b>", summary))
	title.SetXAlign(0)
	title.SetWrap(true)
	title.SetWrapMode(pango.WrapWordChar)
	title.SetEllipsize(pango.EllipsizeEnd)
	title.SetMaxWidthChars(20)

	body := gtk.NewLabel(content)
	body.SetXAlign(0)
	body.SetWrap(true)
	body.SetLines(5)
	body.SetWrapMode(pango.WrapWordChar)
	body.SetEllipsize(pango.EllipsizeEnd)
	body.SetMaxWidthChars(100)

	container.Append(title)
	container.Append(body)
	return container
}

func (ctrl *NotifWindow) setCardGestures(card *gtk.Box, notif *models.Notification, isHovering *bool) {
	gesture := gtk.NewGestureClick()
	gesture.SetButton(0)
	gesture.ConnectPressed(func(nPress int, x, y float64) {
		button := gesture.CurrentButton()
		if button == 1 {
			ctrl.emitDBusSignals(notif)
			openApplication(notif.AppName, notif.DesktopEntry)
		}
		ctrl.dismissCard(card)
	})
	card.AddController(gesture)
}

func (ctrl *NotifWindow) emitDBusSignals(notif *models.Notification) {
	ctrl.Conn.Emit("/org/freedesktop/Notifications",
		"org.freedesktop.Notifications.ActionInvoked",
		uint32(notif.ID), notif.DefaultActionKey,
	)

	ctrl.Conn.Emit("/org/freedesktop/Notifications",
		"org.freedesktop.Notifications.NotificationClosed",
		uint32(notif.ID), uint32(2),
	)
}

func openApplication(appName string, desktopEntry string) {
	target := appName
	if desktopEntry != "" {
		target = desktopEntry
	}

	if isSystemApp(target) {
		return
	}

	go tryFocusOrLaunch(target)
}

func isSystemApp(appName string) bool {
	systemApps := map[string]bool{
		"system": true, "notify-send": true, "networkmanager": true,
		"power-profiles-daemon": true, "xdg-desktop-portal-hyprland": true,
	}
	return systemApps[appName]
}

func tryFocusOrLaunch(appName string) {
	target := strings.ToLower(appName)

	psCmd := exec.Command("sh", "-c", fmt.Sprintf("ps aux | grep -i %s | grep -v grep", target))
	err := psCmd.Run()

	if err == nil {
		log.Printf("發現正在執行的進程: %s，嘗試從 .desktop 獲取指令喚醒...", target)
	} else {
		log.Printf("未發現進程: %s，嘗試從 .desktop 啟動新實體...", target)
	}

	execCmd := findAndExecDesktop(target)
	if execCmd != "" {
		log.Printf("執行指令: %s", execCmd)
		go exec.Command("sh", "-c", execCmd+" &").Run()
	} else {
		go exec.Command("go-notify-manager").Run()
	}
}

func findAndExecDesktop(appName string) string {
	searchPaths := []string{
		"/usr/share/applications",
		"/usr/local/share/applications",
		filepath.Join(os.Getenv("HOME"), ".local/share/applications"),
		"/var/lib/flatpak/exports/share/applications",
	}

	target := strings.ToLower(appName)

	for _, path := range searchPaths {
		files, err := os.ReadDir(path)
		if err != nil {
			continue
		}

		for _, f := range files {
			if strings.Contains(strings.ToLower(f.Name()), target) && strings.HasSuffix(f.Name(), ".desktop") {
				fullPath := filepath.Join(path, f.Name())
				content, err := os.ReadFile(fullPath)
				if err != nil {
					continue
				}

				lines := strings.Split(string(content), "\n")
				for _, line := range lines {
					if strings.HasPrefix(line, "Exec=") {
						execLine := strings.TrimPrefix(line, "Exec=")
						finalCmd := cleanExecLine(execLine)
						log.Printf("從 %s 找到執行指令: %s", f.Name(), finalCmd)
						return finalCmd
					}
				}
			}
		}
	}
	return ""
}

func cleanExecLine(line string) string {
	re := regexp.MustCompile(`%[uUfFiIcKknNvV]`)
	cleaned := re.ReplaceAllString(line, "")
	return strings.TrimSpace(cleaned)
}

func (ctrl *NotifWindow) dismissCard(card *gtk.Box) {
	parentRow := card.Parent()
	if parentRow != nil {
		ctrl.List.Remove(parentRow)
	}
	if ctrl.List.FirstChild() == nil {
		ctrl.Window.SetVisible(false)
	}
}

func (ctrl *NotifWindow) setCardRemoveTimer(isHovering *bool, card *gtk.Box) {
	var status = *isHovering
	var timerID glib.SourceHandle
	var startAutoDismiss = func(delay uint) {
		if timerID > 0 {
			glib.SourceRemove(timerID)
			timerID = 0
		}

		timerID = glib.TimeoutAdd(delay, func() bool {
			if status {
				return false
			}

			parentRow := card.Parent()
			if parentRow != nil {
				ctrl.List.Remove(parentRow)
			}
			if ctrl.List.FirstChild() == nil {
				ctrl.Window.SetVisible(false)
			}

			glib.SourceRemove(timerID)
			timerID = 0
			return false
		})

	}

	motion := gtk.NewEventControllerMotion()
	motion.ConnectEnter(func(x, y float64) {
		if timerID > 0 {
			glib.SourceRemove(timerID)
			timerID = 0
		}
	})
	motion.ConnectLeave(func() {
		status = false
		startAutoDismiss(timeoutDuration)
	})
	card.AddController(motion)
	startAutoDismiss(timeoutDuration)
}
