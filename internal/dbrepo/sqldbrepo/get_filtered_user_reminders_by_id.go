package sqldbrepo

import (
	"context"
	"time"

	"github.com/elkcityhazard/remind-me/internal/models"
)

func (sqdb *SQLDBRepo) GetFilteredUserRemindersByID(id int64, limit, offset int) ([]*models.Reminder, error) {

	ctx, cancel := context.WithTimeout(sqdb.Config.Context, time.Second*10)

	defer cancel()

	reminderChan := make(chan *models.Reminder)
	errorChan := make(chan error)
	reminderDoneChan := make(chan bool)

	sqdb.Config.WG.Add(1)

	go func() {
		defer sqdb.Config.WG.Done()

		stmt := `SELECT ID, Title, Content, UserID, DueDate, CreatedAt, UpdatedAt, Version FROM Reminder WHERE Reminder.UserID = ? AND DueDate > CURRENT_TIMESTAMP() LIMIT ? OFFSET ?`

		rows, err := sqdb.Config.DB.QueryContext(ctx, stmt, id, limit, offset)

		if err != nil {
			errorChan <- err
			return
		}

		defer rows.Close()

		// not sure if this is going to be the best way to do this but I'm not sure if I can accomplish it with a left join or inner join in terms
		// of scanning it to the destination Reminder.Schedule

		for rows.Next() {
			var r models.Reminder

			err = rows.Scan(&r.ID, &r.Title, &r.Content, &r.UserID, &r.DueDate, &r.CreatedAt, &r.UpdatedAt, &r.Version)

			if err != nil {
				errorChan <- err
				return
			}

			stmt := `SELECT ID, ReminderID, CreatedAt, UpdatedAt, DispatchTime, Version FROM Schedule WHERE ReminderID = ?`

			schedRow, err := sqdb.Config.DB.QueryContext(ctx, stmt, r.ID)

			if err != nil {
				errorChan <- err
				return
			}

			defer schedRow.Close()

			for schedRow.Next() {
				var s models.Schedule

				err = schedRow.Scan(&s.ID, &s.ReminderID, &s.CreatedAt, &s.UpdatedAt, &s.DispatchTime, &s.Version)

				if err != nil {
					errorChan <- err
					return
				}

				r.Schedule = append(r.Schedule, &s)
			}

			reminderChan <- &r
		}

		reminderDoneChan <- true

	}()

	// process grabbing the data for the user

	totalReminders := []*models.Reminder{}

	for {
		select {
		case err := <-errorChan:
			return nil, err
		case reminder := <-reminderChan:
			totalReminders = append(totalReminders, reminder)
		case <-reminderDoneChan:
			close(reminderChan)
			close(errorChan)
			close(reminderDoneChan)
			return totalReminders, nil

		}
	}

}
