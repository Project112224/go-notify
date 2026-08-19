package ui

import (
	"github.com/diamondburned/gotk4/pkg/gdk/v4"
	"github.com/diamondburned/gotk4/pkg/glib/v2"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"

	"go-notify-manager/internal/models"
	"go-notify-manager/internal/viewmodel"
)

type Window struct {
	Window       *gtk.ApplicationWindow
	ListBox      *gtk.ListBox
	ClearButton  *gtk.Button
	SearchEntry  *gtk.SearchEntry
	SwitchButton *gtk.Switch
	VM           *viewmodel.HistoryViewModel
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

func NewWindow(app *gtk.Application, cssData []byte, vm *viewmodel.HistoryViewModel) *Window {
	win := CreateWindow(app)
	SetSettings(cssData)

	safeArea := NewSafeArea()
	win.SetChild(safeArea)

	headerResult := NewHeader()
	safeArea.Append(headerResult.Header)

	searchEntry := NewSearchEntry()
	safeArea.Append(searchEntry)

	listViewResult := NewListView()
	safeArea.Append(listViewResult.Scrolled)

	SetGesture(win)

	w := &Window{
		Window:       win,
		ListBox:      listViewResult.ListBox,
		ClearButton:  headerResult.ClearBtn,
		SearchEntry:  searchEntry,
		SwitchButton: headerResult.NotifSwitch,
		VM:           vm,
	}

	w.BindViewModel()
	return w
}

func (w *Window) BindViewModel() {
	if w.VM == nil {
		return
	}

	w.VM.OnItemsUpdated = func(items []models.HistoryItem) {
		glib.IdleAdd(func() {
			w.renderItems(items)
		})
	}

	w.VM.OnFocusModeUpdated = func(enabled bool) {
		glib.IdleAdd(func() {
			w.SwitchButton.SetActive(enabled)
		})
	}

	w.SearchEntry.ConnectChanged(func() {
		w.VM.SetSearchQuery(w.SearchEntry.Text())
	})

	w.ClearButton.ConnectClicked(func() {
		w.VM.ClearAll()
	})

	w.SwitchButton.ConnectStateSet(func(state bool) bool {
		w.VM.ToggleNotificationStatus(state)
		return false
	})

	w.ListBox.SetHeaderFunc(func(row *gtk.ListBoxRow, before *gtk.ListBoxRow) {
		currentFullDate := ""
		if box, ok := row.Child().(*gtk.Box); ok {
			currentFullDate = box.Name()
		}

		if len(currentFullDate) < 10 {
			return
		}
		currentDate := currentFullDate[:10]

		onDeleteDate := func() {
			w.VM.DeleteDateGroup(currentDate)
		}

		if before == nil {
			row.SetHeader(CreateDateHeader(currentDate, w.ListBox, onDeleteDate))
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
			row.SetHeader(CreateDateHeader(currentDate, w.ListBox, onDeleteDate))
		} else {
			row.SetHeader(nil)
		}
	})
}

func (w *Window) renderItems(items []models.HistoryItem) {
	for child := w.ListBox.FirstChild(); child != nil; child = w.ListBox.FirstChild() {
		w.ListBox.Remove(child)
	}

	for _, item := range items {
		itemID := item.ID
		row := NewHistoryRow(item.AppName, item.Summary, item.Body, item.IconPath, item.Urgency, item.CreatedAt, func() {
			w.VM.DeleteSingleItem(itemID)
		})
		w.ListBox.Append(row)
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
	vBox := gtk.NewBox(gtk.OrientationVertical, 8)
	vBox.SetMarginTop(12)
	vBox.SetMarginBottom(12)
	vBox.SetMarginStart(12)
	vBox.SetMarginEnd(12)
	return vBox
}

func NewSearchEntry() *gtk.SearchEntry {
	searchEntry := gtk.NewSearchEntry()
	searchEntry.SetPlaceholderText("搜尋通知 (標題或內容)...")
	searchEntry.SetMarginBottom(8)
	return searchEntry
}

func NewHeader() *HeaderResult {
	header := gtk.NewBox(gtk.OrientationHorizontal, 10)
	header.AddCSSClass("header-box")

	notifSwitch := gtk.NewSwitch()
	notifSwitch.SetActive(true)
	notifSwitch.AddCSSClass("custom-switch")
	notifSwitch.SetVAlign(gtk.AlignCenter)
	notifSwitch.SetHAlign(gtk.AlignStart)
	notifSwitch.SetTooltipText("開啟 / 關閉通知提醒與音效")

	titleLabel := gtk.NewLabel("歷史通知紀錄")
	titleLabel.AddCSSClass("header-title")
	titleLabel.SetHExpand(true)

	clearBtn := gtk.NewButtonWithLabel("清除全部")
	clearBtn.AddCSSClass("clear-button")

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

func CreateDateHeader(dateText string, listbox *gtk.ListBox, onDeleteDate func()) *gtk.Box {
	box := gtk.NewBox(gtk.OrientationVertical, 0)
	box.SetHExpand(true)
	box.AddCSSClass("date-header-box")
	box.AddCSSClass("is-header")

	sep := gtk.NewSeparator(gtk.OrientationHorizontal)
	sep.AddCSSClass("header-separator")
	sep.SetMarginStart(20)
	sep.SetMarginEnd(20)
	box.Append(sep)

	hbox := gtk.NewBox(gtk.OrientationHorizontal, 10)
	hbox.SetMarginStart(12)
	hbox.SetMarginEnd(12)
	hbox.SetMarginTop(10)
	hbox.SetHExpand(true)

	titleBox := gtk.NewBox(gtk.OrientationHorizontal, 10)
	titleBox.SetHExpand(true)

	arrow := gtk.NewImageFromIconName("pan-down-symbolic")
	arrow.SetMarginBottom(2)
	label := gtk.NewLabel(dateText)
	label.SetYAlign(0.5)
	label.AddCSSClass("date-header-label")

	titleBox.Append(arrow)
	titleBox.Append(label)
	hbox.Append(titleBox)

	if onDeleteDate != nil {
		deleteBtn := gtk.NewButtonFromIconName("user-trash-symbolic")
		deleteBtn.AddCSSClass("delete-button")
		deleteBtn.SetVAlign(gtk.AlignCenter)
		deleteBtn.SetHAlign(gtk.AlignEnd)
		deleteBtn.ConnectClicked(func() {
			onDeleteDate()
		})
		hbox.Append(deleteBtn)
	}

	box.Append(hbox)

	click := gtk.NewGestureClick()
	titleBox.AddController(click)

	isCollapsed := false
	click.ConnectReleased(func(n int, x, y float64) {
		isCollapsed = !isCollapsed

		if isCollapsed {
			arrow.SetFromIconName("pan-end-symbolic")
		} else {
			arrow.SetFromIconName("pan-down-symbolic")
		}

		children := listbox.ObserveChildren()
		targetDateClass := "date-" + dateText

		for i := uint(0); i < children.NItems(); i++ {
			item := children.Item(i)
			row, ok := item.Cast().(*gtk.ListBoxRow)
			if !ok {
				continue
			}

			if row.HasCSSClass(targetDateClass) && row.HasCSSClass("is-content") {
				if row.Header() != nil {
					if child := row.Child(); child != nil {
						if widget, ok := child.(*gtk.Widget); ok {
							widget.SetVisible(!isCollapsed)
						}
					}
					if isCollapsed {
						row.AddCSSClass("collapsed-header-row")
					} else {
						row.RemoveCSSClass("collapsed-header-row")
					}
					row.SetVisible(true)

				} else {
					row.SetVisible(!isCollapsed)
				}
			}

		}

		listbox.QueueAllocate()
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
