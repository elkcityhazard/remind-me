package handlers

import "github.com/elkcityhazard/remind-me/internal/config"

var app *config.AppConfig

func NewHandlers(a *config.AppConfig) {
	app = a
}
