package main

import (
	"context"
	"io/ioutil"
	"os"
	"time"

	"github.com/jmbarzee/domain/server"
)

func main() {
	configFileName := os.Getenv("DOMAIN_CONFIG_FILE")
	configFile, err := os.OpenFile(configFileName, os.O_RDONLY, 0777)
	if err != nil {
		panic(err)
	}

	tomlBytes, err := ioutil.ReadAll(configFile)
	if err != nil {
		panic(err)
	}

	config, err := server.ConfigFromTOML(tomlBytes)
	if err != nil {
		panic(err)
	}

	server, err := server.NewDomain(config)
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()
	server.Run(ctx)
}
