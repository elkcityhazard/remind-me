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

	go func() {
		defer close(userChan)
		defer close(errorChan)

		//stmt := `SELECT Reminder.ID, Reminder.Title, Reminder.Content, Reminder.UserID, Reminder.DueDate, Reminder.CreatedAt, Reminder.UpdatedAt, Reminder.Version, Schedule.ID, Schedule.ReminderID, Schedule.CreatedAt, Schedule.UpdatedAt, Schedule.DispatchTime, Schedule.Version FROM Reminder INNER JOIN Schedule ON Schedule.ReminderID = Reminder.ID INNER JOIN User ON User.ID = Reminder.UserID WHERE User.ID = ?`

		stmt := `SELECT User.ID, User.Email, User.CreatedAt, User.UpdAtedAt, User.Scope, User.IsActive, User.Version, PhoneNumber.ID, PhoneNumber.Plaintext, PhoneNumber.UserID, PhoneNumber.CreatedAt, PhoneNumber.UpdatedAt, PhoneNumber.Version FROM User INNER JOIN PhoneNumber ON PhoneNumber.UserID = User.ID WHERE User.ID = ? AND User.IsActive > 0`

		u := models.User{}

		err := sqdb.Config.DB.QueryRowContext(ctx, stmt, id).Scan(&u.ID, &u.Email, &u.CreatedAt, &u.UpdatedAt, &u.Scope, &u.IsActive, &u.Version, &u.PhoneNumber.ID, &u.PhoneNumber.Plaintext, &u.PhoneNumber.UserID, &u.PhoneNumber.CreatedAt, &u.PhoneNumber.UpdatedAt, &u.PhoneNumber.Version)
		if err != nil {
			errorChan <- err
			return
		}

		userChan <- &u
	}()

	err := <-errorChan
	if err != nil {
		return nil, err
	}

	user := <-userChan

	app.InfoChan <- fmt.Sprintf("fetched user %d", user.ID)

	return user, nil
}
