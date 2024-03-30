package handlers

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"

	"github.com/elkcityhazard/remind-me/internal/dbrepo/sqldbrepo"
	"github.com/go-chi/chi/v5"
)

func GetUserByID(w http.ResponseWriter, r *http.Request) {

	idParam, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		app.ErrorChan <- err
		if err := utilWriter.ErrorJSON(w, r, "error", err.Error(), 400); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		return
	}

	dbRepo := sqldbrepo.GetDatabaseConnection()

	user, err := dbRepo.GetUserById(idParam)
	if err != nil {

		if err == sql.ErrNoRows {
			err = errors.New("something went wrong, please try different input")
		}

		app.ErrorChan <- err

		if err := utilWriter.ErrorJSON(w, r, "error", err.Error(), 400); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		return
	}

	if err := utilWriter.WriteJSON(w, r, "user", user, 200); err != nil {
		app.ErrorChan <- err
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}
