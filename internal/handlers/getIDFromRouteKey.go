package handlers

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func GetIDFromRouteKey(r *http.Request) (int64, error) {

	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)

	if err != nil {
		return 0, err
	}

	return id, nil

}
