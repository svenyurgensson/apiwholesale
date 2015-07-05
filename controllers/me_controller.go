package controllers

import (
	"net/http"
	"encoding/json"

	//"apiwholesale/models"

	// s "apiwholesale/system"

	// "gopkg.in/mgo.v2/bson"
	"github.com/zenazn/goji/web"
)

type MyMessages struct {
	Global   []string `json:"global"`
	Personal []string `json:"personal"`
}

type MyResponse struct{
	Rate     float64    `json:"rate"`
	Messages MyMessages `json:"messages"`
}

func Me(c web.C, w http.ResponseWriter, r *http.Request) {
	//customer := c.Env["auth_customer"].(models.Customer)

	resource := &MyResponse{
		Rate: 8.32,
		Messages: MyMessages{
			Global:   []string{"hello"},
			Personal: []string{"world"},
		},
	}


	encoder := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	encoder.Encode(resource)
}
