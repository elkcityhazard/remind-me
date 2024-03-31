package main

import (
	"net/http"

	"github.com/elkcityhazard/remind-me/internal/config"
	"github.com/elkcityhazard/remind-me/internal/handlers"
	"github.com/elkcityhazard/remind-me/pkg/utils"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func routes(app *config.AppConfig) *chi.Mux {
	mux := chi.NewRouter()

	mux.Use(middleware.Logger)
	mux.Use(middleware.Recoverer)
	mux.Use(SessionLoad)

	mux.Mount("/api/v1", PublicRoutes(app))
	mux.Mount("/api/v1/users", UserRoutes())
	mux.Mount("/api/v1/reminders", ReminderRoutes())

	return mux
}

func UserRoutes() http.Handler {
	r := chi.NewRouter()

	r.Post("/add", handlers.InsertUser)
	r.Put("/activation", handlers.HandleActivation)

	//sub route to handle auth required resources
	r.Route("/protected", func(r chi.Router) {

		//r.Use(RequiresAuth)
		r.Get("/id/{id}", handlers.GetUserByID)
		r.Put("/id/{id}", handlers.UpdateUser)
		r.Delete("/id/{id}", handlers.DeleteUser)
		r.Get("/email/{email}", handlers.GetUserByEmail)

	})

	return r
}

func ReminderRoutes() http.Handler {
	r := chi.NewRouter()

	//r.Use(RequiresAuth)

	r.Post("/add", handlers.HandleInsertReminder)
	r.Get("/{id}", handlers.HandleGetUserReminders)
	return r
}

func PublicRoutes(app *config.AppConfig) http.Handler {
	util := utils.NewUtils(app)

	r := chi.NewRouter()

	r.Get("/login", handlers.HandleSignIn)
	r.Get("/logout", handlers.HandleLogoutUser)
	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		type pingStatus struct {
			Code int    `json:"status_code"`
			Msg  string `json:"status_msg"`
		}

		var statusMsg pingStatus

		statusMsg.Code = 200
		statusMsg.Msg = "Successful Ping"

		headers := make(http.Header)

		headers.Add("Content-Type", "application/json;charset=utf-8")

		err := util.WriteJSON(w, r, "payload", statusMsg, 200, headers)
		if err != nil {
			app.ErrorChan <- err
		}
	})

	return r
}
