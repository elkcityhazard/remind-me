package sqldbrepo

import (
	"context"
	"time"

	"github.com/elkcityhazard/remind-me/internal/models"
)

func (sqdb *SQLDBRepo) GetUserById(id int64) (*models.User, error) {

	ctx, cancel := context.WithTimeout(sqdb.Config.Context, time.Second*10)

	defer cancel()

	userChan := make(chan *models.User, 1)
	errorChan := make(chan error, 1)

	go func() {

		defer close(userChan)
		defer close(errorChan)

		stmt := `SELECT User.ID, User.Email, User.IsActive, User.Version, User.Scope, PhoneNumber.Plaintext, Password.Hash FROM User INNER JOIN Password ON User.ID = Password.UserID INNER JOIN PhoneNumber ON User.ID = PhoneNumber.UserID WHERE User.ID = ?`

		var u = models.User{}

		err := sqdb.Config.DB.QueryRowContext(ctx, stmt, id).Scan(&u.ID, &u.Email, &u.IsActive, &u.Version, &u.Scope, &u.PhoneNumber.Plaintext, &u.Password.Hash)

		if err != nil {
			sqdb.Config.ErrorChan <- err
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
