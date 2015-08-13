package main

import (
	"database/sql"

	"github.com/gorilla/mux"
)

func createRouter(db *sql.DB, public_key, private_key []byte) *mux.Router {
	router := mux.NewRouter().StrictSlash(true)

	tokenRouter := router.PathPrefix("/token/").Subrouter()
	tokenRouter.Path("/").Methods("POST").HandlerFunc(tokenCreationHandler(db, private_key))
	tokenRouter.Path("/verify").Methods("POST").HandlerFunc(tokenVerificationHandler(public_key))

	userRouter := router.PathPrefix("/credentials/").Subrouter()
	userRouter.Methods("POST").HandlerFunc(createCredentials)
	userRouter.Methods("PATCH").HandlerFunc(updateCredentials)
	userRouter.Methods("DELETE").HandlerFunc(deleteCredentials)

	return router
}
