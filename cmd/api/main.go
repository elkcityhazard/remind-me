package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/elkcityhazard/remind-me/internal/config"
	"github.com/elkcityhazard/remind-me/internal/dbrepo/sqldbrepo"
	"github.com/elkcityhazard/remind-me/internal/handlers"
	"github.com/elkcityhazard/remind-me/internal/mailer"
	"github.com/elkcityhazard/remind-me/pkg/utils"
)

var (
	app        config.AppConfig
	utilWriter *utils.Utils
)

func main() {
	shutdownCtx, cancel := context.WithCancel(context.Background())
	setupApp()
	setupHTTPServer(shutdownCtx)
	setupMailer()
	setupPollScheduledReminders()
	listenForErrors()
	gracefulShutdown(cancel)
}

func setupApp() {
	app = config.NewAppConfig()
	parseFlags()
	utilWriter = utils.NewUtils(&app)
	handlers.NewHandlers(&app)

	dbConn, err := sqldbrepo.NewSQLDBRepo(&app).NewDatabaseConn()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	app.DB = dbConn
}

func setupHTTPServer(ctx context.Context) {
	srv := &http.Server{
		Addr:    ":8080",
		Handler: routes(&app),
	}

	// Start the server in a goroutine
	app.WG.Add(1)
	go func() {
		defer app.WG.Done()
		fmt.Println("starting server...")

		if err := srv.ListenAndServeTLS("server.crt", "server.key"); err != http.ErrServerClosed {
			log.Fatalf("Server startup failed: %v", err)
		}
	}()

	// Listen for the context to be canceled
	app.WG.Add(1)
	go func() {
		defer app.WG.Done()

		<-ctx.Done()
		// When the context is canceled, initiate a graceful shutdown of the server
		fmt.Println("finalizing shutdown")
		if err := srv.Shutdown(ctx); err != nil {
			log.Fatalf("Server shutdown failed: %v", err)
		}
	}()
}

func setupMailer() {
	mailHandler := mailer.New("localhost", 1025, "web", "password", "web@remind-me.com")
	app.Mailer = mailHandler
	app.WG.Add(1)
	go app.Mailer.ListenForMail(&app.WG)
}

func setupPollScheduledReminders() {
	app.WG.Add(1)
	go pollScheduledReminders()
}

func pollScheduledReminders() {
	defer app.WG.Done()

	ticker := time.NewTicker(time.Second * 15)
	defer ticker.Stop()

	for {
		select {
		case <-app.ReminderDoneChan:
			return
		case t := <-ticker.C:
			app.InfoChan <- fmt.Sprintf("Ticking: %v", t)
			reminders, err := sqldbrepo.GetDatabaseConnection().ProcessAllReminders()
			if err != nil {
				log.Fatalln(err)
				app.ErrorChan <- err
			}
			if len(reminders) > 0 {
				app.InfoChan <- fmt.Sprintf("Processed the following Scheduled Reminders: %v", reminders)
			}
			ticker.Reset(time.Second * 15)

		}
	}
}

func gracefulShutdown(cancel context.CancelFunc) {
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	<-signalCh

	cancel()

	fmt.Println("Shutting down server")

	app.ReminderDoneChan <- true
	app.ErrorDoneChan <- true
	app.Mailer.MailerDoneChan <- true

	app.WG.Wait()

	close(app.InfoChan)
	close(app.ErrorChan)
	close(app.ErrorDoneChan)
	close(app.Mailer.MailerDoneChan)
	close(app.Mailer.MailerDataChan)
	close(app.Mailer.MailerErrorChan)
	close(app.ReminderDoneChan)

	if err := app.DB.Close(); err != nil {
		log.Printf("Failed to close database connection: %v", err)
	}

	fmt.Println("Shutdown Completed")
	os.Exit(0)
	fmt.Println("extra message")
}

func listenForErrors() {
	app.WG.Add(1)

	go func() {
		defer app.WG.Done()

		for {
			select {
			case err := <-app.ErrorChan:
				app.ErrorLog.Println(err)
			case msg := <-app.InfoChan:
				app.InfoLog.Print(msg)
			case <-app.ErrorDoneChan:
				return
			}
		}
	}()
}

func parseFlags() {
	flag.StringVar(&app.DSN, "DSN", "", "the database source name to connect to the database")
	flag.BoolVar(&app.IsProduction, "IsProduction", false, "whether the app is in production or not")

	flag.Parse()
}
