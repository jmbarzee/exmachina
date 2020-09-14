package system

import (
	"net"
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
