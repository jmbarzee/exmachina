package main

import (
	"os"
	"strconv"

	"github.com/jmbarzee/domain/services/neopixelbar"
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

	bar := neopixelbar.NewNeoBar(int(port), int(domainPort))
	bar.Run()
}
