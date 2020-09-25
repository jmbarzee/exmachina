package system

import (
	"context"
	"net"
	"time"
)

// func GetOutboundIP() (net.IP, error) {
// 	ifaces, err := net.Interfaces()
// 	if err != nil {
// 		return nil, err
// 	}
// 	for _, i := range ifaces {
// 		Logf("interface [%s]: %v %v", i.Name, i.Index, i.HardwareAddr)
// 		addrs, err := i.Addrs()
// 		if err != nil {
// 			return nil, err
// 		}
// 		for _, addr := range addrs {
// 			var ip net.IP
// 			switch v := addr.(type) {
// 			case *net.IPNet:
// 				ip = v.IP
// 				v.IP
// 				Logf("IPNet: %s", ip.String())
// 			case *net.IPAddr:
// 				ip = v.IP
// 				Logf("IPAddr: %s", ip.String())
// 			}
// 			// process IP address
// 		}
// 	}
// 	return nil, fmt.Errorf("FAKE ERROR")
// }

func GetOutboundIP() (net.IP, error) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return net.IP{}, err
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP, nil
}

// RoutineCheck offers repeatedly runs check and then waits for wait.
// cancleing ctx will end the check
func RoutineCheck(ctx context.Context, routineName string, wait time.Duration, check func(context.Context)) {
	LogRoutinef(routineName, "Starting routine")
	ticker := time.NewTicker(wait)

Loop:
	for {
		select {
		case <-ticker.C:
			check(ctx)
		case <-ctx.Done():
			break Loop
		}
	}
	LogRoutinef(routineName, "Stopping routine")
}
