package controllers

import (
	"net/http"
	"encoding/json"
	"time"
	"fmt"

	"apiwholesale/models"

	s "apiwholesale/system"

	//"gopkg.in/mgo.v2/bson"
	"github.com/zenazn/goji/web"
)

type MyMessages struct {
	Multicast []string `json:"multicast"`
	Personal  []string `json:"personal"`
}

type MyResponse struct{
	Rate     float64    `json:"rate"`
	RateAt   time.Time  `json:"rateAt"`
	Messages MyMessages `json:"messages"`
}

func Me(c web.C, w http.ResponseWriter, r *http.Request) {
	//customer := c.Env["auth_customer"].(models.Customer)

	rate, err := models.GetLatestRate()
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		s.Log.Err(fmt.Sprintf("[error] GetLatestRate: %s", err.Error()))
		return
	}


	resource := &MyResponse{
		Rate: rate.Rate,
		RateAt: rate.CreatedAt,
		Messages: MyMessages{
			Multicast:   []string{"hello"},
			Personal: []string{"world"},
		},
	}


	encoder := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	encoder.Encode(resource)
}
