package middleware

import (
    "net/http"
    "encoding/base64"
    "strings"
    "../system"

    "github.com/zenazn/goji/web"
)

// SuperSecure is HTTP Basic Auth middleware for super-secret admin page. Shhhh!
func SuperSecure(c *web.C, h http.Handler) http.Handler {
    fn := func(w http.ResponseWriter, r *http.Request) {

        auth := r.Header.Get("Authorization")

        if !strings.HasPrefix(auth, "Basic ") {
            pleaseAuth(w)
            return
        }

        password, err := base64.StdEncoding.DecodeString(auth[6:])

        if err != nil || string(password) != system.AdminCredentials {
            pleaseAuth(w)
            return
        }

        c.Env["Admin"] = true
        h.ServeHTTP(w, r)
    }
    return http.HandlerFunc(fn)
}

func pleaseAuth(w http.ResponseWriter) {
    w.Header().Set("WWW-Authenticate", `Basic realm="Gritter"`)
    w.WriteHeader(http.StatusUnauthorized)
    w.Write([]byte("Wrong or absent login/password!\n"))
}
