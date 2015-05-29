package main

import (
    "fmt"
    "os"
    "time"

    "./system"
    "./routes"

    "github.com/zenazn/goji"
)


func main() {
    system.Boot_time = fmt.Sprintf(time.Now().Format(time.RFC3339))
    system.Hostname, _  = os.Hostname()

    // Add routes
    routes.Include()

    goji.Serve()
}


// https://godoc.org/github.com/zenazn/goji
// EXAMPLES:
// https://github.com/haruyama/golang-goji-sample
// https://github.com/hypebeast/goji-boilerplate
// http://attilaolah.eu/2014/09/10/json-and-struct-composition-in-go/
