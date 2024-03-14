package main

import (
	"os"
	"testing"

	"github.com/elkcityhazard/remind-me/cmd/internal/config"
	"github.com/elkcityhazard/remind-me/cmd/pkg/utils"
)

func TestMain(m *testing.M) {
	app = config.NewAppConfig()
	util = utils.NewUtils(app)

	os.Exit(m.Run())
}
