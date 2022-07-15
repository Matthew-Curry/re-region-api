package bootstrap

/* Public method to run the app, register all routes with the mux */

import (
	"net/http"
	"fmt"

	"github.com/Matthew-Curry/re-region-api/logging"
	"github.com/Matthew-Curry/re-region-api/controller"

)

var logger, _ = logging.GetLogger("file.log")
var port string = ":8080"

func RunApp() {
	// initialize services in the controller package
	err := controller.InitServices()
	if err != nil {
		logger.Fatal("Unable to intiialize core services", err.Error())
	}

	// the multiplexer to handle requests
	mux := http.NewServeMux()

	// geo get endpoints
	mux.HandleFunc("/counties", controller.CountyHandler)
	mux.HandleFunc("/states", controller.StateHandler)

	// list endpoints
	mux.HandleFunc("/county-list", controller.CountyListHandler)
	mux.HandleFunc("/state-list", controller.StateListHandler)

	// general tax info endpoints
	mux.HandleFunc("/county-taxes", controller.CountyTaxesHandler)
	mux.HandleFunc("/state-taxes", controller.StateTaxesHandler)
	mux.HandleFunc("/federal-taxes", controller.FederalTaxesHandler)

	// health endpoint
	mux.HandleFunc("/health", controller.HealthHandler)

	logger.Info(fmt.Sprintf("Listening at %s", port))
	http.ListenAndServe(port, mux)
}