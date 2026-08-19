package ui

import (
	"testing"

	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

func TestResolveAppIconName(t *testing.T) {
	tests := []struct {
		appName  string
		expected string
	}{
		{"Discord", "discord"},
		{"Firefox Developer Edition", "firefox"},
		{"Google-Chrome", "google-chrome"},
		{"TelegramDesktop", "telegram"},
		{"Spotify Premium", "spotify"},
		{"VSCode", "com.visualstudio.code"},
		{"Alacritty", "utilities-terminal"},
		{"MyCustomApp", "mycustomapp"},
		{"", "dialog-information"},
	}

	for _, tt := range tests {
		res := resolveAppIconName(tt.appName)
		if res != tt.expected {
			t.Errorf("resolveAppIconName(%q) = %q; want %q", tt.appName, res, tt.expected)
		}
	}
}

func TestLinkify(t *testing.T) {
	input := "Check https://example.com for info"
	expected := `Check <a href="https://example.com">https://example.com</a> for info`
	if res := linkify(input); res != expected {
		t.Errorf("linkify(%q) = %q; want %q", input, res, expected)
	}
}

func TestNewHistoryRowUI(t *testing.T) {
	app := gtk.NewApplication("com.github.test.row", 0)
	app.ConnectActivate(func() {
		defer app.Quit()

		row := NewHistoryRow("Firefox", "Download Done", "file.zip", "firefox", 1, "2026-08-19T15:30:00Z", func() {})
		if row == nil {
			t.Fatalf("Expected ListBoxRow widget, got nil")
		}

		if !row.HasCSSClass("history-row") {
			t.Errorf("Expected row to have CSS class 'history-row'")
		}
	})
	app.Run(nil)
}
