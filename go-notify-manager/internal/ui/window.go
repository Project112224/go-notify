package ui

import (
	"strings"

	"github.com/diamondburned/gotk4/pkg/gdk/v4"
	"github.com/diamondburned/gotk4/pkg/glib/v2"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

type Window struct {
	Window       *gtk.ApplicationWindow
	ListBox      *gtk.ListBox
	ClearButton  *gtk.Button
	SearchEntry  *gtk.SearchEntry
	SwitchButton *gtk.Switch
}

type HeaderResult struct {
	Header      *gtk.Box
	ClearBtn    *gtk.Button
	NotifSwitch *gtk.Switch
}

type ListViewResult struct {
	Scrolled *gtk.ScrolledWindow
	ListBox  *gtk.ListBox
}

func NewWindow(app *gtk.Application, cssData []byte) *Window {
	win := CreateWindow(app)
	SetSettings(cssData)

	safeArea := NewSafeArea()
	win.SetChild(safeArea)

	// Header 區域
	headerResult := NewHeader()
	safeArea.Append(headerResult.Header)

	searchEntry := NewSearchEntry()
	safeArea.Append(searchEntry)

	// 列表區域
	listViewResult := NewListView()
	safeArea.Append(listViewResult.Scrolled)

	SetGesture(win)

	return &Window{
		Window:       win,
		ListBox:      listViewResult.ListBox,
		ClearButton:  headerResult.ClearBtn,
		SearchEntry:  searchEntry,
		SwitchButton: headerResult.NotifSwitch,
	}
}

func CreateWindow(app *gtk.Application) *gtk.ApplicationWindow {
	win := gtk.NewApplicationWindow(app)
	win.SetTitle("通知歷史管理器")
	win.SetDefaultSize(450, 600)
	return win
}

func SetSettings(cssData []byte) {
	settings := gtk.SettingsGetDefault()
	settings.Object.SetObjectProperty("gtk-application-prefer-dark-theme", true)
	if len(cssData) > 0 {
		ApplyCustomCSS(cssData)
	}

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

func NewSafeArea() *gtk.Box {
	vBox := gtk.NewBox(gtk.OrientationVertical, 10)
	vBox.SetMarginTop(15)
	vBox.SetMarginBottom(15)
	vBox.SetMarginStart(15)
	vBox.SetMarginEnd(15)
	return vBox
}

func NewSearchEntry() *gtk.SearchEntry {
	searchEntry := gtk.NewSearchEntry()
	searchEntry.SetPlaceholderText("搜尋通知 (標題或內容)...")
	searchEntry.SetMarginBottom(10)
	return searchEntry
}

func NewHeader() *HeaderResult {
	header := gtk.NewBox(gtk.OrientationHorizontal, 5)

	notifSwitch := gtk.NewSwitch()
	notifSwitch.SetActive(true)

	titleLabel := gtk.NewLabel("歷史通知紀錄")
	titleLabel.SetHExpand(true)
	clearBtn := gtk.NewButtonWithLabel("清除全部")

	header.Append(notifSwitch)
	header.Append(titleLabel)
	header.Append(clearBtn)
	return &HeaderResult{
		Header:      header,
		ClearBtn:    clearBtn,
		NotifSwitch: notifSwitch,
	}
}

func NewListView() *ListViewResult {
	scrolled := gtk.NewScrolledWindow()
	scrolled.SetVExpand(true)
	scrolled.SetHExpand(true)
	listbox := gtk.NewListBox()
	listbox.SetSelectionMode(gtk.SelectionSingle)
	scrolled.SetChild(listbox)
	return &ListViewResult{
		Scrolled: scrolled,
		ListBox:  listbox,
	}
}

func CreateDateHeader(dateText string, listbox *gtk.ListBox) *gtk.Box {
	box := gtk.NewBox(gtk.OrientationVertical, 0)
	box.SetHExpand(true)
	box.AddCSSClass("date-header-box")

	sep := gtk.NewSeparator(gtk.OrientationHorizontal)
	sep.AddCSSClass("header-separator")
	sep.SetMarginStart(20)
	sep.SetMarginEnd(20)
	box.Append(sep)

	hbox := gtk.NewBox(gtk.OrientationHorizontal, 10)
	hbox.SetMarginStart(12)
	hbox.SetMarginTop(10)

	arrow := gtk.NewImageFromIconName("pan-down-symbolic")
	arrow.SetMarginBottom(2)
	label := gtk.NewLabel(dateText)
	label.SetYAlign(0.5)
	label.AddCSSClass("date-header-label")

	hbox.Append(arrow)
	hbox.Append(label)
	box.Append(hbox)

	click := gtk.NewGestureClick()
	hbox.AddController(click)

	isCollapsed := false
	click.ConnectReleased(func(n int, x, y float64) {
		isCollapsed = !isCollapsed

		// 切換箭頭圖示
		if isCollapsed {
			arrow.SetFromIconName("pan-end-symbolic")
		} else {
			arrow.SetFromIconName("pan-down-symbolic")
		}

		children := listbox.ObserveChildren()
		for i := uint(0); i < children.NItems(); i++ {
			item := children.Item(i)

			row, ok := item.Cast().(*gtk.ListBoxRow)
			if !ok {
				continue
			}

			name := ""
			rowChild := row.Child()
			if rowChild != nil {
				if box, ok := rowChild.(*gtk.Box); ok {
					name = box.Name()
				}
			}

			if len(name) >= 10 && name[:10] == dateText {
				if strings.HasPrefix(name, dateText) {
					if isCollapsed {
						row.SetChildVisible(false)
						row.SetSizeRequest(-1, 0)
						row.SetFocusable(false)
						row.AddCSSClass("collapsed-row")
					} else {
						row.SetChildVisible(true)
						row.SetSizeRequest(-1, -1)
						row.SetFocusable(true)
						row.RemoveCSSClass("collapsed-row")
					}
				}
			}
		}
	})

	return box
}

func SetGesture(win *gtk.ApplicationWindow) {
	click := gtk.NewGestureClick()
	click.SetButton(1)

	click.ConnectPressed(func(n int, x, y float64) {
		win.Window.SetFocus(nil)
	})
	win.AddController(click)
}
