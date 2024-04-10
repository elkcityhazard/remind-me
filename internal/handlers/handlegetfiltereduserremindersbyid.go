package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/elkcityhazard/remind-me/internal/dbrepo/sqldbrepo"
)

func HandleGetFilteredUserRemindersByID(w http.ResponseWriter, r *http.Request) {

	userID := int64(1)

	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))

	if err != nil {
		limit = 10
	}

	if limit == 0 {
		limit = 10
	}

	page, err := strconv.Atoi(r.URL.Query().Get("page"))

	if err != nil {
		page = 1
	}

	if page <= 0 {
		page = 1
	}

	pageSize := (page - 1) * limit

	reminders, err := sqldbrepo.GetDatabaseConnection().GetFilteredUserRemindersByID(userID, limit, pageSize)
	if err != nil {

		if err := utilWriter.WriteJSON(w, r, "data", err.Error(), http.StatusOK); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	if len(reminders) == 0 {
		err := errors.New("no results")
		if err := utilWriter.ErrorJSON(w, r, "error", err.Error(), http.StatusNotFound); err != nil {
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
