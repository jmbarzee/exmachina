package main

import (
	"context"
	"runtime"
	"time"

	"github.com/faiface/pixel/pixelgl"
	"github.com/jmbarzee/dominion/service/config"
	"github.com/jmbarzee/services/lightorchestrator/clients/nptest"
)

func main() {
	runtime.GOMAXPROCS(4)

	pixelgl.Run(run)
}

func run() {
	config, err := config.FromEnv("npBar")
	if err != nil {
		panic(err)
	}

	light, err := NewNPBar(config)
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()

	if err := light.Run(ctx); err != nil {
		panic(err)
	}

}

type (
	NPBar struct {
		*nptest.NPTest
	}
)

const (
	size = 30 //TODO @jmbarzee fetch from config
)

func NewNPBar(config config.ServiceConfig) (NPBar, error) {
	sub, err := nptest.NewNPTest(config, size)
	if err != nil {
		return NPBar{}, err
	}
	return NPBar{sub}, nil
}
