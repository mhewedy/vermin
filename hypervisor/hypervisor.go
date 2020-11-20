package hypervisor

import (
	"github.com/mhewedy/vermin/db"
	"github.com/mhewedy/vermin/hypervisor/base"
	"github.com/mhewedy/vermin/hypervisor/virtualbox"
	"path/filepath"
	"strings"
)

func detect() (base.Hypervisor, error) {
	return virtualbox.Instance, nil
}

func Start(vmName string) error {
	h, err := detect()
	if err != nil {
		return err
	}

	return h.Start(vmName)
}

func Commit(vmName, imageName string) error {
	h, err := detect()
	if err != nil {
		return err
	}

	return h.Commit(vmName, imageName)
}

func Info(vmName string) ([]string, error) {
	h, err := detect()
	if err != nil {
		return nil, err
	}

	return h.Info(vmName)
}

func Create(imageName, vmName string, cpus int, mem int) error {
	h, err := detect()
	if err != nil {
		return err
	}

	return h.Create(imageName, vmName, cpus, mem)
}

func List(all bool) ([]string, error) {
	h, err := detect()
	if err != nil {
		return nil, err
	}

	return h.List(all, func(vmName string) bool {
		return !strings.HasPrefix(vmName, db.VMNamePrefix)
	})
}

func Stop(vmName string) error {
	h, err := detect()
	if err != nil {
		return err
	}

	return h.Stop(vmName)
}

func Remove(vmName string) error {
	h, err := detect()
	if err != nil {
		return err
	}

	return h.Remove(vmName)
}

func Modify(vmName string, cpus int, mem int) error {
	h, err := detect()
	if err != nil {
		return err
	}

	return h.Modify(vmName, cpus, mem)
}

func ShowGUI(vmName string) error {
	h, err := detect()
	if err != nil {
		return err
	}

	return h.ShowGUI(vmName)
}

func AddMount(vmName, hostPath, guestPath string) error {
	absHostPath, err := filepath.Abs(hostPath)
	if err != nil {
		return err
	}

	h, err := detect()
	if err != nil {
		return err
	}

	return h.AddMount(vmName, absHostPath, guestPath)
}

func ListMounts(vmName string) ([]base.MountPath, error) {
	h, err := detect()
	if err != nil {
		return nil, err
	}

	return h.ListMounts(vmName)
}

func RemoveMounts(vmName string) error {
	h, err := detect()
	if err != nil {
		return err
	}

	return h.RemoveMounts(vmName)
}

func SetNetworkAdapterAsBridge(vmName string) error {
	h, err := detect()
	if err != nil {
		return err
	}

	return h.SetNetworkAdapterAsBridge(vmName)
}
