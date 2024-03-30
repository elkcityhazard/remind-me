package handlers

import (
	"github.com/elkcityhazard/remind-me/internal/config"
	"github.com/elkcityhazard/remind-me/pkg/utils"
)

var app *config.AppConfig
var utilWriter *utils.Utils

func NewHandlers(a *config.AppConfig) {
	app = a
}

func PassUtilToHandlers(u *utils.Utils) {
	utilWriter = u
}
