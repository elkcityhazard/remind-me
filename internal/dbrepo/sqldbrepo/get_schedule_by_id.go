package sqldbrepo

import (
	"context"
	"time"

	"github.com/elkcityhazard/remind-me/internal/models"
)

func (sqdb *SQLDBRepo) GetScheduleByID(scheduleID int64) (*models.Schedule, error) {

	ctx, cancel := context.WithTimeout(sqdb.Config.Context, time.Second*10)

	defer cancel()

	scheduleChan := make(chan *models.Schedule, 1)

	errorChan := make(chan error, 1)

	sqdb.Config.WG.Add(1)

	go func() {
		defer sqdb.Config.WG.Done()
		defer close(scheduleChan)
		defer close(errorChan)

		var s models.Schedule

		stmt := `SELECT * FROM Schedule WHERE ID = ?`

		err := sqdb.Config.DB.QueryRowContext(ctx, stmt, scheduleID).Scan(&s.ID, &s.ReminderID, &s.CreatedAt, &s.UpdatedAt, &s.DispatchTime, &s.IsProcessed, &s.Version)

		if err != nil {
			errorChan <- err
			return
		}

		scheduleChan <- &s

	}()

	err := <-errorChan

	if err != nil {
		return nil, err
	}

	foundSchedule := <-scheduleChan

	return foundSchedule, nil

}
