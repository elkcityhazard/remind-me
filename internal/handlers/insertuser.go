package handlers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/mail"
	"strings"

	"github.com/elkcityhazard/remind-me/internal/dbrepo/sqldbrepo"
	cerrors "github.com/elkcityhazard/remind-me/internal/errors"
	"github.com/elkcityhazard/remind-me/internal/models"
	"github.com/elkcityhazard/remind-me/pkg/utils"
)

func InsertUser(w http.ResponseWriter, r *http.Request) {

	util := utils.NewUtils(app)

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
		log.Println(err)
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

	e := cerrors.NewErrors()

	var email = user.Email
	var password1 = user.Password.Plaintext1
	var password2 = user.Password.Plaintext2

	// do passwords match?

	if !strings.EqualFold(password1, password2) {
		e.Add("plaintext_1", "passwords do not match")
		e.Add("plaintext_2", "passwords do not match")
	}

	// is email valid?

	_, err = mail.ParseAddress(email)

	if err != nil {
		e.Add("email", "invalid email address")
	}

	// is phone number valid?

	isValidPhoneNumber := util.ValidatePhoneNumber(user.PhoneNumber.Plaintext)

	if !isValidPhoneNumber {
		e.Add("phone_number", "invalid phone number")
	}

	// salt and hash password

	encodedPW := util.CreateArgonHash(user.Password.Plaintext1)

	validate := util.VerifyArgonHash(user.Password.Plaintext1, encodedPW)

	if !validate {
		e.Add("password_1", "password does not validate")
		e.Add("password_2", "password does not validate")
	}

	user.Password.Hash = []byte(encodedPW)

	// nuke password after hash

	user.Password.Plaintext1 = ""
	user.Password.Plaintext2 = ""

	// check errors first

	if !e.IsValid() {
		if err := utils.NewUtils(app).ErrorJSON(w, r, "error", e, http.StatusBadRequest); err != nil {
			http.Error(w, "error writing json to response writer", http.StatusInternalServerError)
			return
		}
		return
	}

	dbrepo := sqldbrepo.NewSQLDBRepo(app)

	// insert user

	_, err = dbrepo.InsertUser(&user)

	if err != nil {
		if err := utils.NewUtils(app).ErrorJSON(w, r, "error", err.Error(), http.StatusBadRequest); err != nil {
			http.Error(w, "error writing json to response writer", http.StatusInternalServerError)
			return
		}
		return
	}
	if err := utils.NewUtils(app).WriteJSON(w, r, "user", user, http.StatusOK); err != nil {
		http.Error(w, "error writing json to response writer", http.StatusInternalServerError)
		return
	}
}
