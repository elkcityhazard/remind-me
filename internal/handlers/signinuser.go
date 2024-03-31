package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/mail"

	"github.com/elkcityhazard/remind-me/internal/dbrepo/sqldbrepo"
	cerrors "github.com/elkcityhazard/remind-me/internal/errors"
)

func HandleSignIn(w http.ResponseWriter, r *http.Request) {

	defer r.Body.Close()

	jsonErrors := cerrors.NewErrors()

	body, err := io.ReadAll(r.Body)

	if err != nil {
		if err = utilWriter.ErrorJSON(w, r, "error", "error parsing body", http.StatusBadRequest); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		return
	}

	type tmpUser struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var user tmpUser

	err = json.Unmarshal(body, &user)

	if err != nil {
		if err = utilWriter.ErrorJSON(w, r, "error", err.Error(), http.StatusBadRequest); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		return
	}

	if len(user.Email) == 0 {
		jsonErrors.Add("email", "must provide an email")
	}

	_, err = mail.ParseAddress(user.Email)

	if err != nil {
		jsonErrors.Add("email", "invalid email address")
	}

	if !jsonErrors.IsValid() {
		if err := utilWriter.ErrorJSON(w, r, "errors", jsonErrors, http.StatusBadRequest); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		return
	}

	fetchedUser, err := sqldbrepo.GetDatabaseConnection().GetUserByEmail(user.Email)

	if err != nil {
		app.ErrorChan <- err
		if err := utilWriter.ErrorJSON(w, r, "errors", "something went wrong", http.StatusBadRequest); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		return
	}

	passwordMatches := utilWriter.VerifyArgonHash(user.Password, string(fetchedUser.Password.Hash))

	if !passwordMatches {
		jsonErrors.Add("password_1", "invalid password input")
		jsonErrors.Add("password_2", "invalid password input")
		if err := utilWriter.ErrorJSON(w, r, "errors", jsonErrors, http.StatusBadRequest); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		return
	}

	if fetchedUser.IsActive == 0 {
		err := errors.New("invalid user")
		if err := utilWriter.ErrorJSON(w, r, "errors", err.Error(), http.StatusBadRequest); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		return
	}

	app.SessionManager.Put(r.Context(), "id", fetchedUser.ID)

	if err := utilWriter.WriteJSON(w, r, "user", fetchedUser, http.StatusOK); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}
