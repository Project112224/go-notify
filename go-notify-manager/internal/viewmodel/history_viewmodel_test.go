package viewmodel

import (
	"testing"

	"go-notify-manager/internal/models"
)

type MockHistoryRepository struct {
	items []models.HistoryItem
}

func (m *MockHistoryRepository) LoadAll() ([]models.HistoryItem, error) {
	return m.items, nil
}

func (m *MockHistoryRepository) ClearAll() error {
	m.items = nil
	return nil
}

func (m *MockHistoryRepository) DeleteOne(id int) error {
	var filtered []models.HistoryItem
	for _, item := range m.items {
		if item.ID != id {
			filtered = append(filtered, item)
		}
	}
	m.items = filtered
	return nil
}

func (m *MockHistoryRepository) DeleteByDate(dateStr string) error {
	return nil
}

func (m *MockHistoryRepository) Close() error {
	return nil
}

func TestHistoryViewModel_Filter(t *testing.T) {
	mockRepo := &MockHistoryRepository{
		items: []models.HistoryItem{
			{ID: 1, AppName: "Firefox", Summary: "Download complete", Body: "file.zip"},
			{ID: 2, AppName: "Discord", Summary: "New message", Body: "Hello world"},
		},
	}

	vm := NewHistoryViewModel(mockRepo, nil)
	vm.LoadHistory()

	if len(vm.FilteredItems()) != 2 {
		t.Errorf("Expected 2 items, got %d", len(vm.FilteredItems()))
	}

	vm.SetSearchQuery("firefox")
	filtered := vm.FilteredItems()
	if len(filtered) != 1 || filtered[0].ID != 1 {
		t.Errorf("Expected 1 item (Firefox), got %v", filtered)
	}
}

func TestHistoryViewModel_DeleteSingleItem(t *testing.T) {
	mockRepo := &MockHistoryRepository{
		items: []models.HistoryItem{
			{ID: 1, AppName: "Firefox", Summary: "Download complete"},
			{ID: 2, AppName: "Discord", Summary: "New message"},
		},
	}

	vm := NewHistoryViewModel(mockRepo, nil)
	vm.LoadHistory()

	vm.DeleteSingleItem(1)

	if len(vm.FilteredItems()) != 1 || vm.FilteredItems()[0].ID != 2 {
		t.Errorf("Expected 1 item remaining (ID: 2), got %v", vm.FilteredItems())
	}
}
