package models

import (
	"strings"
	"time"
)

type ActivationToken struct {
	ID          int64     `json:"id"`
	UserID      int64     `json:"user_id"`
	Token       []byte    `json:"token"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	IsProcessed int       `json:"is_processed"`
}

type Password struct {
	ID         int64     `json:"id"`
	Plaintext1 string    `json:""`
	Plaintext2 string    `json:""`
	Hash       []byte    `json:"-"`
	Salt       []byte    `json:"-"`
	UserID     int64     `json:"-"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	Version    int       `json:"version"`
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
	Scope       int         `json:"scope"`
	IsActive    int         `json:"is_active"`
	Version     int         `json:"version"`
}

//	PasswordsMatch takes in two provided passwords, and returns whether or not they match

func (u *User) PasswordsMatch(password1, password2 string) bool {
	return strings.EqualFold(password1, password2)
}
