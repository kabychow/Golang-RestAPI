package controllers

import (
	"encoding/json"
	"net/http"
	"two-server/models"
	u "two-server/utils"
)

var CreateAccount = func(w http.ResponseWriter, r *http.Request) {
	var account *models.Account
	if json.NewDecoder(r.Body).Decode(&account) == nil {
		statusCode, response := account.Create()
		u.Respond(w, statusCode, response)
		return
	}
	u.Respond(w, 400, nil)
}

var Authenticate = func(w http.ResponseWriter, r *http.Request) {
	var account *models.Account
	if json.NewDecoder(r.Body).Decode(&account) == nil {
		statusCode, response := models.Login(account.Email, account.Password)
		u.Respond(w, statusCode, response)
		return
	}
	u.Respond(w, 400, nil)
}