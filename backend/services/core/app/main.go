package main

import (
	"github.com/scalent.io/scalent-hrms/pkg/log"
)

func main() {

	err := initConfig()
	if err != nil {
		log.Print(err)
		return
	}

	registry, err := initServer(&config)
	if err != nil {
		log.Print(err)
		return
	}

	err = registry.StartServer()
	if err != nil {
		log.Print(err)
	}
}
