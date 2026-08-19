package ui

import (
	"testing"

	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

func TestSetUrgencyClass(t *testing.T) {
	gtk.Init()

	cardLow := gtk.NewBox(gtk.OrientationHorizontal, 0)
	setUrgencyClass(cardLow, 0)
	if !cardLow.HasCSSClass("urgency-low") {
		t.Errorf("Expected card to have CSS class 'urgency-low'")
	}

	cardCritical := gtk.NewBox(gtk.OrientationHorizontal, 0)
	setUrgencyClass(cardCritical, 2)
	if !cardCritical.HasCSSClass("urgency-critical") {
		t.Errorf("Expected card to have CSS class 'urgency-critical'")
	}

	cardDefault := gtk.NewBox(gtk.OrientationHorizontal, 0)
	setUrgencyClass(cardDefault, 1)
	if !cardDefault.HasCSSClass("notif-card") {
		t.Errorf("Expected card to have CSS class 'notif-card'")
	}
}

func TestCreateTextContainer(t *testing.T) {
	gtk.Init()

	summary := "Test Summary & Special <Char>"
	body := "Test Body Content"

	container := createTextContainer(summary, body)
	if container == nil {
		t.Fatalf("Expected text container widget, got nil")
	}

	if container.FirstChild() == nil {
		t.Fatalf("Expected child label in text container")
	}
}
