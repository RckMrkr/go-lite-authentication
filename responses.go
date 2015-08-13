package main

import (
	"encoding/json"
	"net/http"
)

type response struct {
	Data  interface{}
	Error err
}
type err struct {
	Code     int
	HttpCode int
	Message  string
}

func Respond(w http.ResponseWriter, s interface{}, code int) {
	if s != nil {
		json.NewEncoder(w).Encode(s)
	}
	w.WriteHeader(code)
}

func RespondSuccess(w http.ResponseWriter, s interface{}) {
	Respond(w, s, 200)
}

func RespondBadRequest(w http.ResponseWriter, message string) {
	s := struct {
		Error string `json:"error"`
	}{
		Error: message,
	}
	Respond(w, s, http.StatusBadRequest)
}
