package nplight

import (
	"errors"
	"fmt"
	"log"
	"os/exec"
	"strconv"
)

func Start(port int, domainPort int) error {
	log.Printf("Starting nplight!")

	path := "/usr/local/domain/services/nplight/cmd/run/main.go"
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
