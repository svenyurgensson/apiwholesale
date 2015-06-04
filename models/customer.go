package models

import (
	s "../system"
	"gopkg.in/mgo.v2/bson"

	"crypto/rand"
	"fmt"
	"io"
	"time"
	str "strings"
)

type (
	 Customer struct {
		 Id       bson.ObjectId `json:"id"       bson:"_id"`
		 Email    string        `json:"email"    bson:"email"`
		 Password string        `json:"password" bson:"password"`
		 Token    string        `json:"token"    bson:"token"`
		 TokenTTL time.Time     `json:"tokenTtl" bson:"tokenTtl"`

		 CreatedAt time.Time    `json:"created_at,omitempty"  bson:"created_at,omitempty"`
		 UpdatedAt time.Time    `json:"updated_at,omitempty"  bson:"updated_at,omitempty"`
		 LastSeenAt time.Time   `json:"lastSeen_at,omitempty"  bson:"lastSeen_at,omitempty"`

		 RawData   interface{}  `json:"raw_data"  bson:"raw_data"`
	 }
)

func (c *Customer) RenewToken() string {
	uuid := make([]byte, 16)
	n, err := io.ReadFull(rand.Reader, uuid)
	if n != len(uuid) || err != nil {
		panic(err)
	}
	// variant bits; see section 4.1.1
	uuid[8] = uuid[8]&^0xc0 | 0x80
	// version 4 (pseudo-random); see section 4.1.3
	uuid[6] = uuid[6]&^0xf0 | 0x40
	uuid_string := fmt.Sprintf("%x-%x-%x-%x-%x", uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:])

	c.Token = uuid_string

	// Set TTL Today + 3 month
	var ttl time.Duration = time.Hour * 24 * 31 * 3
	c.TokenTTL = time.Now().Add(ttl)

	return uuid_string
}

func (c *Customer) Upsert() error {
	session := s.GetSession()
	defer session.Close()
	coll := session.DB(s.DB).C("customers")

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

//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

func GetCustomerByToken(token string) (Customer, error) {
	session := s.GetSession()
	defer session.Close()
	coll := session.DB(s.DB).C("customers")
	result := Customer{}

	tk := str.TrimSpace(token)
	err := coll.Find(bson.M{"token": tk}).One(&result)

	return result, err
}

func GetCustomerByCredentials(email string, password string) (Customer, error) {
	session := s.GetSession()
	defer session.Close()
	coll := session.DB(s.DB).C("customers")
	result := Customer{}

	e, p := str.TrimSpace(email), str.TrimSpace(password)
	err := coll.
		Find(bson.M{"email": e, "password": p}).
		One(&result)

	return result, err
}
