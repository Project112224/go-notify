# Testing & Quality Assurance Guide

## 1. 測試矩陣與工具說明

本專案建置了完整的五大測試體系：

| 測試類型 | 執行指令 | 涵蓋目標與說明 |
| :--- | :--- | :--- |
| **1. 單元測試 (Unit Test)** | `make test` | 測試 URL 提取、時間解析、模型、Repository 存取與 ViewModel 狀態邏輯 |
| **2. UI 測試 (UI Test)** | `make test` | 測試 GTK Label 標籤建構、Urgency CSS 類別綁定、App Icon 映射與 Linkify 超連結轉換 |
| **3. 競態條件檢測 (Race Detector)** | `make test-race` | 透過 Go `-race` 標籤測試 D-Bus 訊號併發、Channel 佇列與狀態競爭 |
| **4. 靜態分析與掃描 (Static Analysis)** | `make vet` | 使用 `go vet` 檢查變數格式化、型別轉換與潛在 runtime 錯誤 |
| **5. 模糊測試 (Fuzz Testing)** | `make fuzz` | 使用 Go 1.18+ 原生 Fuzzing 進行邊界字串隨機突變測試 |
| **6. 基準測試 (Benchmark)** | `make bench` | 使用 `go test -bench=. -benchmem` 評估記憶體配置量與執行耗時 (`ns/op`, `B/op`) |
| **7. 一鍵綜合檢查** | `make check-all` | 依序執行 `vet` -> `test` -> `test-race` -> `bench` 流水線 |

---

## 2. Makefile 測試指令說明

```bash
# 1. 執行基礎單元測試與 UI 測試
make test

# 2. 執行競態條件檢測 (Race Detector)
make test-race

# 3. 執行 Go 靜態程式碼檢測
make vet

# 4. 執行模糊測試 (可自訂 FUZZTIME，例如 30s)
make fuzz FUZZTIME=30s

# 5. 執行效能與記憶體 Benchmark
make bench

# 6. 一鍵綜合驗證流水線
make check-all
```

---

## 3. 測試檔案列表

- `go-notify/internal/viewmodel/notification_viewmodel_test.go`
- `go-notify/internal/service/launcher_test.go`
- `go-notify/internal/service/launcher_fuzz_test.go`
- `go-notify/internal/service/launcher_bench_test.go`
- `go-notify/internal/repository/history_repository_test.go`
- `go-notify/internal/dbus/server_test.go`
- `go-notify/internal/ui/card_ui_test.go`
- `go-notify-manager/internal/viewmodel/history_viewmodel_test.go`
- `go-notify-manager/internal/extension/string_helper_test.go`
- `go-notify-manager/internal/extension/string_helper_fuzz_test.go`
- `go-notify-manager/internal/extension/string_helper_bench_test.go`
- `go-notify-manager/internal/repository/history_repository_test.go`
- `go-notify-manager/internal/ui/row_ui_test.go`
- `go-notify-manager/internal/ui/row_fuzz_test.go`
- `go-notify-manager/internal/ui/row_bench_test.go`
