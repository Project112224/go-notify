package ui

import (
	"fmt"
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
	gesture.SetButton(0)
	gesture.ConnectPressed(func(nPress int, x, y float64) {
		button := gesture.CurrentButton()
		if ctrl.VM != nil {
			ctrl.VM.HandleCardClick(notif, button)
		}
		ctrl.dismissCard(card)
	})
	card.AddController(gesture)
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
