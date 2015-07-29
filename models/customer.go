package models

import (
	"crypto/rand"
	"fmt"
	"io"
	"time"
	str "strings"
	"errors"

	s "apiwholesale/system"
	"gopkg.in/mgo.v2/bson"
)

type (
	 Customer struct {
		 Id       bson.ObjectId `json:"id"       bson:"_id"`
		 Email    string        `json:"email"    bson:"email"`
		 Password string        `json:"password" bson:"password"`
		 Token    string        `json:"token"    bson:"token"`
		 TokenTTL time.Time     `json:"tokenTTL" bson:"tokenTTL"`

		 CreatedAt  time.Time   `json:"createdAt,omitempty"  bson:"createdAt,omitempty"`
		 UpdatedAt  time.Time   `json:"updatedAt,omitempty"  bson:"updatedAt,omitempty"`
		 LastSeenAt time.Time   `json:"lastSeenAt,omitempty" bson:"lastSeenAt,omitempty"`

		 Balance int            `json:"balanceTotal,omitempty" bson:"balanceTotal,omitempty"`

		 RawData   interface{}  `json:"rawData"  bson:"rawData"`
	 }
)

func (c *Customer) RenewLastSeen() {
	session := s.GetSession()
	defer session.Close()
	coll := session.DB(s.DB).C("customers")

	coll.Update(bson.M{"_id": c.Id}, bson.M{"$currentDate": bson.M{"lastSeenAt": true}})
}

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

func (c *Customer) Validate() error {

	if len(c.Email) < 1 {
		return errors.New("Wrong customer email!")
	}

	if len(c.Password) < 4 {
		return errors.New("Wrong customer password!")
	}
	return nil
}


func (c *Customer) Upsert() (bson.ObjectId, error) {
	session := s.GetSession()
	defer session.Close()
	coll := session.DB(s.DB).C("customers")

	if error := c.Validate(); error != nil {
		return bson.NewObjectId(), error
	}

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

	return c.Id, err
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

func GetCustomersCount(q bson.M) (int, error) {
	session := s.GetSession()
	defer session.Close()
	coll := session.DB(s.DB).C("customers")

	result, err := coll.Find(q).Count()
	return result, err
}

func GetCustomers(q bson.M, skip, limit int) ([]Customer, error) {
	session := s.GetSession()
	defer session.Close()
	coll := session.DB(s.DB).C("customers")

	result := []Customer{}
	err := coll.Find(q).Limit(limit).Skip(skip).All(&result)

	return result, err
}

func ExistsCustomers(q bson.M) (bool, error) {
	session := s.GetSession()
	defer session.Close()
	coll := session.DB(s.DB).C("customers")

	count, err := coll.Find(q).Count()
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func GetCustomer(q bson.M) (Customer, error) {
	session := s.GetSession()
	defer session.Close()
	coll := session.DB(s.DB).C("customers")

	result := Customer{}
	err := coll.Find(q).One(&result)

	return result, err
}


func DeleteCustomer(q bson.M) error {
	session := s.GetSession()
	defer session.Close()
	coll := session.DB(s.DB).C("customers")

	return coll.Remove(q)
}
