package ui

import (
	"github.com/diamondburned/gotk4-layer-shell/pkg/gtk4layershell"
	"github.com/diamondburned/gotk4/pkg/gdk/v4"
	glib "github.com/diamondburned/gotk4/pkg/glib/v2"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
	"github.com/godbus/dbus/v5"

	"go-notify/internal/models"
)

type NotifWindow struct {
	List   *gtk.ListBox
	Window *gtk.ApplicationWindow
	Conn   *dbus.Conn
}

func NewNotifWindow(app *gtk.Application, cssData []byte) *NotifWindow {

	appWin := gtk.NewApplicationWindow(app)
	window := &appWin.Window

	if len(cssData) > 0 {
		ApplyCustomCSS(cssData)
	}

	window.SetTitle("Go-Notif-Center")

	gtk4layershell.InitForWindow(window)
	gtk4layershell.SetAnchor(window, gtk4layershell.LayerShellEdgeTop, true)
	gtk4layershell.SetAnchor(window, gtk4layershell.LayerShellEdgeRight, true)

	window.SetDefaultSize(350, -1)

	mainBox := gtk.NewBox(gtk.OrientationVertical, 10)
	mainBox.SetMarginTop(10)
	mainBox.SetMarginBottom(10)
	mainBox.SetMarginStart(10)
	mainBox.SetMarginEnd(10)

	list := gtk.NewListBox()
	list.SetSelectionMode(gtk.SelectionNone)

	mainBox.Append(list)
	window.SetChild(mainBox)
	window.Present()

	return &NotifWindow{Window: appWin, List: list}
}

func (ctrl *NotifWindow) Listen(notifChan chan *models.Notification) {
	go func() {
		for n := range notifChan {
			glib.IdleAdd(func() {
				if !ctrl.Window.IsVisible() {
					ctrl.Window.SetVisible(true)
				}
				card := ctrl.NewNotifCard(n)
				ctrl.List.Prepend(card)
			})
		}
	}()
}

func ApplyCustomCSS(cssData []byte) {
	provider := gtk.NewCSSProvider()
	gByte := glib.NewBytes(cssData)
	provider.LoadFromBytes(gByte)

	display := gdk.DisplayGetDefault()
	gtk.StyleContextAddProviderForDisplay(
		display,
		provider,
		gtk.STYLE_PROVIDER_PRIORITY_APPLICATION,
	)
}
