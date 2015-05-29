package controllers

import (
	"fmt"
	"net/http"

	"github.com/zenazn/goji/web"
)


func SessionCreate(c web.C, w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "session created!\n")
}


func SessionDelete(c web.C, w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")

	fmt.Fprintf(w, "session deleted! %s\n", token)
}
