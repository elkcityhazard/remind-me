package main

import "net/http"

func RequiredAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		idExists := app.SessionManager.Exists(r.Context(), "id")

		if !idExists {
			http.Redirect(w, r, "/api/v1/error", http.StatusSeeOther)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func SessionLoad(next http.Handler) http.Handler {
	return app.SessionManager.LoadAndSave(next)
}
