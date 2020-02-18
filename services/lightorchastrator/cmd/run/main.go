package main

import (
	"os"
	"strconv"

	"github.com/jmbarzee/domain/services/lightorchastrator"
)

func main() {

	portString := os.Args[1]
	port, err := strconv.ParseInt(portString, 0, 32)
	if err != nil {
		panic(err)
	}

	domainPortString := os.Args[2]
	domainPort, err := strconv.ParseInt(domainPortString, 0, 32)
	if err != nil {
		panic(err)
	}

	lo := lightorchastrator.NewLightOrch(port, domainPort)
	lo.Run()
}
