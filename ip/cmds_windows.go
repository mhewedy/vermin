package ip

import (
	"github.com/mhewedy/vermin/cmd"
	"strings"
)

func doPing(ip string) error {
	return cmd.Run("ping", "-n", "1", "-w", "0.1", ip)
}

func getArpTable() []addr {

	addrs := make([]addr, 0)

	out, _ := cmd.Execute("arp", "-aN", getIP())
	entries := strings.Split(out, "\n")

	for _, entry := range entries[3:] {
		fields := strings.Fields(entry)
		if len(fields) > 1 {
			addrs = append(addrs, addr{
				ip:  fields[0],
				mac: fields[1],
			})
		}
	}

	return addrs
}
