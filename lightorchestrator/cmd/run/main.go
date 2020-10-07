package main

import (
	"context"
	"runtime"
	"time"

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

	ctx, cancel := context.WithTimeout(context.Background(), 240*time.Second)
	defer cancel()

	if err := lightOrch.Run(ctx); err != nil {
		panic(err)
	}
}
