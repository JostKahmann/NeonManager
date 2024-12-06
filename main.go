package main

import (
	"NeonManager/data"
	"NeonManager/logger"
	"NeonManager/web"
	"log"
)

func main() {

	if err := data.Init(); err != nil {
		log.Fatal(err)
	}
	logger.Info("Initialized database")

	if err := web.Serve(); err != nil {
		log.Fatal(err)
	}
}
