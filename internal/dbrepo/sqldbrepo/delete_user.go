package sqldbrepo

import (
	"context"
	"fmt"
	"time"
)

func (sqldb *SQLDBRepo) DeleteUser(id int64) error {

	ctx, cancel := context.WithTimeout(sqldb.Config.Context, time.Second*10)

	defer cancel()

	errorChan := make(chan error, 1)

	sqldb.Config.WG.Add(1)

	go func() {
		defer sqldb.Config.WG.Done()
		defer close(errorChan)

		row, err := sqldb.Config.DB.ExecContext(ctx, "DELETE FROM User WHERE id = ?", id)

		if err != nil {
			errorChan <- err
			return
		}

		deletedRows, err := row.RowsAffected()

		if err != nil {
			errorChan <- err
			return
		}

		app.InfoChan <- fmt.Sprintf("Deleted Row Count : %d", deletedRows)
	}()

	err := <-errorChan

	if err != nil {
		return err
	}

	return nil

}
