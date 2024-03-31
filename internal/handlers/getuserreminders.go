package handlers

import (
	"net/http"

	"github.com/elkcityhazard/remind-me/internal/dbrepo/sqldbrepo"
)

func HandleGetUserReminders(w http.ResponseWriter, r *http.Request) {

	id, err := GetIDFromRouteKey(r)

	if err != nil {
		if err := utilWriter.ErrorJSON(w, r, "error", err.Error(), http.StatusBadRequest); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		return
	}

	reminders, err := sqldbrepo.GetDatabaseConnection().GetUserRemindersByID(id)

	if err != nil {
		if err := utilWriter.ErrorJSON(w, r, "error", err.Error(), http.StatusBadRequest); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		return
	}

	if err := utilWriter.WriteJSON(w, r, "data", reminders, http.StatusOK); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}
