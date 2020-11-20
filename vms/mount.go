package vms

import (
	"fmt"
	"github.com/mhewedy/vermin/db"
	"github.com/mhewedy/vermin/hypervisor"
	"github.com/mhewedy/vermin/images"
)

var (
	mountFormat = "%-79s%-70s\n"
	mountHeader = fmt.Sprintf(mountFormat, "HOST DIR", "GUEST DIR")
)

func Mount(vmName, hostPath, guestPath string, remove bool) error {
	if err := checkRunningVM(vmName); err != nil {
		return err
	}

	vmdb, err := db.Load(vmName)
	if err != nil {
		return err
	}

	if err = images.CheckCanMount(vmdb.Image); err != nil {
		return err
	}

	if remove {
		if err = hypervisor.RemoveMounts(vmName); err != nil {
			return err
		}
	}

	return hypervisor.AddMount(vmName, hostPath, guestPath)
}

func ListMounts(vmName string) (string, error) {

	out := mountHeader

	paths, err := hypervisor.ListMounts(vmName)
	if err != nil {
		return "", err
	}

	for _, p := range paths {
		out += fmt.Sprintf(mountFormat, p.HostPath, p.GuestPath)
	}

	return out, nil
}
