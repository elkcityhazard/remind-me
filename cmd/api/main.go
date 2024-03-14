package main

import (
	"crypto/tls"
	"log"
	"net/http"

	"github.com/elkcityhazard/remind-me/cmd/internal/config"
	"github.com/elkcityhazard/remind-me/cmd/pkg/utils"
)

var (
	app  *config.AppConfig
	util *utils.Utils
)

func main() {
	app = config.NewAppConfig()

	util = utils.NewUtils(app)

	go listenForErrors(app.ErrorChan, app.ErrorDoneChan)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: routes(),
		TLSConfig: &tls.Config{
			MinVersion:               tls.VersionTLS13,
			PreferServerCipherSuites: true,
		},
	}

	app.InfoLog.Printf("Starting server on %s\n", srv.Addr)

	err := srv.ListenAndServeTLS("localhost.crt", "localhost.key")

	log.Fatal(err)
}

func listenForErrors(eChan <-chan error, eDoneChan <-chan bool) {
	for {
		select {
		case err := <-eChan:
			log.Println(err)
		case <-eDoneChan:
			return
		}
	}
}
