# Security Audit & Privacy Report

## 1. 私密資料與敏感資訊檢查 (Privacy Audit)

經掃描與稽核整個程式庫：

- **密碼 / 金鑰 / Token 檢查**: 專案無任何寫死之 API 密鑰、密碼、Token、Private SSH Key 或憑證。
- **SQLite 資料庫隱私**: 
  - 資料庫檔案儲存於使用者家目錄之 `~/.local/share/go-notify/history.db`。
  - Git 倉庫包含 [.gitignore](file:///home/june/Work/go-notify/.gitignore)，明確排除 `*.db`, `*.db-journal`, `*.db-wal`, `*.db-shm` 與 `history.db`，確保真實通知數據不會上傳至 Version Control。
- **日誌與暫存檔**: 
  - 日誌檔記錄於 `~/.local/share/go-notify/*.log`，`.gitignore` 已排除 `*.log`。
  - D-Bus 接收之圖片暫存檔寫入 `/tmp/notif_img_*.png`，無個人敏感隱私殘留。
- **程式碼硬編碼資訊**: 無任何個人私人 Server 網址、Email 或敏感私密 IP。

---

## 2. 邊界防禦與安全性 (Security Measures)

- **SQL 注入防禦**: SQLite 查詢全數採用參數化查詢 (`db.Query(query, val)`, `db.Exec(query, val)`)，避免 SQL Injection 風險。
- **XSS / HTML Markup 逃逸**: GTK Label 使用 `glib.MarkupEscapeText(summary)` 逃逸 D-Bus 送入之 HTML 特殊字元 (`<`, `>`, `&`)，防止 Markup 解析錯誤與無意攻擊。
- **執行檔清理**: `.desktop` 的 `Exec=` 欄位指令透過 `cleanExecLine` 嚴格過濾 `%u`, `%F`, `%U`, `%i` 等預留佔位符，避免指令注入漏洞。
