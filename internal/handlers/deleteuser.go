package handlers

import (
	"net/http"
	"strconv"

	"github.com/elkcityhazard/remind-me/internal/dbrepo/sqldbrepo"
	"github.com/go-chi/chi/v5"
)

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	idParam, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)

	if err != nil {
		if err = utilWriter.ErrorJSON(w, r, "error", err.Error(), http.StatusBadRequest); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		return
	}

	err = sqldbrepo.GetDatabaseConnection().DeleteUser(idParam)

	if err != nil {
		if err = utilWriter.ErrorJSON(w, r, "error", err.Error(), http.StatusBadRequest); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		return
	}

	if err = utilWriter.WriteJSON(w, r, "message", "the user has been deleted", http.StatusOK); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}
