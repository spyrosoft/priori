package main

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func authorize(handle httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		loggedIn, err := isLoggedIn(w, r)
		if err != nil {
			serve500(w)
		} else if loggedIn {
			handle(w, r, ps)
		} else {
			serveLoginPage(w)
		}
	}
}

func apiAuthorize(handle httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		loggedIn, err := isLoggedIn(w, r)
		if err != nil {
			json.NewEncoder(w).Encode(apiResponse{
				Errors: []string{err.Error()},
			})
		} else if loggedIn {
			handle(w, r, ps)
		} else {
			json.NewEncoder(w).Encode(apiResponse{
				Errors: []string{"You must log in first."},
			})
		}
	}
}
