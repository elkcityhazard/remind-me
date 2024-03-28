package main

import (
	"context"
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/elkcityhazard/remind-me/internal/config"
	"github.com/elkcityhazard/remind-me/internal/dbrepo/sqldbrepo"
	"github.com/elkcityhazard/remind-me/internal/handlers"
	"github.com/elkcityhazard/remind-me/internal/mailer"
)

var (
	app config.AppConfig
)

func main() {

	appInit()

}

func startHTTPServer(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	srv := &http.Server{
		Addr:    ":8080",
		Handler: routes(&app),
		TLSConfig: &tls.Config{
			MinVersion:               tls.VersionTLS13,
			PreferServerCipherSuites: true,
		},
	}

	// Start the server in a goroutine
	go func() {
		fmt.Println("starting server...")
		if err := srv.ListenAndServeTLS("server.crt", "server.key"); err != http.ErrServerClosed {
			log.Fatalf("Server startup failed: %v", err)
		}
	}()

	// Listen for the context to be canceled
	go func() {
		<-ctx.Done()
		// When the context is canceled, initiate a graceful shutdown of the server
		fmt.Println("finalizing shutdown")
		if err := srv.Shutdown(context.Background()); err != nil {
			log.Fatalf("Server shutdown failed: %v", err)
		}
		os.Exit(0)
	}()

}

func appInit() {
	app = config.NewAppConfig()
	parseFlags()
	mailHandler := mailer.New("localhost", 1025, "web", "password", "web@remind-me.com")
	app.Mailer = mailHandler
	app.Session = getSession()
	app.WG.Add(1)

	go app.Mailer.ListenForMail(&app.WG)

	handlers.NewHandlers(&app)

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

	app.WG.Add(1)
	go listenForErrors(app.ErrorChan, app.ErrorDoneChan)

	shutdownCtx, cancel := context.WithCancel(context.Background())

	app.WG.Add(1)

	go startHTTPServer(shutdownCtx, &app.WG)

	signalCh := make(chan os.Signal, 1)

	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)

	<-signalCh

	fmt.Println("Shutting down server")

	cancel()

	app.WG.Wait()

	app.ErrorDoneChan <- true
	app.Mailer.MailerDoneChan <- true
	defer app.DB.Close()

	close(app.ErrorChan)
	close(app.ErrorDoneChan)
	close(app.Mailer.MailerDoneChan)
	close(app.Mailer.MailerDataChan)
	close(app.Mailer.MailerErrorChan)

	fmt.Println("Shutdown Completed")
	os.Exit(1)
}

func listenForErrors(eChan <-chan error, eDoneChan <-chan bool) {
	defer app.WG.Done()
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
