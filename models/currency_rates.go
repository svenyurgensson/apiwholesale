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
		 CreatedAt  time.Time   `json:"createdAt,omitempty"  bson:"createdAt,omitempty"`
	 }
)

func GetLatestRate() (CurrencyRate, error) {
	session := s.GetSession()
	defer session.Close()
	coll := session.DB(s.DB).C("currencyRates")

	result := CurrencyRate{}
	err := coll.Find(bson.M{}).
		Sort("CreatedAt").
		Select(bson.M{"_id": 1, "rate": 1, "createdAt": 1}).
		One(&result)

	return result, err
}
