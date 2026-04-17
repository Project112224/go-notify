package models

import "time"

type HistoryNotif struct {
	ID      uint32    `json:"id"`
	AppName string    `json:"app_name"`
	Summary string    `json:"summary"`
	Body    string    `json:"body"`
	Urgency int       `json:"urgency"`
	Time    time.Time `json:"time"`
	Icon    string    `json:"icon"`
}
