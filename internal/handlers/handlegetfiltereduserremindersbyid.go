package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/elkcityhazard/remind-me/internal/dbrepo/sqldbrepo"
)

func HandleGetFilteredUserRemindersByID(w http.ResponseWriter, r *http.Request) {

	userID := int64(1)

	limit := r.URL.Query().Get("limit")

	if limit == "" {
		limit = "10"
	}

	offset := r.URL.Query().Get("offset")

	if offset == "" {
		offset = "0"
	}

	page := r.URL.Query().Get("page")

	if page == "" {
		page = "1"
	}

	queryPage := ConvertStringToInt(page)

	queryOffset := queryPage - 1

	switch true {
	case queryPage <= 1:
		queryOffset = 0
	case queryPage > 1:
		queryOffset = queryPage * ConvertStringToInt(limit)
	}

	reminders, err := sqldbrepo.GetDatabaseConnection().GetFilteredUserRemindersByID(userID, ConvertStringToInt(limit), queryOffset)

	if err != nil {

		if err := utilWriter.WriteJSON(w, r, "data", err, http.StatusOK); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	if len(reminders) == 0 {
		http.Redirect(w, r, fmt.Sprintf("https://localhost:8080/api/v1/reminders?page=%d&limit=%d", queryPage-1, ConvertStringToInt(limit)), http.StatusSeeOther)
		return
	}

	if err := utilWriter.WriteJSON(w, r, "data", reminders, http.StatusOK); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

func ConvertStringToInt(s string) int {

	if s == "0" || s == "" {
		return 1
	}

	i, err := strconv.ParseInt(s, 10, 64)

	if err != nil {
		return 1
	}

	if int(i) < 0 {
		return 1
	}

	return int(i)

}
