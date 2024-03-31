package sqldbrepo

import (
	"context"
	"log"
	"time"

	"github.com/elkcityhazard/remind-me/internal/models"
)

func (sqdb *SQLDBRepo) ActivateUser(activationToken string, id int64) (*models.User, error) {

	ctx, cancel := context.WithTimeout(sqdb.Config.Context, time.Second*10)

	defer cancel()

	userChan := make(chan *models.User, 1)
	errorChan := make(chan error, 1)

	sqdb.Config.WG.Add(1)

	go func() {
		defer sqdb.Config.WG.Done()
		defer close(userChan)
		defer close(errorChan)

		tx, err := sqdb.Config.DB.Begin()

		if err != nil {
			errorChan <- err
			return
		}

		// fetch user by token
		t := &models.ActivationToken{}
		u := &models.User{}

		u.ID = id

		err = tx.QueryRowContext(ctx, `SELECT User.ID, User.Email, User.CreatedAt, User.UpdatedAt, User.Scope, User.IsActive, User.Version, 
		ActivationToken.ID, ActivationToken.UserId, ActivationToken.Token, ActivationToken.CreatedAt, ActivationToken.UpdatedAt, ActivationToken.IsProcessed 
		FROM User INNER JOIN ActivationToken ON User.ID = ActivationToken.UserID
		WHERE ActivationToken.Token = ? AND User.ID = ? AND ActivationToken.IsProcessed != 1`, activationToken, id).Scan(&u.ID, &u.Email, &u.CreatedAt, &u.UpdatedAt, &u.Scope, &u.IsActive, &u.Version, &t.ID, &t.UserID, &t.Token, &t.CreatedAt, &t.UpdatedAt, &t.IsProcessed)

		if err != nil {
			tx.Rollback()
			errorChan <- err
			return
		}

		// update user

		u.IsActive = 1
		u.UpdatedAt = time.Now()

		_, err = tx.ExecContext(ctx, "UPDATE User SET UpdatedAt = ?, IsActive = ?, Version = Version + 1 WHERE ID=? AND Version = ?", u.UpdatedAt, u.IsActive, u.ID, u.Version)

		if err != nil {
			tx.Rollback()
			errorChan <- err
			return
		}

		_, err = tx.ExecContext(ctx, "UPDATE ActivationToken SET IsProcessed = 1, UpdatedAt=NOW() WHERE ID = ?", t.ID)

		if err != nil {
			tx.Rollback()
			errorChan <- err
			return
		}

		err = tx.Commit()

		if err != nil {
			errorChan <- err
			return
		}

		updatedUser, err := sqdb.GetUserById(id)

		if err != nil {
			log.Println(id, u.ID, err)
			errorChan <- err
			return
		}

		userChan <- updatedUser

	}()

	err := <-errorChan

	if err != nil {
		return nil, err
	}

	user := <-userChan

	return user, nil

}
