# go-notify & go-notify-manager

輕量級 Linux 桌面通知推播服務與通知歷史管理器，基於 Go、GTK4、Wayland LayerShell 及 D-Bus 構建。

---

## 📌 專案簡介

本專案包含兩個核心元件：
1. **`go-notify`**: 遵循 D-Bus `org.freedesktop.Notifications` 規範的桌面通知 Daemon，使用 GTK4 Layer-Shell 在 Wayland 畫面上即時彈出 notification toast，支援手勢控制（單指/左鍵點擊開啟連結與 App、雙指/右鍵關閉通知卡片）、定時自動關閉與 SQLite 歷史紀錄存儲。
2. **`go-notify-manager`**: 桌面端通知歷史管理 GUI 工具，提供歷史通知列表瀏覽、即時搜尋過濾、依日期分組折疊與刪除、一鍵清除歷史，以及 Focus Mode（專注模式 / 音效靜音切換）控制。

---

## 🏗️ 架構設計 (MVVM & Clean Architecture)

專案全面採用 **MVVM (Model-View-ViewModel)** 與服務分層架構設計，將 UI 佈局、狀態邏輯、D-Bus 通訊、SQLite 資料庫與系統命令解耦：

```
go-notify/
├── main.go                     # go-notify 進入點（元件注入與啟動）
├── Makefile                    # 編譯與安裝腳本
├── README.md                   # 說明文件
├── go-notify/                  # 通知推播 Daemon 專案
│   ├── main.go
│   └── internal/
│       ├── models/             # 資料模型 (Notification, HistoryNotif)
│       ├── repository/         # SQLite 資料庫操作介面與實作 (HistoryRepository)
│       ├── service/            # 系統服務層 (LauncherService 開啟 URL/App)
│       ├── dbus/               # D-Bus 服務實現 (org.freedesktop.Notifications)
│       ├── viewmodel/          # NotificationViewModel 狀態與手勢事件處理
│       └── ui/                 # GTK4 LayerShell 視窗與卡片渲染 (NotifWindow, NotifCard)
└── go-notify-manager/          # 通知歷史管理器專案
    ├── main.go
    └── internal/
        ├── models/             # 資料模型 (HistoryItem)
        ├── repository/         # SQLite 歷史查詢與刪除 (HistoryRepository)
        ├── services/           # FocusMode D-Bus 與系統靜音服務 (FocusModeService)
        ├── viewmodel/          # HistoryViewModel 搜尋過濾、分組與刪除邏輯
        └── ui/                 # GTK4 視窗、ListBox 列表、快捷鍵與 Header (ManagerWindow)
```

### 重構架構圖

```mermaid
graph TD
    subgraph UI Layer (Views)
        NotifWin["NotifWindow / NotifCard"]
        MgrWin["ManagerWindow / HistoryRow"]
    end

    subgraph ViewModel Layer
        NotifVM["NotificationViewModel"]
        MgrVM["HistoryViewModel"]
    end

    subgraph Service & Repository Layer
        Repo["SQLiteHistoryRepository"]
        Launcher["DesktopLauncherService"]
        FocusSvc["DBusFocusModeService"]
        DBusServer["DBus NotificationServer"]
    end

    NotifWin --> NotifVM
    MgrWin --> MgrVM

    NotifVM --> Launcher
    NotifVM --> DBusServer
    NotifVM --> Repo

    MgrVM --> Repo
    MgrVM --> FocusSvc
```

---

## ✨ 核心功能特色

- **手勢與點擊互動**:
  - **單指點擊 / 左鍵**: 觸發預設動作、自動解析內文中的 URL 並以瀏覽器開啟；無 URL 時喚醒或啟動對應應用程式。
  - **雙指點擊 / 右鍵**: 關閉 notification card，發送 D-Bus `NotificationClosed` (Dismissed) 訊號。
- **自動關閉與 Hover 暫停**: 滑鼠移入卡片時暫停倒數，移出後自動恢復 5 秒倒數並關閉。
- **歷史紀錄管理**:
  - 自動寫入 SQLite 資料庫 (WAL 模式)。
  - 即時關鍵字搜尋（標題、內文、應用程式名稱）。
  - 按日期分組顯示，支援單筆刪除、按日期批量刪除與一鍵全清。
- **Focus Mode & 系統靜音**:
  - 提供專注模式開關，啟用時可切換系統音效靜音 (`pactl set-sink-mute`)。

---

## 🛠️ 編譯與安裝

### 1. 編譯專案

在根目錄執行：

```bash
make build
```

### 2. 安裝執行檔至系統

執行以下命令可將編譯後的 `go-notify` 與 `go-notify-manager` 安裝至 `~/.local/bin`：

```bash
make install
```

> **提示**: 請確保 `~/.local/bin` 已加入您的 `PATH` 環境變數中。

### 3. 自動化測試與程式碼品質檢查

本專案提供完整的單元測試、UI 測試、競態檢測、靜態分析、模糊測試與效能 Benchmarks：

```bash
# 執行單元與 UI 測試
make test

# 執行靜態程式碼分析 (go vet)
make vet

# 執行競態條件檢測 (Race Detector)
make test-race

# 執行效能與記憶體配置基準測試 (Benchmark & Memory Allocs)
make bench

# 執行模糊測試 (Fuzz Testing，預設 FUZZTIME=10s)
make fuzz FUZZTIME=10s

# 一鍵執行所有品質檢查與測試 (vet, test, test-race, bench)
make check-all
```

### 4. 清理編譯產物

```bash
make clean
```
