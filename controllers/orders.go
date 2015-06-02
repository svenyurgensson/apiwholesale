package controllers

import (
    "net/http"

    "../models"

    "github.com/zenazn/goji/web"
)

func OrdersList(c web.C, w http.ResponseWriter, r *http.Request) {
    customer := c.Env["auth_customer"].(models.Customer)

}

func OrderCreate(c web.C, w http.ResponseWriter, r *http.Request) {
    customer := c.Env["auth_customer"].(models.Customer)
}

func OrderDelete(c web.C, w http.ResponseWriter, r *http.Request) {
    customer := c.Env["auth_customer"].(models.Customer)
}

func OrderUpdate(c web.C, w http.ResponseWriter, r *http.Request) {
    customer := c.Env["auth_customer"].(models.Customer)
}
