package ip

import (
	"fmt"
	"net"
	"strings"
)

func getIPPrefix() string {
	ip := getIP()
	arr := strings.Split(ip, ".")
	return strings.Join(arr[:3], ".") + "."
}

func getIP() string {
	ifaces, err := net.Interfaces()
	if err != nil {
		fmt.Println(err)
		return ""
	}
	strs := make([]string, 0)

	for _, iface := range ifaces {
		addrs, err := iface.Addrs()
		if err != nil {
			fmt.Println(err)
			continue
		}
		for _, addr := range addrs {
			switch v := addr.(type) {
			case *net.IPNet:
				ip := v.IP
				strs = append(strs, ip.String())
			}
		}
	}
	for _, v := range strs {
		if strings.HasPrefix(v, "192.168.") {
			return v
		}
	}
	for _, v := range strs {
		if strings.HasPrefix(v, "10.") {
			return v
		}
	}
	for _, v := range strs {
		if strings.HasPrefix(v, "172.") {
			return v
		}
	}
	for _, v := range strs {
		if v != "127.0.0.1" && v != "::1" {
			return v
		}
	}

	if len(strs) == 0 {
		return ""
	}

	return strs[0]
}
