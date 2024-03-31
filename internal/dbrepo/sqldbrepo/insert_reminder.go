package sqldbrepo

import (
	"context"
	"time"

	"github.com/elkcityhazard/remind-me/internal/models"
)

func (sqdb *SQLDBRepo) InsertReminder(r *models.Reminder) (int64, error) {
	ctx, cancel := context.WithTimeout(sqdb.Config.Context, time.Second*10)

	defer cancel()

	idChan := make(chan int64, 1)
	errorChan := make(chan error, 1)

	sqdb.Config.WG.Add(1)

	go func() {
		defer sqdb.Config.WG.Done()
		defer close(idChan)
		defer close(errorChan)

		tx, err := sqdb.Config.DB.Begin()

		if err != nil {
			tx.Rollback()
			errorChan <- err
			return
		}

		// insert the reminder

		reminderArgs := []any{&r.Title, &r.Content, &r.UserID, &r.DueDate}

		reminderResult, err := tx.ExecContext(ctx, `INSERT INTO Reminder (Title, Content, UserID, DueDate, CreatedAt, UpdatedAt, Version) VALUES(?,?,?,?,NOW(),NOW(),1)`, reminderArgs...)

		if err != nil {
			tx.Rollback()
			errorChan <- err
			return
		}

		reminderID, err := reminderResult.LastInsertId()

		if err != nil {
			tx.Rollback()
			errorChan <- err
			return
		}

		r.ID = reminderID

		// insert the schedule itsems

		scheduleStmt, err := tx.PrepareContext(ctx, "INSERT INTO Schedule(ReminderID, CreatedAt, UpdatedAt, DispatchTime, Version) VALUES (?, NOW(), NOW(), ?, 1)")

		if err != nil {
			tx.Rollback()
			errorChan <- err
			return
		}

		updateSchedule := []*models.Schedule{}

		for _, s := range r.Schedule {
			args := []any{r.ID, s.DispatchTime}
			scheduleResult, err := scheduleStmt.ExecContext(ctx, args...)

			if err != nil {
				tx.Rollback()
				errorChan <- err
				return
			}

			scheduleID, err := scheduleResult.LastInsertId()

			if err != nil {
				tx.Rollback()
				errorChan <- err
				return
			}

			s.ID = scheduleID

			updateSchedule = append(updateSchedule, s)
		}

		r.Schedule = updateSchedule

		err = tx.Commit()

		if err != nil {
			tx.Rollback()
			errorChan <- err
			return
		}

		idChan <- r.ID

	}()

	err := <-errorChan

	if err != nil {
		return 0, err
	}

	id := <-idChan

	return id, nil

}
