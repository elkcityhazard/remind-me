package sqldbrepo

import (
	"context"
	"errors"
	"time"

	"github.com/elkcityhazard/remind-me/internal/models"
	"github.com/elkcityhazard/remind-me/pkg/utils"
	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

const (
	notActive int = iota
	active
)

// InsertUser accepts a pointer declaration to a user, and inserts it into the database.
// It will return the User ID, and an error if there are any
func (sqdb *SQLDBRepo) InsertUser(u *models.User) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)

	defer cancel()

	idChan := make(chan int64, 1)
	errorChan := make(chan error, 1)

	sqdb.Config.WG.Add(1)

	go func() {
		defer close(idChan)
		defer close(errorChan)
		defer sqdb.Config.WG.Done()

		tx, err := sqdb.Config.DB.Begin()
		if err != nil {
			errorChan <- err
			return
		}

		userRow, err := tx.ExecContext(ctx, "INSERT INTO User (Email, CreatedAt, UpdatedAt, Scope, IsActive, Version) VALUES (?,NOW(),NOW(), 2, ?, 1)", u.Email, notActive)
		if err != nil {
			tx.Rollback()
			errorChan <- err
			return
		}

		userID, err := userRow.LastInsertId()
		if err != nil {
			errorChan <- err
			return
		}

		u.ID = userID

		pwRow, err := tx.ExecContext(ctx, "INSERT INTO Password (Hash, UserID, CreatedAt, UpdatedAt, Version) VALUES (?,?, NOW(), NOW(), 1)", u.Password.Hash, u.ID)
		if err != nil {
			tx.Rollback()
			errorChan <- err
			return
		}

		pwID, err := pwRow.LastInsertId()
		if err != nil {
			tx.Rollback()
			errorChan <- err
			return
		}

		u.Password.ID = pwID

		phoneRow, err := tx.ExecContext(ctx, "INSERT INTO PhoneNumber (Plaintext, UserID, CreatedAt, UpdatedAt, Version) Values (?, ?, NOW(), NOW(), 1)", u.PhoneNumber.Plaintext, userID)
		if err != nil {
			tx.Rollback()
			errorChan <- err
			return
		}

		phoneID, err := phoneRow.LastInsertId()
		if err != nil {
			tx.Rollback()
			errorChan <- err
			return
		}

		u.PhoneNumber.ID = phoneID

		activationToken, err := utils.NewUtils(sqdb.Config).GenerateActivationToken()
		if err != nil {
			tx.Rollback()
			errorChan <- err
			return
		}

		encryptedToken, err := bcrypt.GenerateFromPassword([]byte(activationToken), 10)
		if err != nil {
			tx.Rollback()
			errorChan <- err
			return
		}

		_, err = tx.ExecContext(ctx, "INSERT INTO ActivationToken (UserID, Token, CreatedAt, UpdatedAt, IsProcessed) Values(?,?, NOW(), NOW(), false)", userID, encryptedToken)
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

		emailData := make(map[string]interface{})

		emailData["ID"] = activationToken

		ed := models.EmailData{
			Recipient:    u.Email,
			TemplateFile: "user_welcome.tmpl",
			Data:         emailData,
		}

		sqdb.Config.Mailer.MailerDataChan <- &ed
	}()

	err := <-errorChan
	if err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok {
			switch mysqlErr.Number {
			case 1062: // MySQL error code for duplicate key
				return 0, errors.New("something has gone awry")
			}
			return 0, err
		}
	}

	userID := <-idChan

	return userID, nil
}
