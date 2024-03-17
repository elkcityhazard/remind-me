package sqldbrepo

import (
	"database/sql"

	"github.com/elkcityhazard/remind-me/cmd/internal/config"
	"github.com/elkcityhazard/remind-me/cmd/internal/models"
	_ "github.com/go-sql-driver/mysql"
)

type SQLDBRepo struct {
	Config *config.AppConfig
}

func NewSQLDBRepo(ac *config.AppConfig) *SQLDBRepo {
	return &SQLDBRepo{
		Config: ac,
	}
}

//  NewDatabaseConn establishes a new database connection
//  and sets the DB property on the appconfig.  It returns
//  an error if it fails.

func (sqdb *SQLDBRepo) NewDatabaseConn() (*sql.DB, error) {
	sqdb.Config.InfoLog.Println("Opening  database connection...")
	conn, err := sql.Open("mysql", sqdb.Config.DSN)
	if err != nil {
		sqdb.Config.ErrorChan <- err
		return nil, err
	}

	// we pass it back up to app config as well in case we need it later

	sqdb.Config.DB = conn

	return conn, nil
}

// InsertUser accepts a pointer declaration to a user, and inserts it into the database.
// It will return the User ID, and an error if there are any
func (sqdb *SQLDBRepo) InsertUser(*models.User) (int64, error) {
	return 0, nil
}
