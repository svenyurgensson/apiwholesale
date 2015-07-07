package middleware

import (
	"net/http"
	"errors"
	str "strings"
	"fmt"

	"apiwholesale/models"

	s "apiwholesale/system"

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
			s.Log.Err(fmt.Sprintf("[error] client authorization failed: %v", r))
			return
		}

		customer, err := models.GetCustomerByToken(parts[1])
		if err != nil {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			s.Log.Err(fmt.Sprintf("[error] client authorization failed: %v", r))
			return
		}

		s.Log.Err(fmt.Sprintf("authorization failed: %v", r))

		c.Env["auth_customer"] = customer

		// session := s.GetSession()
		// defer session.Close()
		// coll := session.DB(s.DB).C("customers")
		// c.Env["session"]    = session
		// c.Env["collection"] = coll

		// Update lastSeenAt field to current time
		customer.RenewLastSeen()

		h.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}
