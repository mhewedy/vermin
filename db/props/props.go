// This package modeled around showvminfo subcommand which yield a bad performance.
// Consider using db/info package (which depends on parsing the vbox xml file) instead whenever possible
package props

import (
	"github.com/mhewedy/vermin/command"
	"strings"
)

func List(vmName string) ([]string, error) {
	out, err := command.VBoxManage("showvminfo", vmName, "--machinereadable").Call()
	if err != nil {
		return nil, err
	}
	return strings.Fields(out), nil
}

func FindByPrefix(vmName string, prefix string) ([]string, error) {

	entries, err := List(vmName)
	if err != nil {
		return nil, err
	}

	var values []string

	for i := range entries {
		entry := entries[i]
		if strings.HasPrefix(entry, prefix) {
			value := strings.Split(entry, "=")[1]
			value = strings.Trim(value, `""`)
			values = append(values, value)
		}
	}

	return values, nil
}

func FindFirstByPrefix(vmName string, prefix string) (string, bool, error) {
	byPrefix, err := FindByPrefix(vmName, prefix)
	if err != nil {
		return "", false, err
	}

	if len(byPrefix) == 0 {
		return "", false, nil
	}

	return byPrefix[0], true, nil
}
