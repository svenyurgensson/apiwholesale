package models

import (
     s "apiwholesale/system"
    "gopkg.in/mgo.v2/bson"
)

type (
    SearchResponse struct {
        QueryRu  string `json:"queryRu"`
        ResultZh string `json:"resultZh"`
        Source   string `json:"source"`
    }

     SearchTranslation struct {
         Id       bson.ObjectId `bson:"_id"`
         Rus      string  `json:"rus"       bson:"rus"`
         RusNorm  string  `json:"rusNorm"   bson:"rusNorm"`
         TrBing   string  `json:"trBing"    bson:"trBing"`
         TrGoogle string  `json:"trGoogle"  bson:"trGoogle"`
         TrManual string  `json:"trManual"  bson:"trManual"`
     }
)

func SearchInsert( translate SearchResponse) error {
    session := s.GetSession()
    defer session.Close()
    coll := session.DB(s.DB).C("searchTranslations")

    search := &SearchTranslation{}
    search.Id  = bson.NewObjectId()
    search.Rus = translate.QueryRu
    search.RusNorm = ""
    if translate.Source == "bing" {
        search.TrBing  = translate.ResultZh
    }
    if translate.Source == "google" {
        search.TrGoogle  = translate.ResultZh
    }
    search.TrManual = ""

    return coll.Insert(search)
}
