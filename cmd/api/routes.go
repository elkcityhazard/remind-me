package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func routes() *chi.Mux {
	mux := chi.NewRouter()

	mux.Use(middleware.Logger)
	mux.Use(middleware.Recoverer)

	mux.Mount("/api/v1", PingRouter())

	return mux
}

func PingRouter() http.Handler {
	r := chi.NewRouter()

	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		type pingStatus struct {
			Code int    `json:"status_code"`
			Msg  string `json:"status_msg"`
		}

		var statusMsg pingStatus

		statusMsg.Code = 200
		statusMsg.Msg = "Successful Ping"

		err := util.WriteJSON(w, r, "payload", statusMsg)
		if err != nil {
			app.ErrorChan <- err
		}
	})

	return r
}
