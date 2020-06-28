package webserver

import (
	"errors"
	"fmt"
	"log"
	"os/exec"
	"strconv"
)

func Start(port int, domainPort int) error {
	log.Printf("Starting WebServer!")

	path := "/usr/local/domain/services/webserver/cmd/run/main.go"
	gopath := "/usr/local/go/bin/go"

	cmd := exec.Command(gopath, "run", path, strconv.Itoa(port), strconv.Itoa(domainPort), "/usr/local/domain/logs/webserver.log")
	// pgid is same as parents by default

	err := cmd.Start()
	if err != nil {
		fmt.Printf(err.Error())
		return err
	}
	log.Printf("Started WebServer!")

	return nil
}

func Build() error {

	return errors.New("Unimplemented")
}
