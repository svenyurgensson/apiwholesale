package main

import (
		"fmt"
		"net/http"

		"golang.org/x/crypto/bcrypt"
		"github.com/zenazn/goji"
		"github.com/zenazn/goji/web"
)

const Author = "Yury Batenko"
const Version = "1.8 / 2014-09-07"
const ApiVersion = "v1"


func ping(c web.C, w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, %s!", c.URLParams["name"])
}

func session(c web.C, w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "session created!")
}

func main() {
	goji.Get("/" + ApiVersion+ "/ping/:name", ping)

	goji.Post("/" + ApiVersion+ "/session", session)

	goji.Serve()
}


// https://godoc.org/github.com/zenazn/goji
// https://github.com/haruyama/golang-goji-sample
