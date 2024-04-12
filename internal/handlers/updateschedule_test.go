package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_HandleUpdateScheduleByID(t *testing.T) {

	var reminderID = 1

	var scheduleID = 3

	var data any

	p, err := json.Marshal(data)

	if err != nil {
		log.Fatalln(err)
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("/api/v1/%d/schedule/%d", reminderID, scheduleID), bytes.NewReader(p))

	if err != nil {
		log.Fatalln(err)
	}

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(HandleUpdateScheduleByID)

	handler.ServeHTTP(rr, req)

	expected := `{"some": "data"}`

	if rr.Body.String() != expected {
		t.Errorf("expected %s but got %s", expected, rr.Body.String())
	}

}
