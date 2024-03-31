package handlers

import "net/http"

func HandleLogoutUser(w http.ResponseWriter, r *http.Request) {

	err := app.SessionManager.Destroy(r.Context())

	if err != nil {
		if err = utilWriter.ErrorJSON(w, r, "error", "error parsing body", http.StatusBadRequest); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		return
	}

	if err := utilWriter.WriteJSON(w, r, "message", "successfully logged out", http.StatusOK); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
