package sqldbrepo

import (
	"context"
	"time"

	"github.com/elkcityhazard/remind-me/internal/models"
)

func (sqdb *SQLDBRepo) UpdateUser(u *models.User) (int, error) {

	ctx, cancel := context.WithTimeout(sqdb.Config.Context, time.Second*10)

	defer cancel()

	versionChan := make(chan int, 1)
	errorChan := make(chan error, 1)

	sqdb.Config.WG.Add(1)

	go func() {
		defer sqdb.Config.WG.Done()
		defer close(versionChan)
		defer close(errorChan)

		// start a transaction

		tx, err := sqdb.Config.DB.Begin()

		if err != nil {
			errorChan <- err
			return
		}

		_, err = tx.ExecContext(ctx, "UPDATE User SET UpdatedAt=NOW(), IsActive=?,  Version = Version + 1 WHERE ID = ? AND Version = ?", u.IsActive, u.ID, u.Version)

		if err != nil {
			errorChan <- err
			tx.Rollback()
			return
		}

		_, err = tx.ExecContext(ctx, "UPDATE Password SET Hash=?, UpdatedAt=NOW(), Version = Version + 1 WHERE UserID = ? AND Version = ?", u.Password.Hash, u.ID, u.Password.Version)

		if err != nil {
			errorChan <- err
			tx.Rollback()
			return
		}

		_, err = tx.ExecContext(ctx, "UPDATE PhoneNumber SET Plaintext=?, UpdatedAt=NOW(), Version = Version + 1 WHERE UserID = ? AND Version = ?", u.PhoneNumber.Plaintext, u.ID, u.PhoneNumber.Version)

		if err != nil {
			errorChan <- err
			tx.Rollback()
			return
		}

		err = tx.Commit()

		if err != nil {
			errorChan <- err
			return
		}

		fetchedUser, err := sqdb.GetUserById(u.ID)

		if err != nil {
			errorChan <- err
			return

		}

		versionChan <- fetchedUser.Version

	}()

	err := <-errorChan

	if err != nil {
		return 0, err
	}

	version := <-versionChan

	return version, nil
}
