package controllers

import (
	"net/http"
	"encoding/json"

	"../models"

	"github.com/zenazn/goji/web"
)

func OrdersList(c web.C, w http.ResponseWriter, r *http.Request) {
	customer := c.Env["auth_customer"].(models.Customer)

	orders, err := models.GetOrders(&customer)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	encoder := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	encoder.Encode(orders)
}

func OrderCreate(c web.C, w http.ResponseWriter, r *http.Request) {
//    customer := c.Env["auth_customer"].(models.Customer)
}

func OrderDelete(c web.C, w http.ResponseWriter, r *http.Request) {
//    customer := c.Env["auth_customer"].(models.Customer)
}

func OrderUpdate(c web.C, w http.ResponseWriter, r *http.Request) {
//    customer := c.Env["auth_customer"].(models.Customer)
}
