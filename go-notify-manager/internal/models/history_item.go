package models

type HistoryItem struct {
	ID        int    `json:"id"`
	AppName   string `json:"app_name"`
	Summary   string `json:"summary"`
	Body      string `json:"body"`
	Urgency   int    `json:"urgency"`
	IconPath  string `json:"icon_path"`
	CreatedAt string `json:"created_at"`
}
