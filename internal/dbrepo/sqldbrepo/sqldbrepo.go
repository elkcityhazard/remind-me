package sqldbrepo

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/elkcityhazard/remind-me/internal/config"
	_ "github.com/go-sql-driver/mysql"
)

const (
	SuperUser int = iota
	Admin
	User
)

type SQLDBRepo struct {
	Config *config.AppConfig
}

var (
	app    *config.AppConfig
	DBRepo *SQLDBRepo
)

func NewSQLDBRepo(ac *config.AppConfig) *SQLDBRepo {
	app = ac
	return &SQLDBRepo{
		Config: ac,
	}
}

//  NewDatabaseConn establishes a new database connection
//  and sets the DB property on the appconfig.  It returns
//  an error if it fails.

func (sqdb *SQLDBRepo) NewDatabaseConn() (*sql.DB, error) {
	sqdb.Config.InfoLog.Println("Opening database connection...")
	conn, err := sql.Open("mysql", sqdb.Config.DSN)
	if err != nil {
		sqdb.Config.ErrorChan <- err
		return nil, err
	}

	conn.SetMaxIdleConns(10)
	conn.SetMaxOpenConns(100)
	conn.SetConnMaxIdleTime(time.Minute * 5)
	conn.SetConnMaxLifetime(time.Hour)
	// we pass it back up to app config as well in case we need it later

	sqdb.Config.DB = conn

	DBRepo = sqdb

	return conn, nil
}

func GetDatabaseConnection() *SQLDBRepo {
	return DBRepo
}

func (sqdb *SQLDBRepo) SendNotification(dueDate time.Time, title, content string) error {

	errorChan := make(chan error, 1)

	sqdb.Config.WG.Add(1)

	go func() {

		defer sqdb.Config.WG.Done()
		defer close(errorChan)

		client := &http.Client{}

		req, err := http.NewRequest("POST", "https://ntfy.sh/megalawnalien_reminders_3739",
			strings.NewReader(fmt.Sprintf("Due: %s - %s", dueDate.Format("Jan 02, 2006, 03:04:05pm"), content)))

		if err != nil {
			errorChan <- err
			return
		}
		req.Header.Set("Title", title)
		req.Header.Set("Priority", "default")
		req.Header.Set("Tags", "envelope")
		resp, err := client.Do(req)

		if err != nil {
			errorChan <- err
			return
		}

		defer resp.Body.Close()

		if resp.Body != nil {

			body, err := io.ReadAll(resp.Body)

			if err != nil {
				errorChan <- err
				return
			}

			var data interface{}

			err = json.Unmarshal(body, &data)

			if err != nil {
				errorChan <- err
				return
			}

			sqdb.Config.InfoChan <- fmt.Sprintf("we sent the following push notificaiton: %v", data)

		}
	}()

	err := <-errorChan

	if err != nil {
		return err
	}

	return nil

}
