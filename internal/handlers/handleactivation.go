package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/elkcityhazard/remind-me/internal/dbrepo/sqldbrepo"
)

func HandleActivation(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")

	id, err := strconv.ParseInt(r.URL.Query().Get("id"), 10, 64)

	if err != nil {
		if err = utilWriter.ErrorJSON(w, r, "error", err.Error(), 400); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		return
	}

	user, err := sqldbrepo.GetDatabaseConnection().ActivateUser(token, id)

	if err != nil {
		fmt.Println(err)
		if err = utilWriter.ErrorJSON(w, r, "error", err.Error(), 400); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		return
	}

	app.SessionManager.Put(r.Context(), "id", user.ID)

	if err = utilWriter.WriteJSON(w, r, "user", user, http.StatusOK); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

}
