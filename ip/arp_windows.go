package ip

import (
	"strings"
)

func getArpTable() ([]addr, error) {

	addrs := make([]addr, 0)

	s, err := cmd.Arp("-a").Call()
	if err != nil {
		return nil, err
	}
	ss := strings.Split(s, "\n\r")

	for _, out := range ss {
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
	}
	return addrs, nil
}
