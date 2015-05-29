package controllers

import (
	"encoding/json"
	"net/http"
	"../system"

	"github.com/zenazn/goji/web"
)

type Stat struct {
	State           string `json:"state"`
	Version         string `json:"version"`
	Hostname        string `json:"hostname"`
	BootTimestamp   string `json:"bootTimestamp"`
}

func Ping(c web.C, w http.ResponseWriter, r *http.Request) {
	stats := &Stat{
		State: "OK",
		Version: system.Version,
		Hostname: system.Hostname,
		BootTimestamp: system.Boot_time,
	}

	encoder := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	encoder.Encode(stats)
}
