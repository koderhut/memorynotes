package main

import (
	"github.com/gorilla/mux"
	"github.com/koderhut/memorynotes/config"
	"github.com/koderhut/memorynotes/controllers"
	"github.com/koderhut/memorynotes/urlgen"
)

func ConfigRouter(c *config.Context) *mux.Router {
	// init router
	router := mux.NewRouter()
	generator := urlgen.FromConfig(*c)
	ctrl := controllers.NewNotesHandler(generator)

	api := router.PathPrefix(c.PathPrefix).Subrouter()

	// route handlers/endpoints
	api.HandleFunc("/notes/{note}", ctrl.Retrieve).Methods("GET")
	api.HandleFunc("/notes", ctrl.Store).Methods("POST")

	return router
}
