package main

import (
	"context"
	"runtime"
	"time"

	"github.com/jmbarzee/dominion/dominion"
	"github.com/jmbarzee/dominion/dominion/config"
	"github.com/jmbarzee/dominion/system"
)

func main() {
	runtime.GOMAXPROCS(4)

	configFileName := system.RequireEnv("DOMINION_CONFIG_FILE")

	// Check config
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
