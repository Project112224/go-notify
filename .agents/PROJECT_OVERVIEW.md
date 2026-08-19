# Project Overview & Architecture Guide

## 1. 專案概要 (Overview)

本 repository 包含兩個獨立編譯但相互協作的 Go / GTK4 應用程式：

1. **`go-notify`**: 遵循 D-Bus `org.freedesktop.Notifications` 規範的 Linux 桌面通知 Daemon。採用 GTK4 Layer-Shell 浮動於 Wayland 視窗層次最上層，提供動態通知卡片渲染、自動關閉倒數、滑鼠/觸控板手勢（單指開啟連結/應用程式、雙指/右鍵關閉卡片）及 SQLite 歷史紀錄持久化。
2. **`go-notify-manager`**: 通知歷史紀錄桌面端 GUI 管理器。提供關鍵字搜尋過濾、按日期自動分組與折疊/展開、單筆與按日期批量刪除歷史、一鍵清除全部紀錄，以及與 `go-notify` 通訊的 Focus Mode（專注模式與系統音效靜音切換）。

---

## 2. 系統架構 (MVVM Architecture)

本專案採用 **MVVM (Model - View - ViewModel)** 與 **Clean Architecture 服務分層**：

```
                              ┌────────────────────────────────┐
                              │           UI Views             │
                              │ (NotifWindow / ManagerWindow)  │
                              └───────────────┬────────────────┘
                                              │
                                              ▼
                              ┌────────────────────────────────┐
                              │          ViewModels            │
                              │ (NotificationVM / HistoryVM)   │
                              └───────────────┬────────────────┘
                                              │
                                ┌─────────────┴─────────────┐
                                ▼                           ▼
                      ┌───────────────────┐       ┌───────────────────┐
                      │    Repositories   │       │     Services      │
                      │ (SQLite DB Layer) │       │ (Launcher/DBus)   │
                      └───────────────────┘       └───────────────────┘
```

### 2.1 模組角色職責劃分

- **`models`**:
  - `models.Notification`: 通知卡片記憶體資料模型。
  - `models.HistoryNotif`: SQLite 歷史紀錄資料模型。
  - `models.HistoryItem`: Manager 管理器對映紀錄結構體。
- **`repository`**:
  - `HistoryRepository` (Interface): 定義 `Save`, `LoadAll`, `DeleteOne`, `DeleteByDate`, `ClearAll`, `Close` 存取 API。
  - `SQLiteHistoryRepository`: 基於 `database/sql` + `go-sqlite3` WAL 模式實現。
- **`service` / `services`**:
  - `LauncherService`: 處理外設開啟 URL (xdg-open, firefox, gio) 及搜尋 `/usr/share/applications/*.desktop` 喚醒/啟動 App。
  - `FocusModeService`: 透過 D-Bus (`org.freedesktop.Notifications.SetFocusMode`) 傳遞專注模式狀態，並執行系統音效靜音控制 (`pactl set-sink-mute`)。
  - `DBusServer`: 匯出 D-Bus 物件，接收 `Notify` 請求並發送 `ActionInvoked` / `NotificationClosed` 訊號。
- **`viewmodel`**:
  - `NotificationViewModel`: 處理通知卡片點擊按鍵邏輯（Button 1 開啟 URL/App + 發送 ActionInvoked；Button 3 / 其他按鍵僅發送 NotificationClosed 關閉卡片）。
  - `HistoryViewModel`: 管理歷史紀錄載入、即時搜尋過濾、按日期/單筆刪除與專注模式開關狀態同步。
- **`ui`**:
  - 純 UI 視窗與元件渲染層（GTK4 ApplicationWindow, ListBox, Box, Label, GestureController），不直接存取 SQLite 或 D-Bus，全數透過 ViewModel 雙向綁定。

---

## 3. 目錄結構 (Directory Layout)

```
/home/june/Work/go-notify/
├── Makefile                          # 全域編譯、測試與安裝腳本
├── README.md                         # 專案使用者與開發說明文件
├── .agents/                          # AI Agent 專屬架構、測試與資安文件
│   ├── PROJECT_OVERVIEW.md
│   ├── SECURITY_AUDIT.md
│   ├── TESTING_GUIDE.md
│   └── GESTURE_INTERACTION_NOTES.md
├── go-notify/                        # 推播服務 Daemon
│   ├── main.go                       # 依賴注入與進入點
│   └── internal/
│       ├── models/                   # 資料模型
│       ├── repository/               # DB 操作介面與 SQLite 實作
│       ├── service/                  # Launcher 服務
│       ├── dbus/                     # D-Bus Server 實現
│       ├── viewmodel/                # NotificationViewModel & 測試
│       ├── ui/                       # GTK4 LayerShell 視窗 & 卡片 UI
│       └── util/                     # Logger 與進程工具
└── go-notify-manager/                # 通知歷史管理器
    ├── main.go                       # 依賴注入與進入點
    └── internal/
        ├── extension/                # 時間格式化擴充
        ├── logger/                   # 日誌工具
        ├── models/                   # 歷史紀錄模型
        ├── repository/               # DB 歷史查詢介面與 SQLite 實作
        ├── services/                 # FocusMode D-Bus 服務
        ├── viewmodel/                # HistoryViewModel & 測試
        └── ui/                       # GTK4 Manager 視窗、Row、Header、快捷鍵
```
