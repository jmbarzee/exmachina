package main

import (
	"context"
	"runtime"
	"time"

	"github.com/jmbarzee/dominion/service/config"
	example "github.com/jmbarzee/dominion/services/example/service"
)

func main() {
	runtime.GOMAXPROCS(4)

	config, err := config.FromEnv("exampleService")
	if err != nil {
		panic(err)
	}

	example, err := example.NewExampleService(config)
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 240*time.Second)
	defer cancel()

	if err := example.Run(ctx); err != nil {
		panic(err)
	}
}
