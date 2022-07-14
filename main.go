package main

import ("github.com/Matthew-Curry/re-region-api/logging"
	"github.com/Matthew-Curry/re-region-api/bootstrap"

)

var logger, _ = logging.GetLogger("file.log")

func main() {
	logger.Info("App started")
	bootstrap.RunApp()
}