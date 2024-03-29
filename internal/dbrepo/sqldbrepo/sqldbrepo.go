package sqldbrepo

import (
	"database/sql"

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

	// we pass it back up to app config as well in case we need it later

	sqdb.Config.DB = conn

	DBRepo = sqdb

	return conn, nil
}

func GetDatabaseConnection() *SQLDBRepo {
	return DBRepo
}
