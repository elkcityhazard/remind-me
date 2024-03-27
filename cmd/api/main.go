package main

import (
	"crypto/tls"
	"flag"
	"log"
	"net/http"

	"github.com/elkcityhazard/remind-me/internal/config"
	"github.com/elkcityhazard/remind-me/internal/dbrepo/sqldbrepo"
	"github.com/elkcityhazard/remind-me/internal/handlers"
)

var (
	app config.AppConfig
)

func main() {
	app = config.NewAppConfig()

	parseFlags()
	app.Session = getSession()

	go app.Mailer.ListenForMail()

	handlers.NewHandlers(&app)

	go listenForErrors(app.ErrorChan, app.ErrorDoneChan)
	dbConn, err := sqldbrepo.NewSQLDBRepo(&app).NewDatabaseConn()
	if err != nil {
		app.ErrorChan <- err
	}

	app.DB = dbConn

	err = app.DB.Ping()
	if err != nil {
		app.ErrorChan <- err
	}

	handlers.NewHandlers(&app) // we are passing this to the handlers package.  we might want to upgrade it to an interface later

	srv := &http.Server{
		Addr:    ":8080",
		Handler: routes(&app),
		TLSConfig: &tls.Config{
			MinVersion:               tls.VersionTLS13,
			PreferServerCipherSuites: true,
		},
	}

	app.InfoLog.Printf("Starting server on %s\n", srv.Addr)

	log.Fatalln(srv.ListenAndServeTLS("server.crt", "server.key"))

}

func listenForErrors(eChan <-chan error, eDoneChan <-chan bool) {
	for {
		select {
		case err := <-eChan:
			app.ErrorLog.Print(err.Error())
		case <-eDoneChan:
			return
		}
	}
}

func parseFlags() {
	flag.StringVar(&app.DSN, "DSN", "", "the database source name to connect to the database")

	flag.Parse()
}
