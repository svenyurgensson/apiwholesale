package models

import (
    "errors"
    "time"

    s "apiwholesale/system"

    "gopkg.in/mgo.v2/bson"
)

type Order struct {
    Id        bson.ObjectId  `json:"id,omitempty"         bson:"_id"`
    CreatedAt time.Time      `json:"createdAt,omitempty"  bson:"createdAt,omitempty"`
    UpdatedAt time.Time      `json:"updatedAt,omitempty"  bson:"updatedAt,omitempty"`
    CustomerId bson.ObjectId `json:"customerId"           bson:"customerId,omitempty"`
    Status    string         `json:"status"               bson:"status"`
    RawData   interface{}    `json:"rawData"              bson:"rawData"`
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
        c.Status    = "pending"
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

    err := coll.Find(bson.M{"customerId": c.Id, "_id": bson.ObjectIdHex(id)}).One(&order)

    return order, err
}

func GetCustomerOrders(c *Customer) ([]Order, error) {
    orders := []Order{}
    session := s.GetSession()
    defer session.Close()
    coll := session.DB(s.DB).C("orders")

    err := coll.Find(bson.M{"customerId": c.Id}).All(&orders)

    return orders, err
}

func DeleteCustomerOrder(c Customer, id string) error {
    session := s.GetSession()
    defer session.Close()
    coll := session.DB(s.DB).C("orders")

    if ! bson.IsObjectIdHex(id) {
        return errors.New("Wrong order id format")
    }

    return coll.Remove(bson.M{"customerId": c.Id, "_id": bson.ObjectIdHex(id)})
}

func GetOrdersCount(q bson.M) (int, error) {
    session := s.GetSession()
    defer session.Close()
    coll := session.DB(s.DB).C("orders")

    result, err := coll.Find(q).Count()
    return result, err
}



func GetOrders(q bson.M, skip, limit int) ([]Order, error) {
    session := s.GetSession()
    defer session.Close()
    coll := session.DB(s.DB).C("orders")

    result := []Order{}
    err := coll.Find(q).Limit(limit).Skip(skip).All(&result)

    return result, err
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
