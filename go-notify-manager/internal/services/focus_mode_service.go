package services

import (
	"log"
	"os/exec"

	"github.com/godbus/dbus/v5"
)

type FocusModeService interface {
	GetFocusMode() (isLocked bool, err error)
	SetFocusMode(enabled bool) error
	SetSystemMute(muted bool) error
}

type DBusFocusModeService struct{}

func NewDBusFocusModeService() *DBusFocusModeService {
	return &DBusFocusModeService{}
}

func (s *DBusFocusModeService) GetFocusMode() (bool, error) {
	conn, err := dbus.SessionBus()
	if err != nil {
		return false, err
	}
	obj := conn.Object("org.freedesktop.Notifications", "/org/freedesktop/Notifications")

	var isLocked bool
	err = obj.Call("org.freedesktop.Notifications.GetFocusMode", 0).Store(&isLocked)
	if err != nil {
		return false, err
	}
	return isLocked, nil
}

func (s *DBusFocusModeService) SetFocusMode(enabled bool) error {
	conn, err := dbus.SessionBus()
	if err != nil {
		return err
	}
	obj := conn.Object("org.freedesktop.Notifications", "/org/freedesktop/Notifications")

	// SetFocusMode takes locked status (true = focus mode active / notifications locked)
	callErr := obj.Call("org.freedesktop.Notifications.SetFocusMode", 0, !enabled).Err
	if callErr != nil {
		log.Printf("[FocusModeService] 設定 FocusMode 失敗: %v", callErr)
	}

	return s.SetSystemMute(!enabled)
}

func (s *DBusFocusModeService) SetSystemMute(muted bool) error {
	val := "0"
	if muted {
		val = "1"
	}
	cmd := exec.Command("pactl", "set-sink-mute", "@DEFAULT_SINK@", val)
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("[FocusModeService] 靜音切換失敗 (%v): %s", err, string(out))
		return err
	}
	return nil
}
