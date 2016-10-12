package controllers

import (
    "net/http"
    "encoding/json"
    "fmt"
    "time"

    s "apiwholesale/system"

    "apiwholesale/models"
    "github.com/zenazn/goji/web"
)

type Credentials struct {
    Email     string `json:"email"`
    Password  string `json:"password"`
}

type TokenResource struct {
    Token     string `json:"token"`
}


func SessionCreate(c web.C, w http.ResponseWriter, r *http.Request) {
    var cred Credentials
    err := json.NewDecoder(r.Body).Decode(&cred)

    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        s.Log.Err(fmt.Sprintf("[error] session create: %s", err.Error()))
        return
    }

    if cred.Email == "" || cred.Password == "" {
        http.Error(w, "Bad request", http.StatusBadRequest)
        s.Log.Err(fmt.Sprintf("[error] session create: bad credentials %v", cred))
        return
    }

    customer, err := models.GetCustomerByCredentials(cred.Email, cred.Password)

    if err != nil {
        http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
        s.Log.Err(fmt.Sprintf("[error] session create: not found  %v, %s", cred, err.Error()))
        return
    }

    if customer.Token == "" || time.Now().After(customer.TokenTTL) {
        customer.RenewToken()

        if _, err = customer.Upsert(); err != nil {
            http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
            s.Log.Err(fmt.Sprintf("[error] session create: token update %s", err.Error()))
            return
        }
    }

    encoder := json.NewEncoder(w)
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)

    encoder.Encode(TokenResource{customer.Token})
}


func SessionDelete(c web.C, w http.ResponseWriter, r *http.Request) {
    customer := c.Env["auth_customer"].(models.Customer)
    customer.Token = ""
    customer.TokenTTL = time.Now()

    if _, err := customer.Upsert(); err != nil {
        http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
        s.Log.Err(fmt.Sprintf("[error] session delete: %s", err.Error()))
        return
    }

    w.WriteHeader(http.StatusNoContent)
}
