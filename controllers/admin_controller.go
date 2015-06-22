package controllers

import (
	"net/http"
	"encoding/json"
	"strconv"
	"fmt"

	"apiwholesale/models"
	s "apiwholesale/system"

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

	limit, el := strconv.Atoi(r.URL.Query().Get("limit"));
	if (el != nil) { limit = 100 }
	skip,  es := strconv.Atoi(r.URL.Query().Get("skip"));
	if (es != nil) { skip = 0 }

	type Response struct {
		Total int                      `json:"total"`
		Resources []models.Customer    `json:"resources"`
	}
	var response Response
	var err error

	if response.Total, err = models.GetCustomersCount(bson.M{}); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		s.Log.Err(fmt.Sprintf("[error] admin customers list: %s", err.Error()))
		return
	}

	if response.Resources, err = models.GetCustomers(bson.M{}, skip, limit); err != nil  {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		s.Log.Err(fmt.Sprintf("[error] admin customers list: %s", err.Error()))
		return
	}

	encoder := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	encoder.Encode(response)
}



func AdminCustomerCreate(c web.C, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var resource models.Customer
	err := json.NewDecoder(r.Body).Decode(&resource)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		s.Log.Err(fmt.Sprintf("[error] admin customers create: %s", err.Error()))
		return
	}

	var exist bool
	exist, err = models.ExistsCustomers(bson.M{"email": resource.Email})
	if err != nil || exist == true {
		http.Error(w, http.StatusText(http.StatusConflict), http.StatusConflict)
		s.Log.Err(fmt.Sprintf("[error] admin customers create: %s, exists: %t", err.Error(), exist))
		return
	}

	cid, error := resource.Upsert()

	if error != nil {
		http.Error(w, error.Error(), http.StatusBadRequest)
		s.Log.Err(fmt.Sprintf("[error] admin customers create: %s", err.Error()))
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
		s.Log.Err(fmt.Sprintf("[error] admin customers view: %s", error.Error()))
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
		s.Log.Err(fmt.Sprintf("[error] admin customers update: bad customer_id"))
		return
	}

	rid := bson.ObjectIdHex(resource_id)

	presents, error := models.ExistsCustomers(bson.M{"_id": rid})
	if error != nil || presents != true {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		s.Log.Err(fmt.Sprintf("[error] admin customers update: %s, presents: %t", error.Error(), presents))
		return
	}

	var resource models.Customer

	err := json.NewDecoder(r.Body).Decode(&resource)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		s.Log.Err(fmt.Sprintf("[error] admin customers update: %s", err.Error()))
		return
	}

	resource.Id = rid
	_, err = resource.Upsert()

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		s.Log.Err(fmt.Sprintf("[error] admin customers update: %s", err.Error()))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}


func AdminCustomerDelete(c web.C, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	resource_id := c.URLParams["customer_id"]
	if ! bson.IsObjectIdHex(resource_id) {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		s.Log.Err(fmt.Sprintf("[error] admin customers delete: bad customer_id"))
		return
	}

	rid := bson.ObjectIdHex(resource_id)

	presents, error := models.ExistsCustomers(bson.M{"_id": rid})
	if error != nil || presents != true {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		s.Log.Err(fmt.Sprintf("[error] admin customers delete: %s, present: %t", error.Error(), presents))
		return
	}

	err := models.DeleteCustomer(bson.M{"_id": rid})

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		s.Log.Err(fmt.Sprintf("[error] admin customers delete: %s", err.Error()))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}


// ~~~~~~~~  Order CRUD ~~~~~~~~~ //


func AdminOrdersList(c web.C, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	limit, el := strconv.Atoi(r.URL.Query().Get("limit"));
	if (el != nil) { limit = 100 }
	skip,  es := strconv.Atoi(r.URL.Query().Get("skip"));
	if (es != nil) { skip = 0 }

	type Response struct {
		Total int                   `json:"total"`
		Resources []models.Order    `json:"resources"`
	}
	var response Response
	var err error

	if response.Total, err = models.GetOrdersCount(bson.M{}); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		s.Log.Err(fmt.Sprintf("[error] admin orders list count: %s", err.Error()))
		return
	}

	if response.Resources, err = models.GetOrders(bson.M{}, skip, limit); err != nil  {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		s.Log.Err(fmt.Sprintf("[error] admin orders list: %s", err.Error()))
		return
	}

	encoder := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	encoder.Encode(response)
}

func AdminOrderView(c web.C, w http.ResponseWriter, r *http.Request) {
	resource_id := c.URLParams["order_id"]
	if ! bson.IsObjectIdHex(resource_id) {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		s.Log.Err(fmt.Sprintf("[error] admin order: bad order_id"))
		return
	}

	resource, error := models.GetOrder(bson.M{"_id": bson.ObjectIdHex(resource_id)})
	if error != nil {
		http.Error(w, error.Error(), http.StatusNotFound)
		s.Log.Err(fmt.Sprintf("[error] admin order: %s", error.Error()))
		return
	}

	encoder := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	encoder.Encode(resource)
}

func AdminOrderUpdate(c web.C, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	resource_id := c.URLParams["order_id"]
	if ! bson.IsObjectIdHex(resource_id) {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		s.Log.Err(fmt.Sprintf("[error] admin order update: bad order_id"))
		return
	}

	rid := bson.ObjectIdHex(resource_id)

	presents, error := models.ExistsOrders(bson.M{"_id": rid})
	if error != nil || presents != true {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		s.Log.Err(fmt.Sprintf("[error] admin order update: %s, exists: %t", error.Error(), presents))
		return
	}

	var resource models.Order

	err := json.NewDecoder(r.Body).Decode(&resource)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resource.Id = rid
	err = resource.Upsert()

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		s.Log.Err(fmt.Sprintf("[error] admin order update: %s", err.Error()))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func AdminOrderDelete(c web.C, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	resource_id := c.URLParams["order_id"]
	if ! bson.IsObjectIdHex(resource_id) {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		s.Log.Err(fmt.Sprintf("[error] admin order delete: bad order_id"))
		return
	}

	rid := bson.ObjectIdHex(resource_id)

	presents, error := models.ExistsOrders(bson.M{"_id": rid})
	if error != nil || presents != true {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		s.Log.Err(fmt.Sprintf("[error] admin order delete: %s, exists: %t", error.Error(), presents))
		return
	}

	err := models.DeleteOrder(bson.M{"_id": rid})

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		s.Log.Err(fmt.Sprintf("[error] admin order delete: %s", err.Error()))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}


func avoid(){
  s.DEBUG("void")
}
