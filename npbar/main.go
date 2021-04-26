package main

// +build !test

import (
	"context"
	"runtime"

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

	if err := light.Run(context.Background()); err != nil {
		panic(err)
	}
}
