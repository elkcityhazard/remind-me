package dbrepo

import (
	"time"

	"github.com/elkcityhazard/remind-me/internal/models"
)

type DBServicer interface {
	InsertUser(*models.User) (int64, error)
	GetUserById(int64) (*models.User, error)
	GetUserbyEmail(email string) (*models.User, error)
	UpdateUser(*models.User) (int, error) // int is version
	DeleteUser(int64) error
	ActiveUser(activationToken string, id int64) (*models.User, error)

	InsertReminder(r *models.Reminder) (int64, error)
	GetReminderSchedulesByID(reminderID int64) ([]*models.Schedule, error)
	GetReminderByID(reminderID int64) (*models.Reminder, error)
	GetScheduleByID(scheduleID int64) (*models.Schedule, error)
	UpdateReminder(reminder *models.Reminder) ([]*models.Reminder, error)
	UpdateScheduleByID(schedule *models.Schedule) (*models.Schedule, error)
	GetUserRemindersByID(id int64) ([]*models.Reminder, error)
	GetFilteredUserRemindersByID(id int64, limit int, offset int) ([]*models.Reminder, error)
	ProcessRemindersForUser(id int64) ([]models.Reminder, error)
	ProcessAllReminders() ([]models.Reminder, error)
	SendNotification(dueDate time.Time, title, content string) error
}
