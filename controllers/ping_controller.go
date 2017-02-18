package controllers

import (
    "encoding/json"
    "net/http"
    s "apiwholesale/system"

    "github.com/zenazn/goji/web"
    "gopkg.in/mgo.v2"
)

type Stat struct {
    State           string `json:"state"`
    Version         string `json:"version"`
    Env             string `json:"environment"`
    Hostname        string `json:"hostname"`
    RequestsTotal   int    `json:"requestsTotal"`
    RequestsFailed  int    `json:"requestsFailed"`
    BootTimestamp   string `json:"bootTimestamp"`
    DBName          string `json:"dbName"`
    DBSocketsAlive  int `json:"dbSocketsAlive"`
    DBSocketsInUse  int `json:"dbSocketsInUse"`
    DBVersion       string `json:"mongoVersion"`
    SysInfo         string `json:"sysInfo"`
}


func Ping(c web.C, w http.ResponseWriter, r *http.Request) {
    session := s.GetSession()
    defer session.Close()

    state := "OK"
    if session.Ping() != nil {
        state = "FAIL"
    }

    mstats := mgo.GetStats()
    binfo, err := session.BuildInfo()
    if err != nil {
        binfo = mgo.BuildInfo{}
    }

    stats := &Stat{
        State: state,
        Version: s.Version,
        Env: s.Env,
        Hostname: s.Hostname,
        RequestsTotal: s.RequestsTotal,
        RequestsFailed: s.RequestsFailed,
        BootTimestamp: s.Boot_time,
        DBName: s.DB,
        DBSocketsAlive: mstats.SocketsAlive,
        DBSocketsInUse: mstats.SocketsInUse,
        DBVersion: binfo.Version,
        SysInfo: binfo.SysInfo,
    }

    encoder := json.NewEncoder(w)
    w.Header().Set("Content-Type", "application/json")

    encoder.Encode(stats)
}

func Favicon(c web.C, w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "image/x-icon")
    w.WriteHeader(http.StatusOK)
}
