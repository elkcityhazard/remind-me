package dbrepo

import "github.com/elkcityhazard/remind-me/internal/models"

type DBServicer interface {
	InsertUser(*models.User) (int64, error)
	GetUserById(int64) (*models.User, error)
	GetUserbyEmail(email string) (*models.User, error)
	UpdateUser(*models.User) (int, error) // int is version
	DeleteUser(int64) error
}
