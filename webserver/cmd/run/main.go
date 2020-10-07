package main

import (
	"context"
	"runtime"
	"time"

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

	ctx, cancel := context.WithTimeout(context.Background(), 240*time.Second)
	defer cancel()

	if err := example.Run(ctx); err != nil {
		panic(err)
	}
}
