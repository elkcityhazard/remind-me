package sqldbrepo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/elkcityhazard/remind-me/internal/models"
)

func (sqdb *SQLDBRepo) ProcessRemindersForUser(id int64) ([]models.Reminder, error) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)

	defer cancel()

	reminderChan := make(chan *models.Reminder)
	errorChan := make(chan error)
	reminderDoneChan := make(chan bool)

	sqdb.Config.WG.Add(1)

	go func() {
		defer sqdb.Config.WG.Done()

		rows, err := sqdb.Config.DB.QueryContext(ctx, `SELECT Schedule.ID, Schedule.ReminderID, Schedule.DispatchTime, Schedule.Version, 
		Reminder.Title, Reminder.Content, Reminder.DueDate, Reminder.Version, 
		User.Email FROM Schedule 
		INNER JOIN Reminder ON Reminder.ID = Schedule.ReminderID 
		INNER JOIN User ON User.Id = Reminder.UserID  
		WHERE Reminder.UserID = ? AND Schedule.DispatchTime < CURRENT_TIMESTAMP() AND Schedule.IsProcessed < 1 AND Reminder.DueDate > CURRENT_TIMESTAMP() 
		ORDER BY Schedule.DispatchTime ASC`, id)

		if err != nil {

			if err == sql.ErrNoRows {
				errorChan <- errors.New("no rows to process")
				return
			}

			errorChan <- err
			return
		}
		defer rows.Close()

		tx, err := sqdb.Config.DB.Begin()

		if err != nil {
			tx.Rollback()
			errorChan <- err
			return

		}

		type ScheduledReminderStruct struct {
			Reminder models.Reminder
			Schedule models.Schedule
			User     models.User
		}

		scheduledReminders := []*ScheduledReminderStruct{}

		for rows.Next() {

			srm := ScheduledReminderStruct{}

			r := models.Reminder{}
			s := models.Schedule{}
			u := models.User{}

			u.ID = id
			r.UserID = id

			err := rows.Scan(&s.ID, &s.ReminderID, &s.DispatchTime, &s.Version, &r.Title, &r.Content, &r.DueDate, &r.Version, &u.Email)

			if err != nil {

				tx.Rollback()
				errorChan <- err
				return
			}

			srm.Reminder = r
			srm.Schedule = s
			srm.User = u

			scheduledReminders = append(scheduledReminders, &srm)

		}

		for _, srm := range scheduledReminders {

			emailData := make(map[string]any)

			emailData["Reminder"] = srm.Reminder
			emailData["Schedule"] = srm.Schedule
			emailData["User"] = srm.User

			emailPayload := models.EmailData{
				Recipient:    srm.User.Email,
				TemplateFile: "reminder-email.tmpl",
				Data:         emailData,
			}

			srm.Schedule.IsProcessed = 1
			srm.Schedule.UpdatedAt = time.Now()

			updateSched, err := tx.PrepareContext(ctx, "UPDATE Schedule SET IsProcessed = ?, UpdatedAt=NOW(), Version = Version + 1 WHERE ID = ? AND Version = ?")

			if err != nil {
				fmt.Println(err)
				tx.Rollback()
				errorChan <- err
				return
			}

			_, err = updateSched.ExecContext(ctx, 1, srm.Schedule.ID, srm.Schedule.Version)

			if err != nil {
				tx.Rollback()
				errorChan <- err
				return
			}

			reminderChan <- &srm.Reminder

			err = sqdb.SendNotification(srm.Reminder.DueDate, srm.Reminder.Title, srm.Reminder.Content)

			if err != nil {
				tx.Rollback()
				errorChan <- err
				return
			}

			sqdb.Config.Mailer.MailerDataChan <- &emailPayload

		}

		err = tx.Commit()

		if err != nil {
			fmt.Println(err)
			tx.Rollback()
			errorChan <- err
			return
		}

		reminderDoneChan <- true
	}()

	processedReminders := []models.Reminder{}

	for {
		select {
		case err := <-errorChan:
			return nil, err
		case reminder := <-reminderChan:
			sqdb.Config.InfoChan <- "receiving a new reminder"
			processedReminders = append(processedReminders, *reminder)
		case <-reminderDoneChan:
			return processedReminders, nil

		}
	}

}
