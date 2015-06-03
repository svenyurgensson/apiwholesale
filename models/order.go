package models

import (
    s "../system"

//  "gopkg.in/mgo.v2"
    "gopkg.in/mgo.v2/bson"
)

type Order struct {
    Id        bson.ObjectId  `json:"id"          bson:"_id"`
    CreatedAt *s.Timestamp   `json:"created_at"  bson:"created_at"`
    CustomerId bson.ObjectId `json:"-"           bson:"customer_id"`
    RawData   string         `json:"raw_data"    bson:"raw_data"`
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

func GetOrders(c *Customer) ([]Order, error) {
    orders := []Order{}
    session := s.GetSession()
    defer session.Close()
    coll := session.DB(s.DB).C("orders")

    err := coll.Find(bson.M{"customer_id": c.Id}).All(&orders)

    return orders, err
}

func DeleteOrder(c *Customer, id bson.ObjectId) error {
    session := s.GetSession()
    defer session.Close()
    coll := session.DB(s.DB).C("orders")

    return coll.Remove(bson.M{"customer_id": c.Id, "_id": id})
}

// http://stevenwhite.com/building-a-rest-service-with-golang-3/
