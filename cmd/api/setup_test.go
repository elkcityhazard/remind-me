package main

import (
	"os"
	"testing"

	"github.com/elkcityhazard/remind-me/internal/config"
	"github.com/elkcityhazard/remind-me/internal/handlers"
)

func TestMain(m *testing.M) {
	app = config.NewAppConfig()

	handlers.NewHandlers(&app)

	go listenForErrors(app.ErrorChan, app.ErrorDoneChan)
	//  Make this work, all you have to do is
	// M.Run will handle running tests after app setup
	os.Exit(m.Run())
}
