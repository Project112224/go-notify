// internal/dbus/server.go
package dbus

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"os"
	"strings"
	"sync/atomic"
	"time"

	"github.com/godbus/dbus/v5"

	"go-notify/internal/database"
	"go-notify/internal/models"
)

var notifyID uint32

type NotificationServer struct {
	conn      *dbus.Conn
	notifChan chan *models.Notification
	dbService *database.DBService
	Locked    bool
}

func StartServer(notifChan chan *models.Notification, abService *database.DBService) (*dbus.Conn, error) {
	conn, err := dbus.ConnectSessionBus()
	if err != nil {
		return nil, err
	}

	server := &NotificationServer{
		conn:      conn,
		notifChan: notifChan,
		dbService: abService,
	}
	conn.Export(server, "/org/freedesktop/Notifications", "org.freedesktop.Notifications")

	reply, err := conn.RequestName("org.freedesktop.Notifications", dbus.NameFlagReplaceExisting)
	if err != nil || reply != dbus.RequestNameReplyPrimaryOwner {
		return nil, fmt.Errorf("無法取得 D-Bus 名稱")
	}

	return conn, nil
}

// TODO: 核心業務邏輯

// test: gdbus call --session --dest org.freedesktop.Notifications --object-path /org/freedesktop/Notifications --method org.freedesktop.Notifications.SetFocusMode true
func (s *NotificationServer) SetFocusMode(enabled bool) *dbus.Error {
	s.Locked = enabled
	fmt.Printf("通知伺服器狀態更新：鎖定=%v\n", enabled)
	return nil
}

func (s *NotificationServer) GetFocusMode() (bool, *dbus.Error) {
	return s.Locked, nil
}

func (s *NotificationServer) Notify(app_name string, replaces_id uint32, app_icon string, summary string, body string, actions []string, hints map[string]dbus.Variant, expire_timeout int32) (uint32, *dbus.Error) {

	if s.Locked {
		return 0, nil
	}
	// log
	for key, variant := range hints {
		strVal := variantToString(variant)
		log.Printf("Hint Key: %s, Value: %s", key, strVal)
	}

	// 處理 id
	atomic.AddUint32(&notifyID, 1)
	currentID := atomic.LoadUint32(&notifyID)

	urgency := getUrgency(hints)
	desktopEntry := getDesktopEntry(hints)
	defaultKey := getDefaultActionKey(actions)
	appIcon := getAppIcon(hints, app_icon)

	s.notifChan <- &models.Notification{
		AppName:          app_name,
		ID:               currentID,
		Summary:          summary,
		Body:             body,
		Urgency:          urgency,
		Icon:             appIcon,
		DefaultActionKey: defaultKey,
		DesktopEntry:     desktopEntry,
	}

	s.dbService.Save(&models.HistoryNotif{
		AppName: app_name,
		Summary: summary,
		Body:    body,
		Urgency: int(urgency),
		Icon:    appIcon,
		Time:    time.Now(),
	})

	return currentID, nil
}

// Low: 0, Normal: 1, High: 2
func getUrgency(hints map[string]dbus.Variant) uint8 {
	urgency := uint8(1)
	if v, ok := hints["urgency"]; ok {
		if u, ok := v.Value().(uint8); ok {
			urgency = u
		}
	}
	return urgency
}

func getDesktopEntry(hints map[string]dbus.Variant) string {
	desktopEntry := ""
	if v, ok := hints["desktop-entry"]; ok {
		if de, ok := v.Value().(string); ok {
			desktopEntry = de
		}
	}
	return desktopEntry
}

func getDefaultActionKey(actions []string) string {
	defaultKey := ""
	if len(actions) >= 2 {
		defaultKey = actions[0]
		for i := 0; i < len(actions)-1; i += 2 {
			if actions[i] == "default" {
				defaultKey = "default"
				break
			}
		}
	}
	return defaultKey
}

func getAppIcon(hints map[string]dbus.Variant, appIcon string) string {
	icon := appIcon
	if imgPath, ok := hints["image-path"]; ok {
		path := imgPath.Value().(string)
		if strings.HasPrefix(path, "file://") {
			icon = path[7:]
		} else {
			icon = path
		}
	} else if imgData, ok := hints["image-data"]; ok {
		if path, err := saveImageDataToTmp(imgData); err == nil {
			icon = path
		}
	}
	return icon
}

// 將 Notify Hints 轉換為字串
func variantToString(v dbus.Variant) string {
	val := v.Value()

	switch t := val.(type) {
	case string:
		return t
	case []byte:
		b := make([]byte, len(t))
		for i, v := range t {
			b[i] = byte(v)
		}
		return string(b)
	default:
		return fmt.Sprintf("%v", val)
	}
}

func saveImageDataToTmp(hint dbus.Variant) (string, error) {
	data, ok := hint.Value().([]interface{})
	if !ok || len(data) < 7 {
		return "", fmt.Errorf("invalid image-data format")
	}

	width := int(data[0].(int32))
	height := int(data[1].(int32))
	rowstride := int(data[2].(int32))
	hasAlpha := data[3].(bool)
	// bitsPerSample := data[4].(int32)
	channels := int(data[5].(int32))
	pixelBytes := data[6].([]byte)

	img := image.NewNRGBA(image.Rect(0, 0, width, height))

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			srcIdx := y*rowstride + x*channels

			if srcIdx+channels <= len(pixelBytes) {
				r := pixelBytes[srcIdx]
				g := pixelBytes[srcIdx+1]
				b := pixelBytes[srcIdx+2]
				var a uint8 = 255
				if hasAlpha && channels == 4 {
					a = pixelBytes[srcIdx+3]
				}
				img.SetNRGBA(x, y, color.NRGBA{R: r, G: g, B: b, A: a})
			}
		}
	}

	tmpPath := fmt.Sprintf("/tmp/notif_img_%d.png", time.Now().UnixNano())
	f, err := os.Create(tmpPath)
	if err != nil {
		return "", err
	}
	defer f.Close()

	if err := png.Encode(f, img); err != nil {
		return "", err
	}

	return tmpPath, nil
}

// TODO: D-Bus 協定規範方法

func (s *NotificationServer) GetCapabilities() ([]string, *dbus.Error) {
	return []string{"actions", "body", "body-markup", "icon-static"}, nil
}

func (s *NotificationServer) GetServerInformation() (string, string, string, string, *dbus.Error) {
	return "go-notif", "june", "0.1", "1.2", nil
}

func (s *NotificationServer) CloseNotification(id uint32) *dbus.Error {
	return nil
}
