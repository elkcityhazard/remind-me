package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/elkcityhazard/remind-me/internal/dbrepo/sqldbrepo"
	"github.com/go-chi/chi/v5"
)

func HandleUpdateScheduleByID(w http.ResponseWriter, r *http.Request) {

	reminderID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)

	if err != nil {

		log.Println(err)
		if err := utilWriter.ErrorJSON(w, r, "error", err.Error(), http.StatusBadRequest); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		return
	}

	scheduleID, err := strconv.ParseInt(chi.URLParam(r, "scheduleID"), 10, 64)

	if err != nil {

		log.Println(err)
		if err := utilWriter.ErrorJSON(w, r, "error", err.Error(), http.StatusBadRequest); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		return
	}

	type scheduleUpdate struct {
		ReminderID   *int64     `json:"reminder_id"`
		DispatchTime *time.Time `json:"dispatch_time"`
		IsProcessed  *bool      `json:"is_processed"`
	}

	var su scheduleUpdate

	su.ReminderID = &reminderID

	err = json.NewDecoder(r.Body).Decode(&su)

	if err != nil {

		log.Println(err)
		if err := utilWriter.ErrorJSON(w, r, "error", err.Error(), http.StatusBadRequest); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		return
	}

	schedule, err := sqldbrepo.GetDatabaseConnection().GetScheduleByID(scheduleID)

	if err != nil {

		log.Println(err)
		if err := utilWriter.ErrorJSON(w, r, "error", err.Error(), http.StatusBadRequest); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		return
	}

	if su.ReminderID != nil {
		schedule.ReminderID = *su.ReminderID
	}

	if su.DispatchTime != nil {
		schedule.DispatchTime = *su.DispatchTime
	}

	if su.IsProcessed != nil {

		if *su.IsProcessed {
			schedule.IsProcessed = 1
		} else {
			schedule.IsProcessed = 0
		}

	}

	schedule.ID = scheduleID

	schedule.UpdatedAt = time.Now()

	schedule, err = sqldbrepo.DBRepo.UpdateScheduleByID(schedule)

	if err != nil {

		log.Println(err)
		if err := utilWriter.ErrorJSON(w, r, "error", err.Error(), http.StatusBadRequest); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		return
	}

	if err := utilWriter.WriteJSON(w, r, "data", schedule, http.StatusOK); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}
