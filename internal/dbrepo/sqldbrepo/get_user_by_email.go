package sqldbrepo

import (
	"context"
	"time"

	"github.com/elkcityhazard/remind-me/internal/models"
)

func (sqdb *SQLDBRepo) GetUserByEmail(email string) (*models.User, error) {

	ctx, cancel := context.WithTimeout(sqdb.Config.Context, time.Second*10)

	defer cancel()

	userChan := make(chan *models.User, 1)
	errorChan := make(chan error, 1)

	go func() {

		defer close(userChan)
		defer close(errorChan)

		stmt := `SELECT User.ID, User.Email, User.CreatedAt, User.UpdatedAt, User.IsActive, User.Version, User.Scope, PhoneNumber.ID, PhoneNumber.Plaintext, PhoneNumber.CreatedAt, PhoneNumber.UpdatedAt, PhoneNumber.Version, Password.ID, Password.Hash, Password.CreatedAt, Password.UpdatedAt, Password.Version FROM User INNER JOIN Password ON User.ID = Password.UserID INNER JOIN PhoneNumber ON User.ID = PhoneNumber.UserID WHERE User.Email = ?`

		var u = models.User{}

		err := sqdb.Config.DB.QueryRowContext(ctx, stmt, email).Scan(&u.ID, &u.Email, &u.CreatedAt, &u.UpdatedAt, &u.IsActive, &u.Version, &u.Scope, &u.PhoneNumber.ID, &u.PhoneNumber.Plaintext, &u.PhoneNumber.CreatedAt, &u.PhoneNumber.UpdatedAt, &u.PhoneNumber.Version, &u.Password.ID, &u.Password.Hash, &u.Password.CreatedAt, &u.Password.UpdatedAt, &u.Password.Version)

		if err != nil {
			errorChan <- err
		}

		userChan <- &u

	}()

	err := <-errorChan

	if err != nil {
		return nil, err
	}

	user := <-userChan

	return user, nil

}
