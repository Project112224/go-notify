package ui

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/diamondburned/gotk4/pkg/glib/v2"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
	"github.com/diamondburned/gotk4/pkg/pango"

	"go-notify/internal/models"
)

const (
	timeoutDuration = 5000
)

func (ctrl *NotifWindow) NewNotifCard(notif *models.Notification) *gtk.Box {
	card := gtk.NewBox(gtk.OrientationHorizontal, 12)
	card.SetSizeRequest(350, -1)
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
	icon.SetHExpand(false)
	return icon
}

func createTextContainer(summary string, content string) *gtk.Box {
	container := gtk.NewBox(gtk.OrientationVertical, 2)
	container.SetHExpand(true)
	container.SetVAlign(gtk.AlignCenter)

	title := gtk.NewLabel("")
	title.SetMarkup(fmt.Sprintf("<b>%s</b>", glib.MarkupEscapeText(summary)))
	title.AddCSSClass("white-label")
	title.SetXAlign(0)
	title.SetHExpand(true)
	title.SetWrap(true)
	title.SetWrapMode(pango.WrapWordChar)
	title.SetEllipsize(pango.EllipsizeEnd)
	title.SetLines(1)
	title.SetMaxWidthChars(30)

	body := gtk.NewLabel(content)
	body.AddCSSClass("white-label")
	body.SetXAlign(0)
	body.SetHExpand(true)
	body.SetWrap(true)
	body.SetLines(2)
	body.SetWrapMode(pango.WrapWordChar)
	body.SetEllipsize(pango.EllipsizeEnd)
	body.SetMaxWidthChars(40)

	container.Append(title)
	container.Append(body)
	return container
}

func (ctrl *NotifWindow) setCardGestures(card *gtk.Box, notif *models.Notification, isHovering *bool) {
	gesture := gtk.NewGestureClick()
	gesture.SetButton(1)
	gesture.ConnectPressed(func(nPress int, x, y float64) {
		log.Printf("[NotifCard] 點擊通知卡片 (ID: %d, App: %s, Summary: %s)", notif.ID, notif.AppName, notif.Summary)
		ctrl.emitDBusSignals(notif)
		urls := extractURLs(notif.Body)
		if len(urls) > 0 {
			log.Printf("[NotifCard] 偵測到 URL，嘗試開啟: %s", urls[0])
			openURL(urls[0])
		} else {
			openApplication(notif.AppName, notif.DesktopEntry)
		}
		ctrl.dismissCard(card)
	})
	card.AddController(gesture)
}

func openURL(rawURL string) {
	log.Printf("[openURL] 開始嘗試開啟網址: %s", rawURL)

	go func() {
		// 方法 1: xdg-open
		cmd1 := exec.Command("xdg-open", rawURL)
		cmd1.Env = os.Environ()
		out1, err1 := cmd1.CombinedOutput()
		if err1 == nil {
			log.Printf("[openURL] xdg-open 成功開啟: %s", rawURL)
			return
		}
		log.Printf("[openURL] xdg-open 失敗 (%v), 輸出: %s", err1, string(out1))

		// 方法 2: firefox
		cmd2 := exec.Command("firefox", rawURL)
		cmd2.Env = os.Environ()
		out2, err2 := cmd2.CombinedOutput()
		if err2 == nil {
			log.Printf("[openURL] firefox 成功開啟: %s", rawURL)
			return
		}
		log.Printf("[openURL] firefox 失敗 (%v), 輸出: %s", err2, string(out2))

		// 方法 3: gio open
		cmd3 := exec.Command("gio", "open", rawURL)
		cmd3.Env = os.Environ()
		out3, err3 := cmd3.CombinedOutput()
		if err3 == nil {
			log.Printf("[openURL] gio open 成功開啟: %s", rawURL)
			return
		}
		log.Printf("[openURL] gio open 失敗 (%v), 輸出: %s", err3, string(out3))
	}()
}

func extractURLs(text string) []string {
	re := regexp.MustCompile(`https?://[^\s<>"']+|www\.[^\s<>"']+\.[^\s<>"']+`)
	matches := re.FindAllString(text, -1)
	var urls []string
	seen := make(map[string]bool)
	for _, m := range matches {
		url := m
		if strings.HasPrefix(m, "www.") {
			url = "https://" + m
		}
		if !seen[url] {
			seen[url] = true
			urls = append(urls, url)
		}
	}
	return urls
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
	exactName := target + ".desktop"

	// 1. 精確比對：完全符合 <target>.desktop (如 firefox.desktop)
	for _, path := range searchPaths {
		fullPath := filepath.Join(path, exactName)
		if cmd := extractExecFromDesktopFile(fullPath); cmd != "" {
			log.Printf("從精確匹配 %s 找到執行指令: %s", exactName, cmd)
			return cmd
		}
	}

	// 2. ID / 反向網域比對：如 org.mozilla.firefox.desktop
	for _, path := range searchPaths {
		files, err := os.ReadDir(path)
		if err != nil {
			continue
		}
		for _, f := range files {
			nameLower := strings.ToLower(f.Name())
			if !strings.HasSuffix(nameLower, ".desktop") {
				continue
			}
			base := strings.TrimSuffix(nameLower, ".desktop")
			if base == target || strings.HasSuffix(base, "."+target) {
				fullPath := filepath.Join(path, f.Name())
				if cmd := extractExecFromDesktopFile(fullPath); cmd != "" {
					log.Printf("從 ID 匹配 %s 找到執行指令: %s", f.Name(), cmd)
					return cmd
				}
			}
		}
	}

	// 3. 避免誤觸衍生版本 (-developer, -nightly, -esr)，除非 target 原本就指定
	for _, path := range searchPaths {
		files, err := os.ReadDir(path)
		if err != nil {
			continue
		}
		for _, f := range files {
			nameLower := strings.ToLower(f.Name())
			if !strings.HasSuffix(nameLower, ".desktop") {
				continue
			}
			if strings.Contains(nameLower, target) {
				if !strings.Contains(target, "developer") && strings.Contains(nameLower, "developer") {
					continue
				}
				if !strings.Contains(target, "nightly") && strings.Contains(nameLower, "nightly") {
					continue
				}
				fullPath := filepath.Join(path, f.Name())
				if cmd := extractExecFromDesktopFile(fullPath); cmd != "" {
					log.Printf("從模糊匹配 %s 找到執行指令: %s", f.Name(), cmd)
					return cmd
				}
			}
		}
	}

	return target
}

func extractExecFromDesktopFile(fullPath string) string {
	content, err := os.ReadFile(fullPath)
	if err != nil {
		return ""
	}

	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "Exec=") {
			execLine := strings.TrimPrefix(line, "Exec=")
			return cleanExecLine(execLine)
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
