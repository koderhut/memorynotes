package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/koderhut/memorynotes/config"
	"github.com/koderhut/memorynotes/controllers"
	"github.com/koderhut/memorynotes/models"
	"github.com/koderhut/memorynotes/urlgen"
	"net/http"
	"time"
)

func ConfigRouter(c *config.Context) *http.Server {
	// init router
	router := mux.NewRouter()
	generator := urlgen.FromConfig(*c)
	ctrl := controllers.NewNotesHandler(generator)

	api := router.PathPrefix(c.PathPrefix).Subrouter()

	// route handlers/endpoints
	api.HandleFunc("/notes/{note}", ctrl.Retrieve).Methods("GET")
	api.HandleFunc("/notes", ctrl.Store).Methods("POST")

	// a simple stats endpoint
	router.HandleFunc("/stats", func(w http.ResponseWriter, r *http.Request) {
		current, total := models.Stats()
		w.Write([]byte(fmt.Sprintf("Notes: %d\n", current)))
		w.Write([]byte(fmt.Sprintf("Total: %d\n", total)))
	}).Host(fmt.Sprintf("localhost:%s", c.Port))


	srv := &http.Server{
		Addr: c.Addr(),
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler: router,
	}

	return srv
}
