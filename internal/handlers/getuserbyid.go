package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/elkcityhazard/remind-me/internal/models"
	"github.com/elkcityhazard/remind-me/pkg/utils"
)

func InsertUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		headers := make(http.Header)
		headers.Add("Allow", "POST")
		w.WriteHeader(http.StatusBadRequest)
		errorMsg := struct {
			code int
			msg  string
		}{
			code: 400,
			msg:  "This method is not allowed",
		}
		if err := utils.NewUtils(app).ErrorJSON(w, r, "error", errorMsg, 400); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		return
	}

	defer r.Body.Close()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		errorMsg := struct {
			code int
			msg  string
		}{
			code: 400,
			msg:  "error parsing response body",
		}
		if err := utils.NewUtils(app).ErrorJSON(w, r, "error", errorMsg, 400); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	var user models.User
	err = json.Unmarshal(body, &user)

	if err != nil {
		app.ErrorChan <- err
		errorMsg := struct {
			code int
			msg  string
		}{
			code: 400,
			msg:  "error parsing response body into json",
		}
		if err := utils.NewUtils(app).ErrorJSON(w, r, "error", errorMsg, 400); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		return
	}

	if err := utils.NewUtils(app).WriteJSON(w, r, "user", user, http.StatusOK); err != nil {
		http.Error(w, "error writing json to response writer", http.StatusInternalServerError)
		return
	}
}

func GetUserByID(w http.ResponseWriter, r *http.Request) {
}
