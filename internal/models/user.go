package models

import "time"

type Password struct {
	ID        int64     `json:"id"`
	Plaintext string    `json:"-"`
	Hash      []byte    `json:"-"`
	Salt      []byte    `json:"-"`
	UserID    int64     `json:"-"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Version   int       `json:"version"`
}

type PhoneNumber struct {
	ID        int64     `json:"id"`
	Plaintext string    `json:"plaintext"`
	Hash      []byte    `json:"-"`
	UserID    int64     `json:"-"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Version   int       `json:"version"`
}

type User struct {
	ID          int64       `json:"id"`
	Email       string      `json:"email"`
	Password    Password    `json:"password"`
	PhoneNumber PhoneNumber `json:"phone_number"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
	Scope       int         `json:"-"`
	IsActive    bool        `json:"is_active"`
	Version     int         `json:"version"`
}
