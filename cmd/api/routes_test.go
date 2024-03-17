package main

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
)

func TestRoutes(t *testing.T) {
	mux := chi.NewRouter()

	if mux == nil {
		t.Errorf("Expecting a mux router; got nil")
	}

	routes := routes(&app)

	chi.Walk(routes, func(method string, route string, handler http.Handler, middleware ...func(http.Handler) http.Handler) error {
		req := httptest.NewRequest(method, route, nil)
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("Expected status OK; got %v", rec.Code)
		}
		return nil
	})
}

func TestPingRouter(t *testing.T) {
	app.ErrorChan = make(chan error, 1)

	mockPing := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json;encoding=utf-8;")

		w.WriteHeader(http.StatusOK)

		fakeData := make(chan bool, 100)

		if err := json.NewEncoder(w).Encode(fakeData); err != nil {
			app.ErrorChan <- err
		}
	})

	ts := httptest.NewServer(mockPing)

	req, err := http.NewRequest(http.MethodGet, ts.URL+"/api/v1/ping/error", nil)
	if err != nil {
		t.Fatal(err)
	}

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		t.Log(err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(string(body))

	select {
	case err := <-app.ErrorChan:
		if err == nil {
			t.Error("expected error, but got none")
		} else {
			t.Log("Received expected error:", err)
		}
	case <-time.After(time.Second):
		t.Error("timeout waiting for error")
	}
}
