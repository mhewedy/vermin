package vms

import (
	"fmt"
	"github.com/mhewedy/vermin/hypervisor"
)

func checkRunningVM(vmName string) error {
	list, err := hypervisor.List(false)
	if err != nil {
		return err
	}
	if !contains(list, vmName) {
		return fmt.Errorf("%s not running.\nUse the command 'vermin ps' to list running VMs", vmName)
	}
	return nil
}

func checkVM(vmName string) error {
	list, err := hypervisor.List(true)
	if err != nil {
		return err
	}
	if !contains(list, vmName) {
		return fmt.Errorf("%s not found.\nUse the command 'vermin ps -a' to list VMs", vmName)
	}
	return nil
}

func isRunningVM(vmName string) bool {
	list, _ := hypervisor.List(false)
	if contains(list, vmName) {
		return true
	}
	return false
}

func contains(a []string, s string) bool {
	for i := range a {
		if a[i] == s {
			return true
		}
	}
	return false
}
