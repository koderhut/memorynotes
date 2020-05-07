package webapp

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/koderhut/memorynotes/config"
	"github.com/koderhut/memorynotes/contracts"
	"github.com/koderhut/memorynotes/stats"
)

var noteStats = stats.New()

func BootstrapRouter(c *config.Parameters, routing ...WebRouting) *mux.Router {
	// init router
	router := mux.NewRouter()

	router.Use(noteStatsLogger)
	router.Use(mux.CORSMethodMiddleware(router))
	router.Use(corsAllowedHost(c.Web.CorsHost))

	router.
		Path("/stats").
		Methods(http.MethodGet).
		HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
			w.Header().Set("Content-Type", "application/json")

			json.NewEncoder(w).Encode(&contracts.StatsMessage{Status: true, StoredNotes: noteStats.Current, TotalNotes: noteStats.Total})
		}).
		Host(fmt.Sprintf("localhost:%s", c.Port))

	api := router.PathPrefix(c.Web.PathPrefix).Subrouter()

	for _, routerCfg := range routing {
		routerCfg.RegisterRoutes(api)
	}

	printRoutes(router)

	return router
}

func noteStatsLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		route := mux.CurrentRoute(r)

		switch route.GetName() {
		case "notes_store":
			noteStats.Inc()
		case "notes_fetch":
			noteStats.Decr()
		}

		next.ServeHTTP(w, r)
	})
}

func corsAllowedHost(cors string) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", cors)

			next.ServeHTTP(w, r)
		})
	}
}

func printRoutes(r *mux.Router) {
	log.Println(">>> Registered routes:")
	r.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		path, _ := route.GetPathTemplate()
		methods, _ := route.GetMethods()
		log.Printf("Route: %s, Methods: %s\n", path, strings.Join(methods, ","))
		return nil
	})
	log.Println(">>>")
}
