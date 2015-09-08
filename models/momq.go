package models

import (
    "time"

    s "apiwholesale/system"
    "gopkg.in/mgo.v2/bson"
)

type (
     MoMq struct {
         Id         bson.ObjectId  `bson:"_id"`
         Done       bool           `json:"done"                 bson:"done"`
         CreatedAt  time.Time      `json:"createdAt,omitempty"  bson:"createdAt,omitempty"`
         Topic      string         `json:"topic"                bson:"topic"`
         Data       interface{}    `json:"data"                 bson:"data"`
     }
     OrderCreated struct {
         OrderId    bson.ObjectId  `json:"orderId" bson:"orderId"`
     }
)



func NotifyMq(topic string, data interface{}) error {
    session := s.GetSession()
    defer session.Close()
    coll := session.DB(s.DB).C("mq_jobs")

    c := MoMq{}

    c.Id        = bson.NewObjectId()
    c.Done      = false
    c.CreatedAt = time.Now()
    c.Topic     = topic
    c.Data      = data

    return coll.Insert(c)
}

func OrderCreatedNotify(o Order) error {
    oc := OrderCreated{}
    oc.OrderId = o.Id

    return NotifyMq("customer.order.created", oc)
}
