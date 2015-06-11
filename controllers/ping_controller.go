package controllers

import (
    "encoding/json"
    "net/http"
    "../system"

    "github.com/zenazn/goji/web"
//    "gopkg.in/mgo.v2"
)

type Stat struct {
    State           string `json:"state"`
    Version         string `json:"version"`
    Hostname        string `json:"hostname"`
    BootTimestamp   string `json:"bootTimestamp"`
    DBSocketsAlive  int `json:"dbSocketsAlive"`
    DBSocketsInUse  int `json:"dbSocketsInUse"`
}


func Ping(c web.C, w http.ResponseWriter, r *http.Request) {
    state := "OK"
    if system.GetSession().Ping() != nil {
        state = "FAIL"
    }
//    mstats := mgo.GetStats()

    stats := &Stat{
        State: state,
        Version: system.Version,
        Hostname: system.Hostname,
        BootTimestamp: system.Boot_time,
        DBSocketsAlive: 22, //mstats.SocketsAlive,
        DBSocketsInUse: 33, //mstats.SocketsInUse,
    }

    encoder := json.NewEncoder(w)
    w.Header().Set("Content-Type", "application/json")

    encoder.Encode(stats)
}

func Favicon(c web.C, w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "image/x-icon")
    w.WriteHeader(http.StatusOK)
}
