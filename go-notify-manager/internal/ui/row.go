package ui

import (
	"fmt"
	"regexp"
	"strings"

	extension "go-notify-manager/internal/extension"

	"github.com/diamondburned/gotk4/pkg/glib/v2"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
	"github.com/diamondburned/gotk4/pkg/pango"
)

func NewHistoryRow(app, sum, body string, urgency int, timeStr string, onDelete func()) *gtk.ListBoxRow {

	row := NewRow(timeStr)
	fullContent := fmt.Sprintf("%s %s %s", app, sum, body)
	row.SetName(fullContent)

	hbox := NewHBox()
	hbox.SetName(timeStr)

	titleLabel := NewTitleLabel(sum, urgency)
	iconName := getIconName(app)
	img := NewIcon(iconName)

	contentBox := NewContentBox()

	bodyLabel := NewBodyLabel(body)

	gesture := gtk.NewGestureClick()
	gesture.SetPropagationPhase(gtk.PhaseCapture)
	gesture.ConnectPressed(func(n int, x, y float64) {
		bodyLabel.GrabFocus()
	})
	bodyLabel.AddController(gesture)

	timeLabel := NewTimeLabel(timeStr)

	deleteBtn := NewDeleteButton(onDelete)

	// Append
	contentBox.Append(titleLabel)
	contentBox.Append(bodyLabel)
	contentBox.Append(timeLabel)
	hbox.Append(img)
	hbox.Append(contentBox)
	hbox.Append(deleteBtn)
	row.SetChild(hbox)
	return row
}

func linkify(text string) string {
	re := regexp.MustCompile(`(https?://[^\s]+)`)
	return re.ReplaceAllString(text, `<a href="$1">$1</a>`)
}

func getIconName(app string) string {
	switch strings.ToLower(app) {
	case "discord":
		return "discord"
	case "firefox":
		return "firefox"
	case "google-chrome":
		return "google-chrome"
	default:
		return "dialog-information"
	}
}

func NewRow(timeStr string) *gtk.ListBoxRow {
	row := gtk.NewListBoxRow()
	row.SetCanFocus(true)
	row.SetFocusable(true)
	row.SetSelectable(true)
	row.AddCSSClass("history-row")
	row.AddCSSClass("is-content")
	rowDate := extension.DateString(timeStr).FormatToHMS("2006-01-02")
	row.AddCSSClass("date-" + rowDate)
	return row
}

func NewHBox() *gtk.Box {
	hbox := gtk.NewBox(gtk.OrientationHorizontal, 12)
	hbox.SetMarginStart(15)
	hbox.SetMarginEnd(15)
	hbox.SetMarginTop(18)
	hbox.SetMarginBottom(18)

	return hbox
}

func NewTitleLabel(sum string, urgency int) *gtk.Label {
	titleLabel := gtk.NewLabel("")
	titleLabel.SetSelectable(true)
	titleLabel.SetMarkup(fmt.Sprintf("<span size='medium' weight='bold'>%s</span>", sum))
	titleLabel.SetXAlign(0)
	if urgency == 2 {
		titleLabel.SetMarkup(fmt.Sprintf("<span foreground='red'><b>%s</b></span>", sum))
	}

	return titleLabel
}

func NewIcon(iconName string) *gtk.Image {
	img := gtk.NewImageFromIconName(iconName)
	img.SetPixelSize(48)
	return img
}

func NewContentBox() *gtk.Box {
	vbox := gtk.NewBox(gtk.OrientationVertical, 2)
	vbox.SetHExpand(true)
	vbox.SetHAlign(gtk.AlignStart)
	return vbox
}

func NewBodyLabel(body string) *gtk.Label {
	escapedBody := glib.MarkupEscapeText(body)

	bodyLabel := gtk.NewLabel(escapedBody)
	bodyLabel.SetUseMarkup(true)
	bodyLabel.SetMarkup(linkify(escapedBody))

	bodyLabel.SetWrap(true)
	bodyLabel.SetWrapMode(pango.WrapWordChar)
	bodyLabel.SetLines(0)
	bodyLabel.SetEllipsize(pango.EllipsizeNone)
	bodyLabel.SetMaxWidthChars(50)
	bodyLabel.SetHExpand(true)
	bodyLabel.SetXAlign(0)
	bodyLabel.AddCSSClass("dim-label")

	// bodyLabel.SetSelectable(true)
	bodyLabel.ConnectActivateLink(func(uri string) bool {
		return false
	})
	return bodyLabel
}

func NewTimeLabel(displayTime string) *gtk.Label {
	convertedTime := extension.DateString(displayTime).FormatToHMS()
	timeLabel := gtk.NewLabel(convertedTime)
	timeLabel.SetSelectable(true)
	timeLabel.SetVAlign(gtk.AlignStart)
	timeLabel.SetXAlign(0)
	timeLabel.AddCSSClass("time-label")
	return timeLabel
}

func NewDeleteButton(onDelete func()) *gtk.Button {
	deleteBtn := gtk.NewButtonFromIconName("user-trash-symbolic")
	deleteBtn.AddCSSClass("delete-button")
	deleteBtn.SetVAlign(gtk.AlignCenter)
	deleteBtn.SetHAlign(gtk.AlignEnd)
	deleteBtn.ConnectClicked(func() {
		onDelete()
	})
	return deleteBtn
}
