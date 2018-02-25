package models

import (
	"time"
)

// SystemSettings is the system-wide settings
type SystemSettings struct {
	Mode            string `json:"mode"`
	BuildTime       string `json:"buildTime"`
	Version         string `json:"version"`
	Environment     string `json:"environment"`
	GoogleAnalytics string `json:"googleAnalytics"`
	Compiler        string `json:"compiler"`
	Domain          string `json:"domain"`
}

// Notification is the system generated notification entity
type Notification struct {
	ID        int       `json:"id" db:"id"`
	Title     string    `json:"title" db:"title"`
	Link      string    `json:"link" db:"link"`
	Read      bool      `json:"read" db:"read"`
	CreatedOn time.Time `json:"created_on" db:"created_on"`
}
