// +build !windows

package ip

import (
	"github.com/mhewedy/vermin/command"
	"strings"
)

func getArpTable() ([]addr, error) {

	addrs := make([]addr, 0)

	out, err := command.Arp("-an").Call()
	if err != nil {
		return nil, err
	}
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

	return addrs, nil
}
