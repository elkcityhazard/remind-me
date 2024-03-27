package config

import (
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
	ErrorChan     chan error
	ErrorDoneChan chan bool
	WG            sync.WaitGroup
	Mailer        mailer.Mailer
}

// NewAppConfig returns an app config preloaded with a few necessary components
func NewAppConfig() AppConfig {
	return AppConfig{
		InfoLog:       log.New(os.Stdout, "INFO: -> ", log.Ldate|log.Ltime),
		ErrorLog:      log.New(os.Stdout, "ERROR: -> ", log.Ldate|log.Ltime|log.Lshortfile),
		ErrorChan:     make(chan error),
		ErrorDoneChan: make(chan bool),
		WG:            sync.WaitGroup{},
		Mailer:        mailer.New("localhost", 1025, "web@remind-me.com", "", ""),
	}
}
