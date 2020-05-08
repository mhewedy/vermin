package ip

import (
	"fmt"
	"github.com/mhewedy/vermin/command"
	"github.com/mhewedy/vermin/db/info"
	"strconv"
	"strings"
	"sync"
)

type addr struct {
	ip  string
	mac string
}

const max = 255

func Find(vmName string, purge bool) (string, error) {

	mac, err := getMACAddr(vmName)
	if err != nil {
		return "", err
	}

	var pong bool

	if purge {
		ping()
		pong = true
	}

	for {
		arp := getArpTable()
		for i := len(arp) - 1; i >= 0; i-- {
			a := arp[i]
			if a.mac == mac {
				return a.ip, nil
			}
		}

		if pong {
			break
		}

		ping()
		pong = true
	}

	return "", fmt.Errorf("Cannot find ip for %s\nUse the command 'vermin ip -p %s' to purge cache", vmName, vmName)
}

func ping() {
	var wg sync.WaitGroup
	wg.Add(max)

	for i := range [max]int{} {
		go func(i int) {
			ip := getIPPrefix() + strconv.Itoa(i)
			_ = command.Ping(ip).Run()
			wg.Done()
		}(i)
	}

	wg.Wait()
}

func getMACAddr(vmName string) (string, error) {
	box, err := info.GetBoxInfo(vmName)
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
