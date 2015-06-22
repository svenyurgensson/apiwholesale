package middleware

import (
	"net/http"
	"errors"
	str "strings"

	"apiwholesale/models"

	"github.com/zenazn/goji/web"
)

var (
	ErrorCustomerNotFound error = errors.New("Customer not found!")
)

func TokenAuth(c *web.C, h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("X-Middleware", "TokenAuth")

		t := r.Header.Get("Authorization")
		parts := str.Split(t, ":")

		if len(parts) < 2 {
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			return
		}

		customer, err := models.GetCustomerByToken(parts[1])
		if err != nil {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		c.Env["auth_customer"] = customer

		h.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}
