package sqldbrepo

import (
	"context"
	"fmt"
	"time"

	"github.com/elkcityhazard/remind-me/internal/models"
)

func (sqdb *SQLDBRepo) UpdateScheduleByID(schedule *models.Schedule) (*models.Schedule, error) {

	ctx, cancel := context.WithTimeout(sqdb.Config.Context, time.Second*10)

	defer cancel()

	errorChan := make(chan error, 1)
	scheduleChan := make(chan *models.Schedule, 1)

	sqdb.Config.WG.Add(1)

	go func() {
		defer sqdb.Config.WG.Done()
		defer close(errorChan)
		defer close(scheduleChan)

		stmt := `UPDATE Schedule Set 
		UpdatedAt=?, 
		DispatchTime=?,
		IsProcessed=?,
		Version= Version + 1 
		WHERE ID = ? 
		AND Version = ?
		AND IsProcessed < 1
		`

		args := []interface{}{schedule.UpdatedAt, schedule.DispatchTime, schedule.IsProcessed, schedule.ID, schedule.Version}

		tx, err := sqdb.Config.DB.Begin()

		if err != nil {
			errorChan <- err
			tx.Rollback()
			return
		}

		result, err := tx.ExecContext(ctx, stmt, args...)

		if err != nil {
			errorChan <- err
			tx.Rollback()
			return
		}

		affected, err := result.RowsAffected()

		if err != nil {
			errorChan <- err
			tx.Rollback()
			return
		}

		err = tx.Commit()

		if err != nil {
			errorChan <- err
			tx.Rollback()
			return
		}

		sqdb.Config.InfoChan <- fmt.Sprintf("%d rows affected", affected)

		schedule.Version = schedule.Version + 1

		scheduleChan <- schedule

	}()

	err := <-errorChan

	if err != nil {
		return nil, err
	}

	newScheduleItem := <-scheduleChan

	return newScheduleItem, nil

}
