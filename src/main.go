package main

import (
	"embed"
	"fmt"
	"io/fs"
	"net/http"
	"os"

	"github.com/Matthew-Curry/re-region-api/src/controller"
	"github.com/Matthew-Curry/re-region-api/src/logging"
)

const openApiYml = "docs.yml"

var logger logging.Logger
var logFile *os.File

//go:embed static/swagger-ui
var content embed.FS

func main() {
	// the logger and file to close
	var logger, logfile = logging.GetLogger("file.log")
	defer logfile.Close()
	// read in env vars
	port := ":" + os.Getenv("PORT")
	dbUser := os.Getenv("RE_REGION_API_USER")
	dbPassword := os.Getenv("RE_REGION_API_PASSWORD")
	dbName := os.Getenv("RE_REGION_DB")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")

	// initialize services in the controller package
	err := controller.InitServices(dbUser, dbPassword, dbName, dbHost, dbPort)
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

	// swagger ui
	// serve the open API yml for the swagger ui to point at
	mux.HandleFunc("/docs/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		w.Header().Set("Access-Control-Allow-Methods", "*")
		http.ServeFile(w, r, openApiYml)
	})


	// the path for the UI
	fsys, _ := fs.Sub(content, "static/swagger-ui")
	mux.Handle("/", http.FileServer(http.FS(fsys)))

	logger.Info(fmt.Sprintf("Listening at %s", port))
	http.ListenAndServe(port, mux)
}
