package config

import (
	"log"
	"os"
)

type Apper interface{}

var app *AppConfig

type AppConfig struct {
	IsProduction  bool
	DSN           string
	Port          string // ":8080"
	InfoLog       *log.Logger
	ErrorLog      *log.Logger
	ErrorChan     chan error
	ErrorDoneChan chan bool
}

func NewAppConfig() *AppConfig {
	return &AppConfig{
		InfoLog:       log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime),
		ErrorLog:      log.New(os.Stdout, "ERROR", log.Ldate|log.Ltime|log.Lshortfile),
		ErrorChan:     make(chan error),
		ErrorDoneChan: make(chan bool),
	}
}
