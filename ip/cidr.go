package ip

import (
	"encoding/binary"
	"fmt"
	"math"
	"net"
)

type cidr struct {
	base uint32 // base ip address in the cidr block in int32 representation
	ip   uint32 // holds the int32 of ip, used to increment by 1 in next method
	len  int    // length of ips in the subnet
}

func (c cidr) hasNext() bool {
	max := c.base + uint32(c.len)
	if c.ip < max {
		return true
	}
	return false
}

func (c cidr) next() cidr {
	return cidr{
		base: c.base,
		ip:   c.ip + 1,
		len:  c.len,
	}
}

func (c cidr) IP() string {
	return int2ip(c.ip)
}

func getCIDRs() []cidr {

	ifaces, err := net.Interfaces()
	if err != nil {
		fmt.Println(err)
		return []cidr{}
	}

	cidrs := make([]cidr, 0)

	for _, iface := range ifaces {

		addrs, err := iface.Addrs()
		if err != nil {
			fmt.Println(err)
			continue
		}
		for _, addr := range addrs {
			switch v := addr.(type) {
			case *net.IPNet:
				if v.IP.To4() != nil {
					ones, bits := v.IP.DefaultMask().Size()
					ip := v.IP.Mask(v.IP.DefaultMask())
					cidrs = append(cidrs, cidr{
						base: ip2int(ip),
						ip:   ip2int(ip),
						len:  int(math.Pow(2, float64(bits-ones)) - 1),
					})
				}
			}
		}
	}
	return cidrs
}

func ip2int(ip net.IP) uint32 {
	if len(ip) == 16 {
		return binary.BigEndian.Uint32(ip[12:16])
	}
	return binary.BigEndian.Uint32(ip)
}

func int2ip(nn uint32) string {
	ip := make(net.IP, 4)
	binary.BigEndian.PutUint32(ip, nn)
	return ip.String()
}
