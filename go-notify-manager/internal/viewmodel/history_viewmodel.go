package viewmodel

import (
	"log"
	"strings"

	"go-notify-manager/internal/models"
	"go-notify-manager/internal/repository"
	"go-notify-manager/internal/services"
)

type HistoryViewModel struct {
	repo           repository.HistoryRepository
	focusSvc       services.FocusModeService
	items          []models.HistoryItem
	SearchQuery    string
	IsNotifEnabled bool

	OnItemsUpdated     func(items []models.HistoryItem)
	OnFocusModeUpdated func(enabled bool)
}

func NewHistoryViewModel(repo repository.HistoryRepository, focusSvc services.FocusModeService) *HistoryViewModel {
	return &HistoryViewModel{
		repo:     repo,
		focusSvc: focusSvc,
	}
}

func (vm *HistoryViewModel) LoadHistory() {
	items, err := vm.repo.LoadAll()
	if err != nil {
		log.Printf("[HistoryViewModel] 載入歷史紀錄失敗: %v", err)
		return
	}
	vm.items = items
	vm.notifyItemsUpdated()
}

func (vm *HistoryViewModel) MarkAllAsRead() {
	if err := vm.repo.MarkAllAsRead(); err != nil {
		log.Printf("[HistoryViewModel] 標記已讀失敗: %v", err)
	}
	vm.LoadHistory()
}

func (vm *HistoryViewModel) SetSearchQuery(query string) {
	vm.SearchQuery = strings.ToLower(strings.TrimSpace(query))
	vm.notifyItemsUpdated()
}

func (vm *HistoryViewModel) FilteredItems() []models.HistoryItem {
	if vm.SearchQuery == "" {
		return vm.items
	}

	var filtered []models.HistoryItem
	for _, item := range vm.items {
		content := strings.ToLower(item.AppName + " " + item.Summary + " " + item.Body)
		if strings.Contains(content, vm.SearchQuery) {
			filtered = append(filtered, item)
		}
	}
	return filtered
}

func (vm *HistoryViewModel) DeleteSingleItem(id int) {
	if err := vm.repo.DeleteOne(id); err != nil {
		log.Printf("[HistoryViewModel] 刪除項目 %d 失敗: %v", id, err)
		return
	}
	vm.LoadHistory()
}

func (vm *HistoryViewModel) DeleteDateGroup(dateStr string) {
	if err := vm.repo.DeleteByDate(dateStr); err != nil {
		log.Printf("[HistoryViewModel] 刪除日期群組 %s 失敗: %v", dateStr, err)
		return
	}
	vm.LoadHistory()
}

func (vm *HistoryViewModel) ClearAll() {
	if err := vm.repo.ClearAll(); err != nil {
		log.Printf("[HistoryViewModel] 清除所有歷史紀錄失敗: %v", err)
		return
	}
	vm.LoadHistory()
}

func (vm *HistoryViewModel) InitFocusMode() {
	if vm.focusSvc == nil {
		return
	}
	isLocked, err := vm.focusSvc.GetFocusMode()
	if err != nil {
		log.Printf("[HistoryViewModel] 取得 FocusMode 失敗: %v", err)
		vm.IsNotifEnabled = false
	} else {
		vm.IsNotifEnabled = !isLocked
	}
	if vm.OnFocusModeUpdated != nil {
		vm.OnFocusModeUpdated(vm.IsNotifEnabled)
	}
}

func (vm *HistoryViewModel) ToggleNotificationStatus(enabled bool) {
	vm.IsNotifEnabled = enabled
	if vm.focusSvc != nil {
		if err := vm.focusSvc.SetFocusMode(enabled); err != nil {
			log.Printf("[HistoryViewModel] 設定通知狀態失敗: %v", err)
		}
	}
	if vm.OnFocusModeUpdated != nil {
		vm.OnFocusModeUpdated(vm.IsNotifEnabled)
	}
}

func (vm *HistoryViewModel) notifyItemsUpdated() {
	if vm.OnItemsUpdated != nil {
		vm.OnItemsUpdated(vm.FilteredItems())
	}
}
