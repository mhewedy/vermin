// +build !windows

package ip

import (
	"github.com/mhewedy/vermin/cmd"
	"strings"
)

func doPing(ip string) error {
	return cmd.Run("ping", "-c", "1", "-W", "0.1", ip)
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
