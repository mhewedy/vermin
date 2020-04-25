package ip

import (
	"errors"
	"fmt"
	"github.com/mhewedy/vermin/cmd"
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

	return "", errors.New(fmt.Sprintf("cannot find ip for %s", vmName))
}

func ping() {
	var wg sync.WaitGroup
	wg.Add(max)

	for i := range [max]int{} {
		go func(i int) {
			ip := getIPPrefix() + strconv.Itoa(i)
			_ = cmd.Run("ping", "-c", "1", "-W", "0.1", ip)
			wg.Done()
		}(i)
	}

	wg.Wait()
}

func getMACAddr(vmName string) (string, error) {

	out, _ := cmd.Execute("vboxmanage", "showvminfo", vmName, "--machinereadable")
	entries := strings.Fields(out)

	for _, entry := range entries {
		if strings.HasPrefix(entry, "macaddress1") {
			mac := strings.Split(entry, "=")[1]
			mac = strings.ToLower(strings.Trim(mac, `""`))
			return formatMACAddr(mac), nil
		}
	}
	return "", errors.New(fmt.Sprintf("unable to get mac address for %s", vmName))
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

func getArpTable() []addr {

	addrs := make([]addr, 0)

	out, _ := cmd.Execute("arp", "-an")
	entries := strings.Split(out, "\n")

	for _, entry := range entries {
		fields := strings.Fields(entry)
		if len(fields) > 4 {
			ip := strings.TrimFunc(fields[1], func(r rune) bool {
				return r == '(' || r == ')'
			})
			addrs = append(addrs, addr{
				ip:  ip,
				mac: fields[3],
			})
		}
	}

	return addrs
}
