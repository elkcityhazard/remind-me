package handlers

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/elkcityhazard/remind-me/internal/dbrepo/sqldbrepo"
	"github.com/elkcityhazard/remind-me/pkg/utils"
	"github.com/go-chi/chi/v5"
)

func GetUserByEmail(w http.ResponseWriter, r *http.Request) {

	email := chi.URLParam(r, "email")

	if email == "" {
		if err := utils.NewUtils(app).ErrorJSON(w, r, "error", "invalid input", 400); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		return
	}

	dbRepo := sqldbrepo.GetDatabaseConnection()

	user, err := dbRepo.GetUserByEmail(email)

	if err != nil {
		if err == sql.ErrNoRows {
			err = errors.New("something went wrong, please try different input")
		}
		if err := utils.NewUtils(app).ErrorJSON(w, r, "error", err.Error(), 400); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		return
	}

	if err := utils.NewUtils(app).WriteJSON(w, r, "user", user, 200); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

}
