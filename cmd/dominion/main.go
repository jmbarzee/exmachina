package main

import (
	"context"
	"os"
	"time"

	"github.com/jmbarzee/dominion/dominion"
)

func main() {
	configFileName := os.Getenv("DOMINION_CONFIG_FILE")

	dominion, err := dominion.NewDominion(configFileName)
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()
	dominion.Run(ctx)
}
