package main

import (
	"context"
	"os"
	"runtime"
	"time"

	"github.com/jmbarzee/dominion/domain"
	"github.com/jmbarzee/dominion/domain/config"
)

func main() {
	runtime.GOMAXPROCS(4)

	configFileName := os.Getenv("DOMAIN_CONFIG_FILE")

	// Check config
	if err := config.SetupFromTOML(configFileName); err != nil {
		panic(err)
	}

	domain, err := domain.NewDomain(config.GetDomainConfig())
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 240*time.Second)
	defer cancel()

	if err := domain.Run(ctx); err != nil {
		panic(err)
	}
}
