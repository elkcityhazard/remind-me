package models

import "time"

type Schedule struct {
	ID           int64     `json:"id"`
	ReminderID   int64     `json:"reminder_id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	DispatchTime time.Time `json:"dispatch_time"`
	Version      int       `json:"version"`
	IsProcessed  int       `json:"is_processed"`
}
