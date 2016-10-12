package models

import (
	"time"

	s "apiwholesale/system"
	"gopkg.in/mgo.v2/bson"
)

type (
	 CurrencyRate struct {
		 Id       bson.ObjectId `bson:"_id"`
		 Rate       float64     `json:"rate"  bson:"rate"`
		 CreatedAt  time.Time   `json:"createdAt,omitempty"  bson:"created_at,omitempty"`
	 }
)

func GetLatestRate() (CurrencyRate, error) {
	session := s.GetSession()
	defer session.Close()
	coll := session.DB(s.DB).C("currency_rates")

	result := CurrencyRate{}

	err := coll.Find(bson.M{}).
		Sort("-created_at").
		Select(bson.M{"_id": 1, "rate": 1, "created_at": 1}).
		One(&result)

	return result, err
}
