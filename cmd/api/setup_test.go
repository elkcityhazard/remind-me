package main

import (
	"os"
	"testing"

	"github.com/elkcityhazard/remind-me/cmd/internal/config"
)

func TestMain(m *testing.M) {
	app = config.NewAppConfig()
	//  Make this work, all you have to do is
	// M.Run will handle running tests after app setup
	os.Exit(m.Run())
}
