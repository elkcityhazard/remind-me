package sqldbrepo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
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
			reminderDoneChan <- true
			return
		}

		fetchScheduledReminders, err := tx.QueryContext(ctx, "SELECT * FROM Schedule WHERE Schedule.DispatchTime < CURRENT_TIMESTAMP() AND Schedule.IsProcessed < 1 ORDER BY Schedule.DispatchTime ASC")

		if err != nil {
			tx.Rollback()
			errorChan <- err
			reminderDoneChan <- true
			return
		}

		defer fetchScheduledReminders.Close()

		var remindersToSend = []*models.Schedule{}

		// fetch the scheduled reminders

		for fetchScheduledReminders.Next() {
			s := models.Schedule{}

			err := fetchScheduledReminders.Scan(&s.ID, &s.ReminderID, &s.CreatedAt, &s.UpdatedAt, &s.DispatchTime, &s.IsProcessed, &s.Version)

			if err != nil {
				if err == sql.ErrNoRows {
					tx.Rollback()
					errorChan <- errors.New("no scheduled reminders right now")
					reminderDoneChan <- true
					return
				}
				tx.Rollback()
				errorChan <- err
				reminderDoneChan <- true
				return
			}

			remindersToSend = append(remindersToSend, &s)
		}

		if len(remindersToSend) == 0 {
			reminderDoneChan <- true
			return
		}

		// fetch the reminders associated with the scheduled outgoing notifications

		var idsToFetch = []int64{}

		for _, s := range remindersToSend {

			if !sliceContainsID(idsToFetch, s.ID) {
				idsToFetch = append(idsToFetch, s.ReminderID)
			}

		}

		reminderArgs := int64SliceToAnySlice(idsToFetch)

		substitutionString := createPlaceholdersFromSlice(idsToFetch)

		reminderStmt := fmt.Sprintf("SELECT * FROM Reminder WHERE ID IN (%s)", strings.TrimSpace(substitutionString))

		fetchReminderData, err := tx.QueryContext(ctx, reminderStmt, reminderArgs...)

		if err != nil {
			tx.Rollback()
			errorChan <- err
			reminderDoneChan <- true
			return
		}

		defer fetchScheduledReminders.Close()

		reminderData := []*models.Reminder{}

		for fetchReminderData.Next() {

			r := models.Reminder{}

			err := fetchReminderData.Scan(&r.ID, &r.Title, &r.Content, &r.UserID, &r.DueDate, &r.CreatedAt, &r.UpdatedAt, &r.Version)

			if err != nil {
				if err == sql.ErrNoRows {
					tx.Rollback()
					errorChan <- errors.New("no scheduled reminders right now")
					reminderDoneChan <- true
					return
				}
				tx.Rollback()
				errorChan <- err
				reminderDoneChan <- true
				return
			}

			reminderData = append(reminderData, &r)

		}

		if len(reminderData) == 0 {
			reminderDoneChan <- true
			return
		}

		userIDSlice := []int64{}

		for _, r := range reminderData {

			// check if the reminder has a user id foreign key

			if r.UserID > 0 {
				if !sliceContainsID(userIDSlice, r.UserID) {
					userIDSlice = append(userIDSlice, r.UserID)

				}
			}

			// loop through the reminders to send and append them to the reminder

			for _, s := range remindersToSend {

				if s.ReminderID == r.ID {
					r.Schedule = append(r.Schedule, s)
				}

			}
		}

		formattedUserIDs := int64SliceToAnySlice(userIDSlice)

		userSubstitutionString := createPlaceholdersFromSlice(userIDSlice)

		userStmt := fmt.Sprintf("SELECT ID, Email FROM User WHERE ID IN (%s)", strings.TrimSpace(userSubstitutionString))

		fetchUserData, err := tx.QueryContext(ctx, userStmt, formattedUserIDs...)

		if err != nil {
			tx.Rollback()
			errorChan <- err
			reminderDoneChan <- true
			return
		}

		defer fetchUserData.Close()

		userEmails := []*models.User{}

		for fetchUserData.Next() {
			u := models.User{}

			err := fetchUserData.Scan(&u.ID, &u.Email)
			if err != nil {
				if err == sql.ErrNoRows {
					tx.Rollback()
					errorChan <- errors.New("no scheduled reminders right now")
					reminderDoneChan <- true
					return
				}
				tx.Rollback()
				errorChan <- err
				reminderDoneChan <- true
				return
			}

			userEmails = append(userEmails, &u)

		}

		sqdb.Config.InfoChan <- fmt.Sprintf("%+v", userEmails)

		for _, u := range userEmails {
			emailData := make(map[string]any)
			for _, r := range reminderData {
				emailData["Reminder"] = r
				for _, s := range r.Schedule {
					if u.ID == r.UserID {

						emailData["Schedule"] = s
						emailData["User"] = u

						emailPayload := models.EmailData{
							Recipient:    u.Email,
							TemplateFile: "reminder-email.tmpl",
							Data:         emailData,
						}

						s.IsProcessed = 1
						s.UpdatedAt = time.Now()

						updateSched, err := tx.PrepareContext(ctx, "UPDATE Schedule SET IsProcessed = ?, UpdatedAt=NOW(), Version = Version + 1 WHERE ID = ? AND Version = ?")

						if err != nil {
							tx.Rollback()
							errorChan <- err
							reminderDoneChan <- true
							return
						}

						_, err = updateSched.ExecContext(ctx, 1, s.ID, s.Version)

						if err != nil {
							tx.Rollback()
							errorChan <- err
							reminderDoneChan <- true
							return
						}

						reminderChan <- r

						err = sqdb.SendNotification(r.DueDate, r.Title, r.Content)

						if err != nil {
							tx.Rollback()
							errorChan <- err
							reminderDoneChan <- true
							return
						}

						sqdb.Config.Mailer.MailerDataChan <- &emailPayload
					}
				}
			}
		}

		err = tx.Commit()

		if err != nil {
			tx.Rollback()
			errorChan <- err
			reminderDoneChan <- true
			return
		}

		reminderDoneChan <- true

	}()

	processedReminders := []models.Reminder{}

	for {
		select {
		case err := <-errorChan:
			sqdb.Config.ErrorChan <- err
		case reminder := <-reminderChan:
			sqdb.Config.InfoChan <- "receiving a new reminder"
			processedReminders = append(processedReminders, *reminder)
		case <-reminderDoneChan:
			return processedReminders, nil
		default:
			// Optionally, add logic here to handle the case where no values are ready to be received.
			sqdb.Config.InfoChan <- "nothing to do right now"
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

func createPlaceholdersFromSlice(s []int64) string {

	if len(s) == 0 {
		return ""
	}

	ph := []string{}

	for range s {
		ph = append(ph, "?")
	}

	if len(ph) > 0 {
		return strings.Join(ph, ",")
	}

	return ""
}

func int64SliceToAnySlice(intSlice []int64) []interface{} {
	anySlice := make([]interface{}, len(intSlice))
	for i, v := range intSlice {
		anySlice[i] = v
	}
	return anySlice
}
