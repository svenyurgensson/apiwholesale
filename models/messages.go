package models

import (
	"time"

	s "apiwholesale/system"
	"gopkg.in/mgo.v2/bson"
)

type (
	 Message struct {
		 Id         bson.ObjectId  `bson:"_id"`
		 Type       string         `json:"type"  bson:"type"`
		 Message    string         `json:"message"  bson:"message"`
		 ProducerId bson.ObjectId  `json:"producerId,omitempty" bson:"producerId,omitempty"`
		 RecipientId bson.ObjectId `json:"recipientId,omitempty" bson:"recipientId,omitempty"`
		 CreatedAt  time.Time      `json:"createdAt,omitempty"  bson:"createdAt,omitempty"`
	 }
)

func GetMulticastMessagesSince(timestamp time.Time) ([]Message, error) {
	session := s.GetSession()
	defer session.Close()
	coll := session.DB(s.DB).C("messages")

	result := []Message{}

	err := coll.Find(bson.M{"type": "multicast", "createdAt": bson.M{"$gte": timestamp}}).
		Sort("CreatedAt").
		Select(bson.M{"_id": 1, "message": 1, "createdAt": 1}).
		Limit(20).
		All(&result)

	return result, err
}

func GetDirectMessagesSince(c Customer, timestamp time.Time) ([]Message, error) {
	session := s.GetSession()
	defer session.Close()
	coll := session.DB(s.DB).C("messages")

	result := []Message{}

	err := coll.Find(bson.M{"type": "direct", "recipientId": c.Id, "createdAt": bson.M{"$gte": timestamp}}).
		Sort("CreatedAt").
		Select(bson.M{"_id": 1, "message": 1, "createdAt": 1}).
		Limit(20).
		All(&result)

	return result, err
}
