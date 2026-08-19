package viewmodel

import (
	"testing"

	"go-notify/internal/models"
)

type MockLauncherService struct {
	openedURL string
	openedApp string
}

func (m *MockLauncherService) OpenURL(rawURL string) {
	m.openedURL = rawURL
}

func (m *MockLauncherService) OpenApplication(appName string, desktopEntry string) {
	m.openedApp = appName
}

func (m *MockLauncherService) ExtractURLs(text string) []string {
	if text == "Check https://example.com" {
		return []string{"https://example.com"}
	}
	return nil
}

func TestHandleCardClick_LeftClickURL(t *testing.T) {
	mockLauncher := &MockLauncherService{}
	vm := NewNotificationViewModel(mockLauncher)

	notif := &models.Notification{
		ID:      1,
		AppName: "TestApp",
		Body:    "Check https://example.com",
	}

	vm.HandleCardClick(notif, 1)

	if mockLauncher.openedURL != "https://example.com" {
		t.Errorf("Expected openedURL to be 'https://example.com', got '%s'", mockLauncher.openedURL)
	}
}

func TestHandleCardClick_LeftClickApp(t *testing.T) {
	mockLauncher := &MockLauncherService{}
	vm := NewNotificationViewModel(mockLauncher)

	notif := &models.Notification{
		ID:      1,
		AppName: "TestApp",
		Body:    "No url here",
	}

	vm.HandleCardClick(notif, 1)

	if mockLauncher.openedApp != "TestApp" {
		t.Errorf("Expected openedApp to be 'TestApp', got '%s'", mockLauncher.openedApp)
	}
}

func TestHandleCardClick_RightClick(t *testing.T) {
	mockLauncher := &MockLauncherService{}
	vm := NewNotificationViewModel(mockLauncher)

	notif := &models.Notification{
		ID:      1,
		AppName: "TestApp",
		Body:    "Check https://example.com",
	}

	// Right click (button 3) should not open app or URL
	vm.HandleCardClick(notif, 3)

	if mockLauncher.openedURL != "" || mockLauncher.openedApp != "" {
		t.Errorf("Expected no launcher call on right click, got url='%s', app='%s'", mockLauncher.openedURL, mockLauncher.openedApp)
	}
}
