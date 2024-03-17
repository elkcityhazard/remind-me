package sqldbrepo

import (
	"testing"

	"github.com/elkcityhazard/remind-me/cmd/internal/config"
)

func Test_NewSQLDBRepo(t *testing.T) {
	app := &config.AppConfig{}

	repo := NewSQLDBRepo(app)

	if repo == nil {
		t.Errorf("expected a new sql db repo, bot nil")
	}
}

func Test_NewDatabaseConn(t *testing.T) {
}

func Test_InsertUser(t *testing.T) {
}
