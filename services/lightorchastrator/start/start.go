package lightorchastrator

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"strconv"
)

func Start(port int, domainPort int, log *log.Logger) error {
	log.Printf("Starting neopixelbar!")

	lightOrchPath := "src/github.com/jmbarzee/domain/services/lightorchastrator/cmd/run/main.go"
	path := path.Join(os.Getenv("GOPATH"), lightOrchPath)
	gopath := "/usr/local/go/bin/go"

	cmd := exec.Command(gopath, "run", path, strconv.Itoa(port), strconv.Itoa(domainPort))
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
