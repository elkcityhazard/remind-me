package config

import (
	"context"
	"database/sql"
	"log"
	"os"
	"sync"

	"github.com/alexedwards/scs/v2"
	"github.com/elkcityhazard/remind-me/internal/mailer"
)

type Apper interface{}

type AppConfig struct {
	IsProduction  bool
	DSN           string
	DB            *sql.DB
	Session       *scs.SessionManager
	Port          string // ":8080"
	InfoLog       *log.Logger
	ErrorLog      *log.Logger
	InfoChan      chan (string)
	ErrorChan     chan error
	ErrorDoneChan chan bool
	WG            sync.WaitGroup
	Context       context.Context
	Mailer        mailer.Mailer
}

// NewAppConfig returns an app config preloaded with a few necessary components
func NewAppConfig() AppConfig {

	return AppConfig{
		InfoLog:       log.New(os.Stdout, "INFO: -> ", log.Ldate|log.Ltime),
		ErrorLog:      log.New(os.Stdout, "ERROR: -> ", log.Ldate|log.Ltime|log.Lshortfile),
		ErrorChan:     make(chan error),
		InfoChan:      make(chan string),
		ErrorDoneChan: make(chan bool),
		WG:            sync.WaitGroup{},
		Context:       context.Background(),
	}

}
