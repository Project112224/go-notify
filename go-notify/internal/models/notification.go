package models

type Notification struct {
	ID               uint32
	AppName          string
	Icon             string
	Summary          string
	Body             string
	Urgency          byte
	DesktopEntry     string
	DefaultActionKey string
}
