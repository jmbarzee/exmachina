package main

import (
	"github.com/jmbarzee/domain/services"
	"github.com/jmbarzee/domain/services/nplight/service"
)

func main() {
	port, domainPort, logger, err := services.GatherStandardArgs()
	if err != nil {
		panic(err)
	}

	light := service.NewNPLight(port, domainPort, logger)
	light.Run()
}
