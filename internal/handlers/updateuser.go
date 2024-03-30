package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/elkcityhazard/remind-me/internal/dbrepo/sqldbrepo"
	cerrors "github.com/elkcityhazard/remind-me/internal/errors"
	"github.com/go-chi/chi/v5"
)

func UpdateUser(w http.ResponseWriter, r *http.Request) {

	type tmpPassword struct {
		ID         *int64    `json:"id"`
		Hash       []byte    `json:"-"`
		Salt       []byte    `json:"-"`
		Plaintext1 *string   `json:"plaintext_1"`
		Plaintext2 *string   `json:"plaintext_2"`
		UserID     *int64    `json:"user_id"`
		UpdatedAt  time.Time `json:"updated_at"`
		IsActive   *int      `json:"is_active"`
		Version    *int      `json:"version"`
	}

	type tmpPhoneNumber struct {
		ID        *int64    `json:"id"`
		Hash      []byte    `json:"-"`
		Plaintext *string   `json:"plaintext"`
		UserID    *int64    `json:"user_id"`
		UpdatedAt time.Time `json:"updated_at"`
		Version   *int      `json:"version"`
	}

	type tmpUser struct {
		ID          *int64         `json:"id"`
		Email       *string        `json:"email"`
		Password    tmpPassword    `json:"password"`
		PhoneNumber tmpPhoneNumber `json:"phone_number"`
		UpdatedAt   time.Time      `json:"updated_at"`
		Scope       *int           `json:"scope"`
		IsActive    *int           `json:"is_active"`
		Version     int            `json:"version"`
	}

	idParam, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)

	if err != nil {
		if err = utilWriter.ErrorJSON(w, r, "error", err.Error(), http.StatusBadRequest); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		return
	}

	body, err := io.ReadAll(r.Body)

	if err != nil {
		if err = utilWriter.ErrorJSON(w, r, "error", err.Error(), http.StatusBadRequest); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		return
	}

	dbConn := sqldbrepo.GetDatabaseConnection()

	fetchedUser, err := dbConn.GetUserById(idParam)

	if err != nil {
		if err = utilWriter.ErrorJSON(w, r, "error", err.Error(), http.StatusBadRequest); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		return
	}

	originalPW := fetchedUser.Password.Hash

	var tmp tmpUser

	tmp.ID = &idParam

	if tmp.ID == nil {

		err = errors.New("error parsing id param")

		if err = utilWriter.ErrorJSON(w, r, "error", err.Error(), http.StatusBadRequest); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		return
	}

	err = json.Unmarshal(body, &tmp)

	if err != nil {
		if err = utilWriter.ErrorJSON(w, r, "error", err.Error(), http.StatusBadRequest); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		return
	}

	customErrors := cerrors.NewErrors()

	// we need to check for nil pointer values, and construct an updated user based on the ID

	// update the password if a password has been provided
	if !strings.EqualFold(*tmp.Password.Plaintext1, *tmp.Password.Plaintext2) {
		customErrors.Add("plaintext_1", "passwords do not match")
		customErrors.Add("plaintext_2", "passwords do not match")
	} else {

		if len(*tmp.Password.Plaintext1) > 0 || len(*tmp.Password.Plaintext2) > 0 {

			passwordIsTheSame := utilWriter.VerifyArgonHash(*tmp.Password.Plaintext1, string(fetchedUser.Password.Hash))

			app.InfoChan <- fmt.Sprintf("Is the password the same? %v", passwordIsTheSame)

			if !passwordIsTheSame {
				fetchedUser.Password.Plaintext1 = *tmp.Password.Plaintext1
				fetchedUser.Password.Plaintext2 = *tmp.Password.Plaintext2

				hashedPassword := utilWriter.CreateArgonHash(fetchedUser.Password.Plaintext1)

				tmp.Password.Hash = []byte(hashedPassword)

				// update the vals in Password

				fetchedUser.Password.Hash = []byte(hashedPassword)
				fetchedUser.Password.UpdatedAt = time.Now()
				// note pass pasword into update op to check for version in db and only increment if it is still equal to the current val
			}

		}

	}

	// handle phone number

	if len(*tmp.PhoneNumber.Plaintext) > 0 && !strings.EqualFold(*tmp.PhoneNumber.Plaintext, fetchedUser.PhoneNumber.Plaintext) {

		fetchedUser.PhoneNumber.Plaintext = *tmp.PhoneNumber.Plaintext
		fetchedUser.PhoneNumber.UpdatedAt = time.Now()
	}

	originalActiveStatus := fetchedUser.IsActive

	if *tmp.IsActive != fetchedUser.IsActive {

		if *tmp.IsActive < 0 {
			*tmp.IsActive = 0
		}

		if *tmp.IsActive > 1 {
			*tmp.IsActive = 1
		}

		if *tmp.IsActive != originalActiveStatus {
			fetchedUser.IsActive = *tmp.IsActive
		}
	}

	if !strings.EqualFold(string(originalPW), string(fetchedUser.Password.Hash)) || !strings.EqualFold(fetchedUser.PhoneNumber.Plaintext, *tmp.PhoneNumber.Plaintext) || originalActiveStatus != *tmp.IsActive {
		version, err := dbConn.UpdateUser(fetchedUser)

		if err != nil {
			if err = utilWriter.ErrorJSON(w, r, "error", err.Error(), http.StatusBadRequest); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			return
		}

		dataMap := make(map[string]interface{})

		dataMap["id"] = idParam
		dataMap["version"] = version
		dataMap["user"] = fetchedUser

		if err = utilWriter.WriteJSON(w, r, "data", dataMap, http.StatusOK); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	} else {

		dataMap := make(map[string]interface{})

		dataMap["id"] = idParam
		dataMap["version"] = fetchedUser.Version
		dataMap["user"] = fetchedUser

		if err = utilWriter.WriteJSON(w, r, "data", dataMap, http.StatusOK); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	}

}
