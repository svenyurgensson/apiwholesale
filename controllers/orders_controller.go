package controllers

import (
    "net/http"
    "encoding/json"

    "../models"

    s "../system"

    "gopkg.in/mgo.v2/bson"
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

func OrderGet(c web.C, w http.ResponseWriter, r *http.Request) {
    customer := c.Env["auth_customer"].(models.Customer)

    order_id := c.URLParams["order_id"]
    order, err := models.GetOrder(customer, order_id)

    if err != nil {
        http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
        return
    }

    encoder := json.NewEncoder(w)
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)

    encoder.Encode(order)
}



func OrderCreate(c web.C, w http.ResponseWriter, r *http.Request) {
    customer := c.Env["auth_customer"].(models.Customer)
    var order models.Order
    err := json.NewDecoder(r.Body).Decode(&order)

    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    order.CustomerId = customer.Id

    if err = order.Upsert(); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    w.WriteHeader(http.StatusCreated)
}

func OrderDelete(c web.C, w http.ResponseWriter, r *http.Request) {
    customer := c.Env["auth_customer"].(models.Customer)

    order_id := c.URLParams["order_id"]
    err := models.DeleteOrder(customer, order_id)

    if err != nil {
        http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
        return
    }

    w.WriteHeader(http.StatusNoContent)
}

func OrderUpdate(c web.C, w http.ResponseWriter, r *http.Request) {
    customer := c.Env["auth_customer"].(models.Customer)
    order_id := c.URLParams["order_id"]

    var order models.Order
    err := json.NewDecoder(r.Body).Decode(&order)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    if ! bson.IsObjectIdHex(order_id) {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    order.Id = bson.ObjectIdHex(order_id)
    order.CustomerId = customer.Id
    exists, error := models.Exists(bson.M{"customer_id": customer.Id, "_id": order.Id})

    if error != nil || exists == false {
        http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
        return
    }

    if err = order.Upsert(); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    w.WriteHeader(http.StatusOK)
}

func void(){
  s.DEBUG("void")
}