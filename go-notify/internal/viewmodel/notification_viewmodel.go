package viewmodel

import (
	"log"

	"github.com/godbus/dbus/v5"

	"go-notify/internal/models"
	"go-notify/internal/service"
)

type NotificationViewModel struct {
	launcherSvc service.LauncherService
	Conn        *dbus.Conn
}

func NewNotificationViewModel(launcherSvc service.LauncherService) *NotificationViewModel {
	return &NotificationViewModel{
		launcherSvc: launcherSvc,
	}
}

func (vm *NotificationViewModel) HandleCardClick(notif *models.Notification, button uint) {
	log.Printf("[NotifViewModel] 處理卡片點擊 (Button: %d, ID: %d, App: %s, Summary: %s)", button, notif.ID, notif.AppName, notif.Summary)

	if button == 1 {
		vm.emitActionInvoked(notif)
		vm.emitNotificationClosed(notif)

		urls := vm.launcherSvc.ExtractURLs(notif.Body)
		if len(urls) > 0 {
			log.Printf("[NotifViewModel] 偵測到 URL，嘗試開啟: %s", urls[0])
			vm.launcherSvc.OpenURL(urls[0])
		} else {
			vm.launcherSvc.OpenApplication(notif.AppName, notif.DesktopEntry)
		}
	} else {
		vm.emitNotificationClosed(notif)
	}
}

func (vm *NotificationViewModel) emitActionInvoked(notif *models.Notification) {
	if vm.Conn == nil {
		return
	}
	vm.Conn.Emit("/org/freedesktop/Notifications",
		"org.freedesktop.Notifications.ActionInvoked",
		uint32(notif.ID), notif.DefaultActionKey,
	)
}

func (vm *NotificationViewModel) emitNotificationClosed(notif *models.Notification) {
	if vm.Conn == nil {
		return
	}
	vm.Conn.Emit("/org/freedesktop/Notifications",
		"org.freedesktop.Notifications.NotificationClosed",
		uint32(notif.ID), uint32(2),
	)
}
