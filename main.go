package main

import (
	"NeonManager/data"
	"NeonManager/logger"
	"NeonManager/web"
)

func main() {

	logger.FatalOrLog("Failed to init DB: %v", data.Init(), "Initialized database")
	logger.FatalIfErr("Failed to serve http: %v", web.Serve())
}
