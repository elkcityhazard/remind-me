package main

import "net/http"

func RequiresAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		idExists := app.SessionManager.Exists(r.Context(), "id")

		if !idExists {
			if err := utilWriter.ErrorJSON(w, r, "error", "this is a protected resource, please login", http.StatusBadRequest); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			return
		}

		next.ServeHTTP(w, r)
	})
}

func SessionLoad(next http.Handler) http.Handler {
	return app.SessionManager.LoadAndSave(next)
}
