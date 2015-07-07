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
	customer := c.Env["auth_customer"].(models.Customer)

	since, e := time.Parse(time.RFC3339, c.URLParams["since"])
	if e != nil {
		since = time.Now().AddDate(0, 0, -1)
	}

	rate, err := models.GetLatestRate()
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		s.Log.Err(fmt.Sprintf("[error] GetLatestRate: %s", err.Error()))
		return
	}

	var glob_messages []string
	var pers_messages []string

	gm, err1 := models.GetMulticastMessagesSince(since)
	if err1 != nil {
		s.Log.Err(fmt.Sprintf("[error] GetMulticastMessagesSince: %s", err.Error()))
	} else {
		for _, c := range gm {
			glob_messages = append(glob_messages, fmt.Sprintf("%s : %s", c.CreatedAt.Format(time.RFC3339), c.Message))
		}
	}

	pm, err2 := models.GetDirectMessagesSince(customer, since)
	if err2 != nil {
		s.Log.Err(fmt.Sprintf("[error] GetPersonalMessagesSince: %s", err.Error()))
	} else {
		for _, c := range pm {
			pers_messages = append(pers_messages, fmt.Sprintf("%s : %s", c.CreatedAt.Format(time.RFC3339), c.Message))
		}
	}

	resource := &MyResponse{
		Rate: rate.Rate,
		RateAt: rate.CreatedAt,
		Messages: MyMessages{
			Multicast: glob_messages,
			Personal:  pers_messages,
		},
	}


	encoder := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	encoder.Encode(resource)
}
