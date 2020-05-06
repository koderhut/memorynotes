package webapp

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/koderhut/memorynotes/config"
	"github.com/koderhut/memorynotes/stats"
	"net/http"
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

			w.Write([]byte(fmt.Sprintf("Stored Notes: %d\n", noteStats.Current)))
			w.Write([]byte(fmt.Sprintf("Total Notes: %d\n", noteStats.Total)))
		}).
		Host(fmt.Sprintf("localhost:%s", c.Port))

	api := router.PathPrefix(c.Web.PathPrefix).Subrouter()

	for _, routerCfg := range routing {
		routerCfg.RegisterRoutes(api)
	}

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