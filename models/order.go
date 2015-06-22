package models

import (
    s "../system"

    "errors"
    "time"
    "gopkg.in/mgo.v2/bson"
)

type Order struct {
    Id        bson.ObjectId  `json:"id,omitempty"          bson:"_id"`
    CreatedAt time.Time      `json:"created_at,omitempty"  bson:"created_at,omitempty"`
    UpdatedAt time.Time      `json:"updated_at,omitempty"  bson:"updated_at,omitempty"`
    CustomerId bson.ObjectId `json:"customer_id"           bson:"customer_id,omitempty"`
    RawData   interface{}    `json:"raw_data"              bson:"raw_data"`
}


func (c *Order) Upsert() error {
    session := s.GetSession()
    defer session.Close()
    coll := session.DB(s.DB).C("orders")

    var err error

    if ! c.Id.Valid() {
        c.Id = bson.NewObjectId()
        c.CreatedAt = time.Now()
        c.UpdatedAt = c.CreatedAt
        err = coll.Insert(c)
    } else {
        c.UpdatedAt = time.Now()
        err = coll.Update(bson.M{"_id": c.Id}, bson.M{"$set": c})
    }

    return err
}

func ExistsOrders(q bson.M) (bool, error) {
    session := s.GetSession()
    defer session.Close()
    coll := session.DB(s.DB).C("orders")

    count, err := coll.Find(q).Count()
    if err != nil {
        return false, err
    }

    return count > 0, nil
}

func GetCustomerOrder(c Customer, id string) (Order, error) {
    order := Order{}
    session := s.GetSession()
    defer session.Close()
    coll := session.DB(s.DB).C("orders")

    if ! bson.IsObjectIdHex(id) {
        return order, errors.New("Wrong order id format")
    }

    err := coll.Find(bson.M{"customer_id": c.Id, "_id": bson.ObjectIdHex(id)}).One(&order)

    return order, err
}

func GetCustomerOrders(c *Customer) ([]Order, error) {
    orders := []Order{}
    session := s.GetSession()
    defer session.Close()
    coll := session.DB(s.DB).C("orders")

    err := coll.Find(bson.M{"customer_id": c.Id}).All(&orders)

    return orders, err
}

func DeleteCustomerOrder(c Customer, id string) error {
    session := s.GetSession()
    defer session.Close()
    coll := session.DB(s.DB).C("orders")

    if ! bson.IsObjectIdHex(id) {
        return errors.New("Wrong order id format")
    }

    return coll.Remove(bson.M{"customer_id": c.Id, "_id": bson.ObjectIdHex(id)})
}


func GetOrders(q bson.M, skip, limit int) ([]Order, error) {
    orders := []Order{}
    session := s.GetSession()
    defer session.Close()
    coll := session.DB(s.DB).C("orders")

    err := coll.Find(q).Limit(limit).Skip(skip).All(&orders)

    return orders, err
}

func GetOrder(q bson.M) (Order, error) {
    order := Order{}
    session := s.GetSession()
    defer session.Close()
    coll := session.DB(s.DB).C("orders")

    err := coll.Find(q).One(&order)

    return order, err
}


func DeleteOrder(q bson.M) error {
    session := s.GetSession()
    defer session.Close()
    coll := session.DB(s.DB).C("orders")

    return coll.Remove(q)
}







// http://stevenwhite.com/building-a-rest-service-with-golang-3/
