package main

import (
	"NeonManager/data"
	"NeonManager/web"
	"log"
)

func main() {

	if err := data.Init(); err != nil {
		log.Fatal(err)
	}

	if err := web.Serve(); err != nil {
		log.Fatal(err)
	}
}
