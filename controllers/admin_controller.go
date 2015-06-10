package controllers

import (
	"net/http"
	"encoding/json"

	"../models"
	s "../system"

	"gopkg.in/mgo.v2/bson"
	"github.com/zenazn/goji/web"
)


func AdminApplication(c web.C, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")


	w.WriteHeader(http.StatusOK)
}


// ~~~~~~~~~ Customers CRUD ~~~~~~~~~~ //


func AdminCustomersList(c web.C, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	resources, err := models.GetCustomers()
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	encoder := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	encoder.Encode(resources)
}



func AdminCustomerCreate(c web.C, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var resource models.Customer
	err := json.NewDecoder(r.Body).Decode(&resource)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var exist bool
	exist, err = models.ExistsCustomers(bson.M{"email": resource.Email})
	if err != nil || exist == true {
		http.Error(w, http.StatusText(http.StatusConflict), http.StatusConflict)
		return
	}

	cid, error := resource.Upsert()

	if error != nil {
		http.Error(w, error.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Location", "/v1/admin/customer/" + cid.Hex())

	w.WriteHeader(http.StatusCreated)
}


func AdminCustomerView(c web.C, w http.ResponseWriter, r *http.Request) {
	resource_id := c.URLParams["customer_id"]
	if ! bson.IsObjectIdHex(resource_id) {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	resource, error := models.GetCustomer(bson.M{"_id": bson.ObjectIdHex(resource_id)})
	if error != nil {
		http.Error(w, error.Error(), http.StatusNotFound)
		return
	}

	encoder := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	encoder.Encode(resource)
}

func AdminCustomerUpdate(c web.C, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	resource_id := c.URLParams["customer_id"]
	if ! bson.IsObjectIdHex(resource_id) {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	cid := bson.ObjectIdHex(resource_id)

	presents, error := models.ExistsCustomers(bson.M{"_id": cid})
	if error != nil || presents != true {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	var resource models.Customer

	err := json.NewDecoder(r.Body).Decode(&resource)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resource.Id = cid
	_, err = resource.Upsert()

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}


func AdminCustomerDelete(c web.C, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	resource_id := c.URLParams["customer_id"]
	if ! bson.IsObjectIdHex(resource_id) {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	cid := bson.ObjectIdHex(resource_id)

	presents, error := models.ExistsCustomers(bson.M{"_id": cid})
	if error != nil || presents != true {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	err := models.DeleteCustomer(bson.M{"_id": cid})

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}


// ~~~~~~~~  Order CRUD ~~~~~~~~~ //


func AdminOrdersList(c web.C, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	resources, err := models.GetOrders(bson.M{})
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	encoder := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	encoder.Encode(resources)
}

func AdminOrderView(c web.C, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")


	w.WriteHeader(http.StatusOK)
}

func AdminOrderUpdate(c web.C, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")


	w.WriteHeader(http.StatusOK)
}

func AdminOrderDelete(c web.C, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")


	w.WriteHeader(http.StatusOK)
}


func avoid(){
  s.DEBUG("void")
}
