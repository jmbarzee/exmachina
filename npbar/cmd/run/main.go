package main

import (
	"context"
	"runtime"
	"time"

	"github.com/jmbarzee/dominion/service/config"
	"github.com/jmbarzee/services/npbar/service"
)

func main() {
	runtime.GOMAXPROCS(4)

	config, err := config.FromEnv("npBar")
	if err != nil {
		panic(err)
	}

	light, err := service.NewNPBar(config)
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()

	if err := light.Run(ctx); err != nil {
		panic(err)
	}
}
