package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/elkcityhazard/remind-me/internal/dbrepo/sqldbrepo"
	"github.com/go-chi/chi/v5"
)

type tempReminder struct {
	ID      *int64     `json:"id"`
	Title   *string    `json:"title"`
	Content *string    `json:"content"`
	DueDate *time.Time `json:"due_date"`
}

func HandleUpdateReminder(w http.ResponseWriter, r *http.Request) {

	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)

	if err != nil {
		if err := utilWriter.ErrorJSON(w, r, "error", err.Error(), http.StatusBadRequest); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		return
	}

	var t tempReminder

	dbConn := sqldbrepo.GetDatabaseConnection()

	rmdr, err := dbConn.GetReminderByID(id)
	if err != nil {
		if err := utilWriter.ErrorJSON(w, r, "error", err.Error(), http.StatusBadRequest); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		return
	}

	err = json.NewDecoder(r.Body).Decode(&t)

	if err != nil {
		if err := utilWriter.ErrorJSON(w, r, "error", err.Error(), http.StatusBadRequest); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		return
	}

	if t.Title != nil {
		rmdr.Title = *t.Title
	}

	if t.Content != nil {
		rmdr.Content = *t.Content
	}

	if t.DueDate != nil {
		rmdr.DueDate = *t.DueDate
	}

	rmdr.UpdatedAt = time.Now()

	updatedReminder, err := dbConn.UpdateReminder(rmdr)

	if err != nil {
		if err := utilWriter.ErrorJSON(w, r, "error", err.Error(), http.StatusBadRequest); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		return
	}
	if err := utilWriter.WriteJSON(w, r, "data", updatedReminder, http.StatusOK); err != nil {
		if err := utilWriter.ErrorJSON(w, r, "error", err.Error(), http.StatusBadRequest); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		return
	}

}
