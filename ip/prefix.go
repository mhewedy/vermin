package ip

import (
	"fmt"
	"net"
	"strings"
)

func getIPPrefixes() []string {

	ipPrefixes := make([]string, 0)

	for _, ip := range getAllLocalIPAddresses() {
		arr := strings.Split(ip, ".")
		ipPrefixes = append(ipPrefixes, strings.Join(arr[:3], ".")+".")
	}
	return ipPrefixes
}

func getAllLocalIPAddresses() []string {
	ifaces, err := net.Interfaces()
	if err != nil {
		fmt.Println(err)
		return []string{}
	}
	ips := make([]string, 0)

	for _, iface := range ifaces {
		if strings.Contains(iface.Name, "VirtualBox") {
			continue
		}

		addrs, err := iface.Addrs()
		if err != nil {
			fmt.Println(err)
			continue
		}
		for _, addr := range addrs {
			switch v := addr.(type) {
			case *net.IPNet:
				if v.IP.To4() != nil {
					ips = append(ips, v.IP.String())
				}
			}
		}
	}
	return ips
}
