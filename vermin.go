package main

import (
	"strings"
	"vermin/cmd"
	"vermin/info"
)

func ps(all bool) (string, error) {

	var args = [2]string{"list"}
	if all {
		args[1] = "vms"
	} else {
		args[1] = "runningvms"
	}

	r, err := cmd.Execute("vboxmanage", args[:]...)
	if err != nil {
		return "", err
	}

	var vms []string
	fields := strings.Fields(r)

	for i := range fields {
		if i%2 == 0 {
			vms = append(vms, strings.ReplaceAll(fields[i], `"`, ""))
		}
	}

	return info.List(vms), nil
}
