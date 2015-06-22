package controllers

import (
    "net/http"
    "encoding/json"

    "apiwholesale/models"
    "fmt"

    s "apiwholesale/system"

    "gopkg.in/mgo.v2/bson"
    "github.com/zenazn/goji/web"
)


func OrdersList(c web.C, w http.ResponseWriter, r *http.Request) {
    customer := c.Env["auth_customer"].(models.Customer)

    resources, err := models.GetCustomerOrders(&customer)
    if err != nil {
        http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
        s.Log.Err(fmt.Sprintf("[error] orders list: %s", err.Error()))
        return
    }

    encoder := json.NewEncoder(w)
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)

    encoder.Encode(resources)
}

func OrderGet(c web.C, w http.ResponseWriter, r *http.Request) {
    customer := c.Env["auth_customer"].(models.Customer)

    order_id := c.URLParams["order_id"]
    resource, err := models.GetCustomerOrder(customer, order_id)

    if err != nil {
        http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
        s.Log.Err(fmt.Sprintf("[error] order not found: %s", err.Error()))
        return
    }

    encoder := json.NewEncoder(w)
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)

    encoder.Encode(resource)
}



func OrderCreate(c web.C, w http.ResponseWriter, r *http.Request) {
    customer := c.Env["auth_customer"].(models.Customer)
    var order models.Order
    err := json.NewDecoder(r.Body).Decode(&order)

    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        s.Log.Err(fmt.Sprintf("[error] order create: bad format %s", err.Error()))
        return
    }

    order.CustomerId = customer.Id

    if err = order.Upsert(); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        s.Log.Err(fmt.Sprintf("[error] order create: %s", err.Error()))
        return
    }

    w.WriteHeader(http.StatusCreated)
}

func OrderDelete(c web.C, w http.ResponseWriter, r *http.Request) {
    customer := c.Env["auth_customer"].(models.Customer)

    order_id := c.URLParams["order_id"]
    err := models.DeleteCustomerOrder(customer, order_id)

    if err != nil {
        http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
        s.Log.Err(fmt.Sprintf("[error] order delete: %s", err.Error()))
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
        s.Log.Err(fmt.Sprintf("[error] order update: %s", err.Error()))
        return
    }

    if ! bson.IsObjectIdHex(order_id) {
        http.Error(w, err.Error(), http.StatusBadRequest)
        s.Log.Err(fmt.Sprintf("[error] order update: bad order_id"))
        return
    }

    order.Id = bson.ObjectIdHex(order_id)
    order.CustomerId = customer.Id
    exists, error := models.ExistsOrders(bson.M{"customer_id": customer.Id, "_id": order.Id})

    if error != nil {
        http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
        s.Log.Err(fmt.Sprintf("[error] order update: %s", error.Error()))
        return
    }

    if exists == false {
        http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
        s.Log.Err(fmt.Sprintf("[error] order update: exists: %t", exists))
        return
    }

    if err = order.Upsert(); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        s.Log.Err(fmt.Sprintf("[error] order update: %s", err.Error()))
        return
    }

    w.WriteHeader(http.StatusOK)
}
