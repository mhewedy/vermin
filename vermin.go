package main

import (
	"vermin/info"
)

func ps(all bool) (string, error) {
	vms, err := info.List(all)
	if err != nil {
		return "", err
	}
	return info.Get(vms), nil
}
