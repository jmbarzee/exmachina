package main

import (
	"context"
	"io/ioutil"
	"os"
	"strconv"
	"time"

	"github.com/jmbarzee/domain"
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

	config, err := domain.ConfigFromTOML(tomlBytes)
	if err != nil {
		panic(err)
	}

	for i := 0; i < 3; i++ {
		config.UUID = "GUID-" + strconv.Itoa(i)
		config.Port += 100
		_, err = domain.NewDomain(context.TODO(), config)
		if err != nil {
			panic(err)
		}
	}

	c := time.After(time.Second * 20)
	<-c
}
