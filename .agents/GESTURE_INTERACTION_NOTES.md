# Touchpad & Gesture Interaction Implementation Notes

## 1. 觸控板與手勢交互機制 (Gesture Handling)

在 GTK4 Wayland / Hyprland 環境下，使用者點擊或操作觸控板通知卡片時，libinput / GTK 可能會發送多種不同的事件類型。為確保手勢與點擊 100% 成功觸發，[card.go](file:///home/june/Work/go-notify/go-notify/internal/ui/card.go#L93-L127) 實現了以下機制：

### 1.1 `gtk.PhaseCapture` 搶先捕獲階段

- **問題**: GTK4 預設採用氣泡階段 (`PhaseBubble`)，當點擊發生在 `card` (`GtkBox`) 內部的文字標籤 (`GtkLabel`) 或圖示 (`GtkImage`) 時，子元件可能會先消化點擊事件，導致父容器的手勢控制器收不到通知。
- **解決方案**: 每個手勢控制器皆呼叫 `SetPropagationPhase(gtk.PhaseCapture)`，在事件向下傳遞給子元件前**搶先捕獲**。

---

### 1.2 三重手勢控制器組合 (Triple Controller Stack)

```go
// 1. 滑鼠 / 觸控板點擊 (Button 1: 左鍵開啟 App/URL, Button 3: 右鍵/雙指關閉卡片)
click := gtk.NewGestureClick()
click.SetButton(0) // 0 表示監聽所有按鍵
click.SetPropagationPhase(gtk.PhaseCapture)

// 2. 觸控板雙指撥動 / 滑動 (Swipe)
swipe := gtk.NewGestureSwipe()
swipe.SetButton(0)
swipe.SetPropagationPhase(gtk.PhaseCapture)

// 3. 觸控板雙指滾動 (Scroll)
scroll := gtk.NewEventControllerScroll(gtk.EventControllerScrollBothAxes)
scroll.SetPropagationPhase(gtk.PhaseCapture)
```

1. **`GtkGestureClick`**: 聽取所有滑鼠與觸控板點擊按鍵。
   - `button == 1`（單指點擊/左鍵）：觸發 `ActionInvoked` D-Bus 訊號，自動開啟內文 URL 或喚醒應用程式，並關閉卡片。
   - `button == 3`（雙指點擊/右鍵）或 `button == 2`（中鍵）：發送 `NotificationClosed` (Reason 2: Dismissed) D-Bus 訊號並關閉卡片（不安裝/不開啟應用程式）。
2. **`GtkGestureSwipe`**: 捕獲觸控板雙指滑動/撥動 (Swipe) 手勢，觸發卡片關閉。
3. **`EventControllerScroll`**: 捕獲觸控板雙指滾動 (Scroll) 手勢，觸發卡片關閉。

---

## 2. D-Bus 訊號規範 (Desktop Notifications Spec)

遵循 FreeDesktop Notification 規範：
- **`ActionInvoked`**: 當使用者點擊卡片觸發預設動作（單指/左鍵）時發送。
- **`NotificationClosed(id, reason)`**:
  - `reason = 2`: User dismissed the notification (使用者主動點擊或雙指關閉卡片)。
