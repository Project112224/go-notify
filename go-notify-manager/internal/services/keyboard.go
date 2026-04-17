package service

import (
	"log"

	"go-notify-manager/internal/ui"

	"github.com/diamondburned/gotk4/pkg/gdk/v4"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

type KeyboardService struct {
	win *ui.Window
}

func NewKeyboardService(win *ui.Window) *KeyboardService {
	return &KeyboardService{win: win}
}

func (s *KeyboardService) BindShortcuts() {
	keyCtrl := gtk.NewEventControllerKey()
	s.win.Window.AddController(keyCtrl)

	keyCtrl.ConnectKeyPressed(func(keyval uint, keycode uint, state gdk.ModifierType) bool {
		// 1. 處理 Q 關閉 (排除搜尋框輸入)
		if keyval == gdk.KEY_q && !s.win.SearchEntry.IsFocus() {
			s.win.Window.Close()
			return true
		}

		// 2. 處理 Ctrl + C 複製
		if keyval == gdk.KEY_c && state.Has(gdk.ControlMask) {
			return s.handleCopy()
		}

		return false
	})
}

// 私有方法處理複製邏輯
func (s *KeyboardService) handleCopy() bool {
	focusWidget := s.win.Window.Focus()
	if focusWidget == nil {
		return false
	}

	if label, ok := focusWidget.(*gtk.Label); ok {
		start, end, hasSelection := label.SelectionBounds()
		if hasSelection {
			fullText := label.Text()
			selectedText := fullText[start:end]
			s.win.Window.Clipboard().SetText(selectedText)
			log.Println("已複製選取文字:", selectedText)
			return true
		}
	}

	return false
}
