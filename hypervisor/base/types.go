package base

import (
	"encoding/binary"
	"fmt"
	"math"
	"net"
)

type Hypervisor interface {
	Start(vmName string) error

	Commit(vmName, imageName string) error

	Create(imageName, vmName string, cpus int, mem int) error

	List(all bool, excludeFunc func(vmName string) bool) ([]string, error)

	Stop(vmName string) error

	Remove(vmName string) error

	Modify(vmName string, cpus int, mem int) error

	ShowGUI(vmName string) error

	AddMount(vmName, ipAddr, hostPath, guestPath string) error

	RemoveMounts(vmName, ipAddr string) error

	ListMounts(vmName, ipAddr string) ([]MountPath, error)

	SetNetworkAdapterAsBridge(vmName string) error

	GetBoxInfo(vmName string) (*Box, error)

	GetSubnet() (*Subnet, error)
}

type MountPath struct {
	HostPath  string
	GuestPath string
}

type Box struct {
	CPU      string
	Mem      string
	DiskSize string
	MACAddr  string
}

type Subnet struct {
	start   uint32 // start ip address in the subnet in int32 representation
	current uint32 // holds the int32 of ip, used to increment by 1 in next method
	Len     int    // length of ips in the subnet
}

// anIp is some ip inside the subnet
func NewSubnet(anIp, netmask string) (*Subnet, error) {

	ones, bits := net.IPMask(net.ParseIP(netmask).To4()).Size()
	_, ipNet, err := net.ParseCIDR(fmt.Sprintf("%s/%d", anIp, ones))
	if err != nil {
		return nil, err
	}

	ipAddr := ip2int(ipNet.IP)
	l := int(math.Pow(2, float64(bits-ones)) - 1)
	return &Subnet{
		start:   ipAddr,
		current: ipAddr,
		Len:     l,
	}, nil
}

func (c *Subnet) HasNext() bool {
	max := c.start + uint32(c.Len)
	if c.current < max {
		return true
	}
	return false
}

func (c *Subnet) Next() *Subnet {
	return &Subnet{
		start:   c.start,
		current: c.current + 1,
		Len:     c.Len,
	}
}

func (c *Subnet) IP() string {
	return int2ip(c.current)
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
