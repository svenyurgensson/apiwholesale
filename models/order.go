package models

import (
	"time"

	s "../system"

//  "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Order struct {
	Id        bson.ObjectId  `json:"id"          bson:"_id"`
	CreatedAt time.Time      `json:"createdAt"   bson:"createdAt"`
	CustomerId bson.ObjectId `json:"customer_id" bson:"customer_id"`
	RawData  string          `json:"raw_data"    bson:"raw_data"`
}


func (c *Order) Upsert() error {
	session := s.GetSession()
	defer session.Close()
	coll := session.DB(s.DB).C("orders")

	var err error

	if ! c.Id.Valid() {
		c.Id = bson.NewObjectId()
		err = coll.Insert(c)
	} else {
		err = coll.Update(bson.M{"_id": c.Id}, c)
	}

	return err
}

// func GetOrders(c *Customer) []Order {

// }

func DeleteOrder(c *Customer, id bson.ObjectId) error {
	session := s.GetSession()
	defer session.Close()
	coll := session.DB(s.DB).C("orders")

	return coll.Remove(bson.M{"customer_id": c.Id, "_id": id})
}

// http://stevenwhite.com/building-a-rest-service-with-golang-3/
