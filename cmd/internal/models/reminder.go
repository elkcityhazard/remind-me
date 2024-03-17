package models

import (
	"html/template"
	"time"
)

type Reminderer interface {
	InsertReminder(*Reminder) (int64, error)
}

type Reminder struct {
	ID          int64  `json:"id"`
	Title       string `json:"title"`
	HTML        *template.HTML
	Plaintext   string
	UserID      int64
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Version     int
	IsProcessed int // 0 false, 1 true
}
