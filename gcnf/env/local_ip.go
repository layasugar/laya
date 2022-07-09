package env

import (
	"log"
	"net"
)

var (
	envLocalIP = "127.0.0.1"
)

func LocalIP() string {
	if envLocalIP == "" {
		ips := getLocalIPs()
		if len(ips) > 0 {
			envLocalIP = ips[0] + ":80"
		} else {
			envLocalIP = defaultHttpListen
		}
	}

	return envLocalIP
}

func getLocalIPs() (ips []string) {
	interfaceAddr, err := net.InterfaceAddrs()
	if err != nil {
		log.Printf("fail to get net interface addrs: %v", err)
		return ips
	}

	for _, address := range interfaceAddr {
		ipNet, isValidIpNet := address.(*net.IPNet)
		if isValidIpNet && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				ips = append(ips, ipNet.IP.String())
			}
		}
	}
	return ips
}
