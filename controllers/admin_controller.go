package controllers

import (
	"net/http"

	"github.com/zenazn/goji/web"
)


func AdminEntry(c web.C, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")


	w.WriteHeader(http.StatusOK)
}
