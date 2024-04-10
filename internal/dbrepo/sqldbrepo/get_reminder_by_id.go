package sqldbrepo

import (
	"context"
	"time"

	"github.com/elkcityhazard/remind-me/internal/models"
)

func (sqdb *SQLDBRepo) GetReminderSchedulesByID(reminderID int64) ([]*models.Schedule, error) {

	ctx, cancel := context.WithTimeout(sqdb.Config.Context, time.Second*10)

	defer cancel()

	scheduleChan := make(chan *models.Schedule)
	scheduleDoneChan := make(chan bool)
	errorChan := make(chan error)

	sqdb.Config.WG.Add(1)

	go func() {

		defer sqdb.Config.WG.Done()

		stmt := `SELECT * FROM Schedule WHERE reminderID = ? AND IsProcessed < 1`

		rows, err := sqdb.Config.DB.QueryContext(ctx, stmt, reminderID)

		if err != nil {
			errorChan <- err
			scheduleDoneChan <- true
			return
		}

		defer rows.Close()

		for rows.Next() {
			s := models.Schedule{}

			rows.Scan(&s.ID, &s.ReminderID, &s.CreatedAt, &s.UpdatedAt, &s.DispatchTime, &s.IsProcessed, &s.Version)

			scheduleChan <- &s

		}

		scheduleDoneChan <- true

	}()

	var schedules []*models.Schedule

	for {
		select {
		case err := <-errorChan:
			if err != nil {
				close(errorChan)
				close(scheduleChan)
				close(scheduleDoneChan)
				return nil, err
			}
		case s := <-scheduleChan:
			schedules = append(schedules, s)
		case <-scheduleDoneChan:
			close(errorChan)
			close(scheduleChan)
			close(scheduleDoneChan)
			return schedules, nil

		}
	}

}

func (sqdb *SQLDBRepo) GetReminderByID(reminderID int64) (*models.Reminder, error) {

	ctx, cancel := context.WithTimeout(sqdb.Config.Context, time.Second*10)

	defer cancel()

	reminderChan := make(chan *models.Reminder, 1)

	errorChan := make(chan error, 1)

	sqdb.Config.WG.Add(1)

	go func() {
		defer sqdb.Config.WG.Done()
		defer close(reminderChan)
		defer close(errorChan)

		var r models.Reminder

		stmt := `SELECT * FROM Reminder WHERE ID = ?`

		err := sqdb.Config.DB.QueryRowContext(ctx, stmt, reminderID).Scan(&r.ID, &r.Title, &r.Content, &r.UserID, &r.DueDate, &r.CreatedAt, &r.UpdatedAt, &r.Version)

		if err != nil {
			errorChan <- err
			return
		}

		r.Schedule, err = sqdb.GetReminderSchedulesByID(reminderID)

		if err != nil {
			errorChan <- err
			return
		}

		reminderChan <- &r

	}()

	err := <-errorChan

	if err != nil {
		return nil, err
	}

	foundReminder := <-reminderChan

	return foundReminder, nil

}
