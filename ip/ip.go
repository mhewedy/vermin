package ip

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/mhewedy/vermin/hypervisor"
	"github.com/mhewedy/vermin/log"
)

type addr struct {
	ip  string
	mac string
}

func GetIpAddress(vmName string) (string, error) {

	health, err := hypervisor.HealthCheck(vmName)
	if err != nil {
		return "", err
	}

	// Wait for health status to be "Up" for up to 5 minutes
	timeout := time.After(5 * time.Minute)
	for *health != "Up" {
		select {
		case <-timeout:
			return "", fmt.Errorf("timeout waiting for VM health status to be 'Up'")
		default:
			time.Sleep(10 * time.Second) // Check every 10 seconds
			health, err = hypervisor.HealthCheck(vmName)
			if err != nil {
				return "", err
			}
		}
	}

	ipAddr, err := hypervisor.GetVMProperty(vmName, "ip")
	if err != nil {
		return "", err
	}
	return *ipAddr, nil
}

// Find will try to find IP for the VM.
//
// If the purge flag if true, it will invalidate the cache first then start the search process.
// Otherwise the search will start without clearing the cache, but if no result found, the cache will be cleared and the search
// will executed again.
func Find(vmName string, purge bool) (string, error) {

	mac, err := getMACAddr(vmName)
	if err != nil {
		return "", err
	}

	log.Debug("found mac: %s for vm: %s", mac, vmName)
	var pong bool

	if purge {
		log.Debug("purge=1, purging...")
		if err := ping(); err != nil {
			return "", err
		}
		pong = true
	}

	for {
		arp, err := getArpTable()

		if log.IsDebugEnabled() {
			log.Debug("here's the arp table:")

			sort.Slice(arp, func(i, j int) bool {
				return arp[i].mac > arp[j].mac
			})

			for _, e := range arp {
				log.Debug("IP: %s, MAC: %s", e.ip, e.mac)
			}
		}

		if err != nil {
			return "", err
		}

		for i := len(arp) - 1; i >= 0; i-- {
			a := arp[i]
			if a.mac == mac {
				return a.ip, nil
			}
		}

		if pong {
			break
		}

		if err := ping(); err != nil {
			return "", err
		}
		pong = true
	}

	return "", fmt.Errorf("cannot find ip for %s\nuse the command 'vermin ip -p %s' to purge cache", vmName, vmName)
}

func getMACAddr(vmName string) (string, error) {
	box, err := hypervisor.GetBoxInfo(vmName)
	if err != nil {
		return "", err
	}

	return formatMACAddr(strings.ToLower(box.MACAddr)), nil
}

func formatMACAddr(mac string) string {
	ret := make([]rune, 0)

	for i := range mac {
		if i%2 == 0 && mac[i] == '0' {
			continue
		}
		ret = append(ret, rune(mac[i]))
		if i%2 == 1 && i < len(mac)-1 {
			ret = append(ret, ':')
		}
	}
	return string(ret)
}
