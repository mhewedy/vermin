package vms

import (
	"fmt"

	"github.com/mhewedy/vermin/hypervisor"
	"github.com/mhewedy/vermin/ip"
)

var (
	mountFormat = "%-40s%s\n"
	mountHeader = fmt.Sprintf(mountFormat, "HOST DIR", "GUEST DIR")
)

func Unmount(vmName string) error {
	ipAddr, err := ip.GetIpAddress(vmName)
	if err != nil {
		return err
	}

	if err := checkRunningVM(vmName); err != nil {
		return err
	}

	return hypervisor.RemoveMounts(vmName, ipAddr)
}

func Mount(vmName, hostPath, guestPath string) error {
	ipAddr, err := ip.GetIpAddress(vmName)
	if err != nil {
		return err
	}

	if err := checkRunningVM(vmName); err != nil {
		return err
	}

	return hypervisor.AddMount(vmName, ipAddr, hostPath, guestPath)
}

func ListMounts(vmName string) (string, error) {
	ipAddr, err := ip.GetIpAddress(vmName)
	if err != nil {
		return "", err
	}

	out := mountHeader

	paths, err := hypervisor.ListMounts(vmName, ipAddr)
	if err != nil {
		return "", err
	}

	for _, p := range paths {
		out += fmt.Sprintf(mountFormat, p.HostPath, p.GuestPath)
	}

	return out, nil
}
