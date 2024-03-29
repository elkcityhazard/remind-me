package sqldbrepo

import (
	"testing"

	"github.com/elkcityhazard/remind-me/internal/config"
)

func Test_GetUserByEmail(t *testing.T) {

	var mockApp = config.NewAppConfig()

	mockApp.DSN = ""

	mockStore := NewSQLDBRepo(&mockApp)

	mockConn, _ := mockStore.NewDatabaseConn()

	mockStore.Config = app
	mockStore.Config.DB = mockConn

	user, err := mockStore.GetUserByEmail("mario@google.com")

	if err != nil {
		t.Log(err)
	}

	if user == nil {
		t.Error("Expected user, but got none")
	}

	t.Log(user)

}
