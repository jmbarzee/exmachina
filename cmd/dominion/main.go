package main

import (
	"context"
	"os"
	"runtime"
	"time"

	"github.com/jmbarzee/dominion/dominion"
	"github.com/jmbarzee/dominion/dominion/config"
)

func main() {
	runtime.GOMAXPROCS(4)

	// Check config
	configFileName := os.Getenv("DOMINION_CONFIG_FILE")
	if err := config.SetupFromTOML(configFileName); err != nil {
		panic(err)
	}

	dominion, err := dominion.NewDominion(config.GetDominionConfig())
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 240*time.Second)
	defer cancel()

	if err := dominion.Run(ctx); err != nil {
		panic(err)
	}
}
