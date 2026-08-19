package dbus

import (
	"testing"

	"github.com/godbus/dbus/v5"
)

func TestGetUrgency(t *testing.T) {
	hints := map[string]dbus.Variant{
		"urgency": dbus.MakeVariant(uint8(2)),
	}
	if urgency := getUrgency(hints); urgency != 2 {
		t.Errorf("Expected urgency 2, got %d", urgency)
	}

	emptyHints := map[string]dbus.Variant{}
	if urgency := getUrgency(emptyHints); urgency != 1 {
		t.Errorf("Expected default urgency 1, got %d", urgency)
	}
}

func TestGetDesktopEntry(t *testing.T) {
	hints := map[string]dbus.Variant{
		"desktop-entry": dbus.MakeVariant("firefox"),
	}
	if entry := getDesktopEntry(hints); entry != "firefox" {
		t.Errorf("Expected desktop entry 'firefox', got '%s'", entry)
	}
}

func TestGetDefaultActionKey(t *testing.T) {
	actions := []string{"default", "Default Action", "cancel", "Cancel"}
	if key := getDefaultActionKey(actions); key != "default" {
		t.Errorf("Expected action key 'default', got '%s'", key)
	}

	actions2 := []string{"open", "Open"}
	if key := getDefaultActionKey(actions2); key != "open" {
		t.Errorf("Expected action key 'open', got '%s'", key)
	}

	if key := getDefaultActionKey(nil); key != "" {
		t.Errorf("Expected empty action key for nil actions, got '%s'", key)
	}
}

func TestGetAppIcon(t *testing.T) {
	hints := map[string]dbus.Variant{
		"image-path": dbus.MakeVariant("file:///usr/share/pixmaps/app.png"),
	}
	icon := getAppIcon(hints, "default-icon")
	if icon != "/usr/share/pixmaps/app.png" {
		t.Errorf("Expected '/usr/share/pixmaps/app.png', got '%s'", icon)
	}
}

func TestVariantToString(t *testing.T) {
	vStr := dbus.MakeVariant("hello")
	if res := variantToString(vStr); res != "hello" {
		t.Errorf("Expected 'hello', got '%s'", res)
	}

	vBytes := dbus.MakeVariant([]byte("byte-data"))
	if res := variantToString(vBytes); res != "byte-data" {
		t.Errorf("Expected 'byte-data', got '%s'", res)
	}
}
