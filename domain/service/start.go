package service

import (
	"fmt"
	"net"
	"os/exec"
	"path"
	"strconv"
	"strings"

	"github.com/jmbarzee/dominion/system"
)

func Start(serviceType string, ip net.IP, dominionPort int, domainUUID string, servicePort int) error {
	system.Logf("Starting %v!", serviceType)

	rootPath := "/usr/local/dominion/services"
	makefilePath := path.Join(rootPath, strings.ToLower(serviceType))
	makePath, err := exec.LookPath("make")
	if err != nil {
		return fmt.Errorf("make was not found in path: %w", err)
	}

	cmd := exec.Command(
		makePath,
		"-C", makefilePath,
		"DOMINION_IP="+ip.String(),
		"DOMINION_PORT="+strconv.Itoa(dominionPort),
		"DOMAIN_UUID="+domainUUID,
		"SERVICE_PORT="+strconv.Itoa(servicePort))
	// pgid is same as parents by default

	err = cmd.Start()
	if err != nil {
		fmt.Printf(err.Error())
		return err
	}

	// fmt.Println(cmd.Dir)
	// fmt.Println(cmd.Path)
	// fmt.Println(cmd.Args)

	// bytes, err := cmd.Output()
	// if err != nil {
	// 	fmt.Println("Error: ", err)
	// 	return err
	// }
	// fmt.Printf("Output: %s", bytes)

	return nil
}
