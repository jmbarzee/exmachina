package main

import (
	"github.com/jmbarzee/domain/services"
	"github.com/jmbarzee/domain/services/webserver/service"
)

func main() {
	port, domainPort, logger, err := services.GatherStandardArgs()
	if err != nil {
		panic(err)
	}

	server, err := service.NewWebServer(port, domainPort, logger, "/usr/local/domain/services/webserver/service/static")
	if err != nil {
		panic(err)
	}

	server.Run()
}
