package cli

import "github.com/mhewedy/vermin/commands/info"

func getStoppedVms() []string {
	stopped := make([]string, 0)
	all, _ := info.List(true)
	running, _ := info.List(false)

	for i := range all {
		vm := all[i]
		if !contains(running, vm) {
			stopped = append(stopped, vm)
		}
	}
	return stopped
}

func contains(a []string, s string) bool {
	for i := range a {
		if a[i] == s {
			return true
		}
	}
	return false
}
