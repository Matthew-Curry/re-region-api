package bootstrap

/* Public method to run the app, register all routes with the mux */

import (
	"fmt"
	"net/http"
	"io/fs"
	"strings"
	"embed"

	"github.com/Matthew-Curry/re-region-api/controller"
	"github.com/Matthew-Curry/re-region-api/logging"
)

const openApiYml = "docs.yml"

var logger, _ = logging.GetLogger("file.log")
var port string = ":8080"

//go:embed static
var content embed.FS

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

	// swagger ui
	// serve the open API yml for the swagger ui to point at
	mux.HandleFunc("/static/", func( res http.ResponseWriter, req *http.Request) {
		// restrict openapi doc to local, just for use by swagger ui
		a := strings.Split(req.RemoteAddr, ":")[0]
		if a == "127.0.0.1" {
			http.ServeFile(res, req, "docs.yml")
		}
	})

	// the path for the UI
	fsys, _ := fs.Sub(content, "static")
	mux.Handle("/swagger-ui/", http.FileServer(http.FS(fsys)))

	logger.Info(fmt.Sprintf("Listening at %s", port))
	http.ListenAndServe(port, mux)
}