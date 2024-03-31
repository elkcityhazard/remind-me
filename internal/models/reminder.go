package models

import (
	"time"
)

type Reminder struct {
	ID        int64       `json:"id"`
	Title     string      `json:"title"`
	Content   string      `json:"content"`
	UserID    int64       `json:"user_id"`
	DueDate   time.Time   `json:"due_date"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
	Schedule  []*Schedule `json:"schedule"`
	Version   int         `json:"version"`
}
