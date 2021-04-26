package main

import (
	"context"
	"runtime"

	"github.com/jmbarzee/dominion/service/config"
	"github.com/jmbarzee/services/webserver/service"
)

func main() {
	runtime.GOMAXPROCS(4)

	config, err := config.FromEnv("webServer")
	if err != nil {
		panic(err)
	}

	example, err := service.NewWebServer(config)
	if err != nil {
		panic(err)
	}

	if err := example.Run(context.Background()); err != nil {
		panic(err)
	}
}
