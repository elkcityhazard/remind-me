package sqldbrepo

import (
	"testing"

	"github.com/elkcityhazard/remind-me/cmd/internal/config"
)

func Test_NewSQLDBRepo(t *testing.T) {

	var app = &config.AppConfig{}

	repo := NewSQLDBRepo(app)

	if repo == nil {
		t.Errorf("expected a new sql db repo, bot nil")
	}

}
