package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/elkcityhazard/remind-me/internal/dbrepo/sqldbrepo"
	cerrors "github.com/elkcityhazard/remind-me/internal/errors"
	"github.com/elkcityhazard/remind-me/internal/models"
)

func HandleInsertReminder(w http.ResponseWriter, r *http.Request) {

	customErrors := cerrors.NewErrors()

	userID := app.SessionManager.GetInt64(r.Context(), "id")

	if userID == 0 {
		//customErrors.Add("user_id", "no userID in context")
		userID = 1
	}

	var payload *models.Reminder

	err := json.NewDecoder(r.Body).Decode(&payload)

	if err != nil {
		app.ErrorChan <- err
		if err := utilWriter.ErrorJSON(w, r, "error", err.Error(), http.StatusBadRequest); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		return
	}

	payload.UserID = 1

	if customErrors.IsValid() {

		insertedReminder, err := sqldbrepo.GetDatabaseConnection().InsertReminder(payload)

		if err != nil {
			app.ErrorChan <- err
			if err := utilWriter.ErrorJSON(w, r, "error", err.Error(), http.StatusBadRequest); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			return
		}

		if err = utilWriter.WriteJSON(w, r, "payload", insertedReminder, http.StatusOK); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	} else {
		if err = utilWriter.ErrorJSON(w, r, "errors", customErrors, http.StatusBadRequest); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

}
