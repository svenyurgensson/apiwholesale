package middleware

import (
	"net/http"
	"errors"
	str "strings"
	"fmt"
	"time"

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

		c.Env["auth_customer"] = customer

		var emptyTokenTTL  = customer.TokenTTL.IsZero()
		var tokenRotten = !emptyTokenTTL && time.Now().After(customer.TokenTTL)

		if emptyTokenTTL || tokenRotten {
			new_token := customer.RenewToken()
			w.Header().Set("X-Renew-Token", new_token)
			s.Log.Info(fmt.Sprintf("[renewToken] client %X new token generated", customer.Id.Hex ))
		}
		w.Header().Set("X-Token-TTL", customer.TokenTTL.Format(time.RFC822))

		// Update lastSeenAt field to current time
		customer.RenewLastSeen()

		h.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}
