package main

import (
	"context"
	"os"
	"time"

	"github.com/jmbarzee/dominion/domain"
)

func main() {
	configFileName := os.Getenv("DOMAIN_CONFIG_FILE")

	domain, err := domain.NewDomain(configFileName)
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()
	domain.Run(ctx)
}
