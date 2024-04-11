package sqldbrepo

import (
	"context"
	"fmt"
	"time"

	"github.com/elkcityhazard/remind-me/internal/models"
)

func (sqdb *SQLDBRepo) UpdateReminder(reminder *models.Reminder) (*models.Reminder, error) {

	ctx, cancel := context.WithTimeout(sqdb.Config.Context, time.Second*10)

	defer cancel()

	errChan := make(chan error, 1)
	reminderChan := make(chan *models.Reminder, 1)

	sqdb.Config.WG.Add(1)

	go func() {
		defer sqdb.Config.WG.Done()
		defer close(errChan)
		defer close(reminderChan)

		stmt := `UPDATE Reminder Set 
		Title=?, 
		Content=?, 
		DueDate=?, 
		UpdatedAt=?, 
		Version= Version + 1 
		WHERE ID = ? 
		AND Version = ?`

		args := []interface{}{reminder.Title, reminder.Content, reminder.DueDate, reminder.UpdatedAt, reminder.ID, reminder.Version}

		tx, err := sqdb.Config.DB.Begin()

		if err != nil {
			errChan <- err
			tx.Rollback()
			return
		}

		result, err := tx.ExecContext(ctx, stmt, args...)

		if err != nil {
			errChan <- err
			tx.Rollback()
			return
		}

		affected, err := result.RowsAffected()

		if err != nil {
			errChan <- err
			tx.Rollback()
			return
		}

		err = tx.Commit()

		if err != nil {
			errChan <- err
			tx.Rollback()
			return
		}

		sqdb.Config.InfoChan <- fmt.Sprintf("%d rows affected", affected)

		reminder.Version = reminder.Version + 1

		reminderChan <- reminder

	}()

	err := <-errChan

	if err != nil {
		return nil, err
	}

	newReminder := <-reminderChan

	return newReminder, nil

}
