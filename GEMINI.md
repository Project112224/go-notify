# Workspace Rules & Coding Style Guide

本文件為 `go-notify` 與 `go-notify-manager` 專案的團隊開發規範與程式碼風格 (Coding Style Guidelines)。

---

## 1. 專案架構與路徑規範 (Architecture & Path Conventions)

本 repository 由兩個核心 Go 專案組成：

```
/
├── go-notify/                # 推播服務 Daemon
│   ├── main.go               # 依賴注入與進入點
│   └── internal/
│       ├── models/           # 通知與歷史紀錄結構體 (Notification, HistoryNotif)
│       ├── repository/       # SQLite 存取介面與實作 (HistoryRepository)
│       ├── service/          # 應用程式與 URL 開啟服務 (LauncherService)
│       ├── dbus/             # D-Bus Server 實現 (org.freedesktop.Notifications)
│       ├── viewmodel/        # NotificationViewModel (點擊與動作處理)
│       └── ui/               # GTK4 LayerShell 視窗與卡片 (NotifWindow, NotifCard)
└── go-notify-manager/        # 通知歷史管理器
    ├── main.go               # 依賴注入與進入點
    └── internal/
        ├── extension/        # 日期與字串擴充輔助
        ├── models/           # 歷史紀錄模型 (HistoryItem)
        ├── repository/       # SQLite 歷史查詢與刪除介面 (HistoryRepository)
        ├── services/         # FocusMode D-Bus 與系統靜音服務 (FocusModeService)
        ├── viewmodel/        # HistoryViewModel (搜尋過濾與刪除狀態)
        └── ui/               # GTK4 視窗、ListBox 列表、Header 與快捷鍵
```

---

## 2. 程式碼風格 (Coding Style Guidelines)

### 2.1 MVVM 分層職責與解耦

1. **`models`**: 僅包含純粹的 Go Struct，不包含任何 DB 操作或 GTK UI 邏輯。
2. **`repository`**: 所有資料庫存取必須定義介面 (Interface)，實作層採用參數化查詢（防範 SQL 注入），禁止直接在 UI 或 ViewModel 中編寫原生 SQL。
3. **`service` / `services`**: 封裝系統外設呼叫（如 `.desktop` 解析、`xdg-open`、`pactl`）與 D-Bus 請求。
4. **`viewmodel`**: 管理狀態、即時搜尋過濾與發送 D-Bus 訊號。UI 與 ViewModel 之間透過 Callback 或方法呼叫溝通。
5. **`ui`**: 純粹的 GTK4 視窗與元件渲染層。**禁止直接在 UI 元件內操作 SQLite 或發送 D-Bus 請求**。

---

### 2.2 GTK4 與手勢控制器規範 (GTK4 & Event Controllers)

1. **`PhaseCapture` 事件捕獲階段**:
   - 在容器元件（如 `card *gtk.Box`）綁定手勢控制器（`GtkGestureClick`, `GtkGestureSwipe`, `EventControllerScroll`）時，必須呼叫 `.SetPropagationPhase(gtk.PhaseCapture)`，避免內層 `GtkLabel` 或 `GtkBox` 氣泡攔截點擊事件。
2. **手勢按鍵區分**:
   - `button == 1`（單指/左鍵）：觸發預設動作（開啟 URL 或應用程式）、發送 D-Bus `ActionInvoked` 並關閉卡片。
   - `button == 3`（雙指點擊/右鍵）或 `button == 2`（中鍵）：發送 D-Bus `NotificationClosed(reason=2)` 並關閉卡片（不安裝/不開啟應用程式）。

---

### 2.3 錯誤處理與日誌記錄 (Error Handling & Logging)

1. **Prefix 日誌格式**: 日誌統一包含模組前綴，如 `[NotifViewModel]`, `[HistoryRepository]`, `[LauncherService]`。
2. **安全逃逸**: GTK Label 使用 `glib.MarkupEscapeText()` 逃逸可能包含 HTML 標籤的標題與內文。
3. **路徑安全**: 檔案存取統一使用 `filepath.Join`，避免硬編碼 `/` 或 `\`。

---

## 3. 自動化測試與驗證規範 (Testing Rules)

1. **單元與 UI 測試**: 凡新增/重構 ViewModel, Repository, Service 或 UI 元件，必須編寫對應之 `*_test.go` 與 `*_ui_test.go`。
2. **執行完整測試套件**:
   - 變更後必須執行 `make test` 確保單元與 UI 測試 100% 通過。
   - 執行 `make vet` 確保程式碼通過 Go 靜態分析。
   - 執行 `make test-race` 確保核心邏輯無競爭條件。
   - 執行 `make bench` 評估記憶體配置量與效能變化。
3. **系統執行檔更新**: 修改後需執行 `make install` 確保更新 `~/.local/bin/go-notify` 與 `~/.local/bin/go-notify-manager`。
