package utils

import (
	"encoding/json"
	"net/http"

	"github.com/elkcityhazard/remind-me/cmd/internal/config"
)

var app *config.AppConfig

type Utilser interface{}

type Utils struct {
	app         *config.AppConfig
	maxFileSize int
}

func NewUtils(a *config.AppConfig) *Utils {
	app = a

	return &Utils{
		app:         app,
		maxFileSize: 1024 * 1024 * 1024 * 30,
	}
}

func (u *Utils) WriteJSON(w http.ResponseWriter, r *http.Request, envelope string, data interface{}) error {
	payload := make(map[string]interface{})

	payload[envelope] = data

	w.Header().Set("Content-Type", "application/json;encoding=utf-8;")

	w.WriteHeader(http.StatusOK)

	err := json.NewEncoder(w).Encode(payload)
	if err != nil {
		return err
	}

	return nil
}
