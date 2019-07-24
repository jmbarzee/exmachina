package lightorchastrator

import (
	"errors"
	"fmt"
	"log"
	"os/exec"
	"strconv"
)

func Start(port int, domainPort int, log *log.Logger) error {
	log.Printf("Starting neopixelbar!")

	cmdString := "go run $GOPATH/src/github.com/jmbarzee/domain/services/lightorchastrator/cmd/run/main.go"

	cmd := exec.Command(cmdString, strconv.Itoa(port), strconv.Itoa(domainPort))
	// pgid is same as parents by default

	err := cmd.Start()
	if err != nil {
		fmt.Printf(err.Error())
		return err
	}

	return nil
}

func Build() error {

	return errors.New("Unimplemented")
}
