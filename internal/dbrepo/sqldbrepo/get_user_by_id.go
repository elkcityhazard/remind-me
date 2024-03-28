package sqldbrepo

import (
	"context"
	"fmt"
	"time"

	"github.com/elkcityhazard/remind-me/internal/models"
)

func (sqdb *SQLDBRepo) GetUserById(id int64) (*models.User, error) {

	ctx, cancel := context.WithTimeout(sqdb.Config.Context, time.Second*10)

	defer cancel()

	userChan := make(chan *models.User, 1)
	errorChan := make(chan error, 1)

	defer close(userChan)
	defer close(errorChan)

	go func() {

		stmt := `SELECT user.ID, user.Email, user.CreatedAt, user.UpdatedAt, user.Scope, user.IsActive, user.Version FROM User WHERE ID = ?`

		var u = models.User{}

		err := sqdb.Config.DB.QueryRowContext(ctx, stmt, id).Scan(&u.ID, &u.Email, nil, nil, &u.CreatedAt, &u.UpdatedAt, &u.Scope, &u.IsActive, &u.Version)

		if err != nil {
			fmt.Println(u)
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
