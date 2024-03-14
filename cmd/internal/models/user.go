package models

import "time"

type Password struct {
	Plaintext string `json:"-"`
	Hash      []byte `json:"-"`
}

type PhoneNumber struct {
	Plantext string `json:"plaintext"`
	Hash     []byte `json:"-"`
}

type User struct {
	ID          int64       `json:"id"`
	Email       string      `json:"email"`
	Password    Password    `json:"-"`
	PhoneNumber PhoneNumber `json:"phone_number"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
	Scope       int         `json:"-"`
	IsActive    bool        `json:"is_active"`
	Version     int         `json:"version"`
}
