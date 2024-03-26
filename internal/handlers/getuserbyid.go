package handlers

import (
	"net/http"
)

func GetUserByID(w http.ResponseWriter, r *http.Request) {

	w.Write([]byte("hello, world"))
}
