package sqldbrepo

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/elkcityhazard/remind-me/internal/models"
)

func (sqdb *SQLDBRepo) ProcessAllReminders() ([]models.Reminder, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)

	defer cancel()

	reminderChan := make(chan *models.Reminder)
	errorChan := make(chan error)
	reminderDoneChan := make(chan bool)

	sqdb.Config.WG.Add(1)

	go func() {
		defer sqdb.Config.WG.Done()

		tx, err := sqdb.Config.DB.Begin()

		if err != nil {
			tx.Rollback()
			errorChan <- err
			return
		}

		fetchScheduledReminders, err := tx.QueryContext(ctx, "SELECT * FROM Schedule WHERE Schedule.DispatchTime < CURRENT_TIMESTAMP() AND Schedule.IsProcessed < 1 ORDER BY Schedule.DispatchTime ASC")

		if err != nil {
			tx.Rollback()
			errorChan <- err
			return
		}

		defer fetchScheduledReminders.Close()

		var remindersToSend = []*models.Schedule{}

		// fetch the scheduled reminders

		for fetchScheduledReminders.Next() {
			s := models.Schedule{}

			err := fetchScheduledReminders.Scan(&s.ID, &s.ReminderID, &s.CreatedAt, &s.UpdatedAt, &s.DispatchTime, &s.IsProcessed, &s.Version)

			if err != nil {
				tx.Rollback()
				errorChan <- err
				return
			}

			remindersToSend = append(remindersToSend, &s)
		}

		// fetch the reminders associated with the scheduled outgoing notifications

		var idsToFetch = []int64{}

		for _, s := range remindersToSend {

			if !sliceContainsID(idsToFetch, s.ID) {
				idsToFetch = append(idsToFetch, s.ID)
			}

		}

		formattedIDs := joinIntsToString(idsToFetch)

		fetchReminderData, err := tx.QueryContext(ctx, fmt.Sprintf("SELECT * FROM Reminder WHERE id IN (%s)", strings.Join(formattedIDs, ",")))

		if err != nil {
			tx.Rollback()
			errorChan <- err
			return
		}

		defer fetchScheduledReminders.Close()

		reminderData := []*models.Reminder{}

		for fetchReminderData.Next() {

			r := models.Reminder{}

			err := fetchReminderData.Scan(&r.ID, &r.Title, &r.Content, &r.UserID, &r.DueDate, &r.CreatedAt, &r.UpdatedAt, &r.Version)

			if err != nil {
				tx.Rollback()
				errorChan <- err
				return
			}

			reminderData = append(reminderData, &r)

		}

		for _, r := range reminderData {

			for _, s := range remindersToSend {

				if s.ReminderID == r.ID {
					r.Schedule = append(r.Schedule, s)
				}

			}
		}

		userIDSlice := []int64{}

		for _, r := range reminderData {

			if !sliceContainsID(userIDSlice, r.UserID) {
				userIDSlice = append(userIDSlice, r.UserID)
			}

		}

		formattedUserIDs := joinIntsToString(userIDSlice)

		fetchUserData, err := tx.QueryContext(ctx, fmt.Sprintf("SELECT ID, Email FROM User WHERE ID in (%s)", strings.Join(formattedUserIDs, ",")))

		if err != nil {
			tx.Rollback()
			errorChan <- err
			return
		}

		defer fetchUserData.Close()

		userEmails := []*models.User{}

		for fetchUserData.Next() {
			u := models.User{}

			err := fetchUserData.Scan(&u.ID, &u.Email)
			if err != nil {
				tx.Rollback()
				errorChan <- err
				return
			}

			userEmails = append(userEmails, &u)

		}

		for _, u := range userEmails {
			for _, r := range reminderData {
				for _, s := range r.Schedule {
					if u.ID == r.UserID {
						emailData := make(map[string]any)

						emailData["Reminder"] = r
						emailData["Schedule"] = s
						emailData["User"] = u.Email

						emailPayload := models.EmailData{
							Recipient:    u.Email,
							TemplateFile: "reminder-email.tmpl",
							Data:         emailData,
						}

						s.IsProcessed = 1
						s.UpdatedAt = time.Now()

						updateSched, err := tx.PrepareContext(ctx, "UPDATE Schedule SET IsProcessed = ?, UpdatedAt=NOW(), Version = Version + 1 WHERE ID = ? AND Version = ?")

						if err != nil {
							fmt.Println(err)
							tx.Rollback()
							errorChan <- err
							return
						}

						_, err = updateSched.ExecContext(ctx, 1, s.ID, s.Version)

						if err != nil {
							tx.Rollback()
							errorChan <- err
							return
						}

						reminderChan <- r

						err = sqdb.SendNotification(r.DueDate, r.Title, r.Content)

						if err != nil {
							tx.Rollback()
							errorChan <- err
							return
						}

						sqdb.Config.Mailer.MailerDataChan <- &emailPayload
					}
				}
			}
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

// sliceContainsID is a help function to only fetch the id's we need for the scheduled reminders
// it accepts a slice of int53, and a value to check. and returns a bool
func sliceContainsID(s []int64, v int64) bool {
	for _, vins := range s {
		if vins == v {
			return true
		}
	}

	return false
}

func joinIntsToString(s []int64) []string {

	tmp := []string{}

	for _, v := range s {

		itos := strconv.FormatInt(v, 10)

		tmp = append(tmp, itos)

	}

	return tmp

}
