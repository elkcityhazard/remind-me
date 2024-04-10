package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/elkcityhazard/remind-me/internal/dbrepo/sqldbrepo"
	"github.com/go-chi/chi/v5"
)

func HandleGetReminderByID(w http.ResponseWriter, r *http.Request) {

	reminderID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)

	if err != nil {
		if err := utilWriter.ErrorJSON(w, r, "error", err.Error(), http.StatusBadRequest); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		return
	}

	fmt.Println(reminderID)

	reminders, err := sqldbrepo.GetDatabaseConnection().GetReminderByID(reminderID)

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
