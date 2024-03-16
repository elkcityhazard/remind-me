package utils

import (
	"encoding/json"
	"net/http"

	"github.com/elkcityhazard/remind-me/cmd/internal/config"
)

var app *config.AppConfig

type Utilser interface {
	WriteJSON(w http.ResponseWriter, r *http.Request, envelope string, data interface{}) error
	ErrorJSON(w http.ResponseWriter, r *http.Request, enveloper string, data interface{}, statusCode int, headers ...http.Header) error
}

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

// WriteJSON takes in a responseWriter, request, enveolor, and data and write json to the response writer.
// it can return a potential error
// to pass in headers use headers := make(http.Header)
func (u *Utils) WriteJSON(w http.ResponseWriter, r *http.Request, envelope string, data interface{}, statusCode int, headers ...http.Header) error {
	payload := make(map[string]interface{})

	payload[envelope] = data

	w.Header().Set("Content-Type", "application/json;encoding=utf-8;")

	if len(headers) > 0 {

		for _, header := range headers {
			for key, value := range header {
				for _, v := range value {
					w.Header().Add(key, v)
				}
			}
		}

	}

	w.WriteHeader(statusCode)

	return json.NewEncoder(w).Encode(payload)

}

//  ErrorJSON is provided a response writer to write to, a pointer to a request, an envelope string, and some data
//  and returns an error json output.  It can potentially return an error

func (u *Utils) ErrorJSON(w http.ResponseWriter, r *http.Request, envelope string, data interface{}, statusCode int, headers ...http.Header) error {
	payload := map[string]interface{}{}

	payload[envelope] = data

	w.Header().Set("Content-Type", "application/json;encoding:utf-8;")

	for _, header := range headers {
		for key, values := range header {
			for _, value := range values {
				w.Header().Add(key, value)
			}
		}
	}

	w.WriteHeader(statusCode)

	return json.NewEncoder(w).Encode(payload)
}
