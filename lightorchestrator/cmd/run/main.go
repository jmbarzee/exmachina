package main

import (
	"context"
	"runtime"

	"github.com/jmbarzee/dominion/service/config"
	lightorch "github.com/jmbarzee/services/lightorchestrator/service"
)

func main() {
	runtime.GOMAXPROCS(4)

	config, err := config.FromEnv("lightOrchestrator")
	if err != nil {
		panic(err)
	}

	lightOrch, err := lightorch.NewLightOrch(config)
	if err != nil {
		panic(err)
	}

	if err := lightOrch.Run(context.Background()); err != nil {
		panic(err)
	}
}
