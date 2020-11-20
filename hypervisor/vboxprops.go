package hypervisor

import (
	"strings"
)

func findByPrefix(vmName string, prefix string) ([]string, error) {

	entries, err := Info(vmName)
	if err != nil {
		return nil, err
	}

	var values []string

	for i := range entries {
		entry := entries[i]
		if strings.HasPrefix(entry, prefix) || strings.HasPrefix(strings.Trim(entry, `"`), prefix) {
			value := strings.Split(entry, "=")[1]
			value = strings.Trim(value, `"`)
			values = append(values, value)
		}
	}

	return values, nil
}

func findFirstByPrefix(vmName string, prefix string) (string, bool, error) {
	byPrefix, err := findByPrefix(vmName, prefix)
	if err != nil {
		return "", false, err
	}

	if len(byPrefix) == 0 {
		return "", false, nil
	}

	return byPrefix[0], true, nil
}
