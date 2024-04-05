package hypervisor

import (
	"path/filepath"
	"reflect"
	"strings"

	"github.com/mhewedy/vermin/db"
	"github.com/mhewedy/vermin/hypervisor/base"
	"github.com/mhewedy/vermin/hypervisor/virtualbox"
	"github.com/mhewedy/vermin/progress"
)

func detect() (base.Hypervisor, error) {
	h := virtualbox.Instance
	return h, nil
}

func GetHypervisorName(showDetectedMsg bool) (string, error) {
	h, err := detect()
	if err != nil {
		return "", err
	}

	if showDetectedMsg {
		progress.Immediate(reflect.TypeOf(h).Elem().Name(), "hypervisor detected       ")
	}
	return reflect.TypeOf(h).Elem().Name(), nil
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

func GetVMProperty(vmName, property string) (*string, error) {
	h, err := detect()
	if err != nil {
		return nil, err
	}

	return h.GetVMProperty(vmName, property)
}

func HealthCheck(vmName, property string) (*string, error) {
	h, err := detect()
	if err != nil {
		return nil, err
	}

	return h.HealthCheck(vmName, property)
}

func ShowGUI(vmName string) error {
	h, err := detect()
	if err != nil {
		return err
	}

	return h.ShowGUI(vmName)
}

func AddMount(vmName, ipAddr, hostPath, guestPath string) error {
	absHostPath, err := filepath.Abs(hostPath)
	if err != nil {
		return err
	}

	h, err := detect()
	if err != nil {
		return err
	}

	return h.AddMount(vmName, ipAddr, absHostPath, guestPath)
}

func ListMounts(vmName, ipAddr string) ([]base.MountPath, error) {
	h, err := detect()
	if err != nil {
		return nil, err
	}

	return h.ListMounts(vmName, ipAddr)
}

func RemoveMounts(vmName, ipAddr string) error {
	h, err := detect()
	if err != nil {
		return err
	}

	return h.RemoveMounts(vmName, ipAddr)
}

func SetNetworkAdapterAsBridge(vmName string) error {
	h, err := detect()
	if err != nil {
		return err
	}

	return h.SetNetworkAdapterAsBridge(vmName)
}

func GetBoxInfo(vmName string) (*base.Box, error) {
	h, err := detect()
	if err != nil {
		return nil, err
	}

	return h.GetBoxInfo(vmName)
}

func GetSubnet() (*base.Subnet, error) {
	h, err := detect()
	if err != nil {
		return nil, err
	}

	return h.GetSubnet()
}

func ShrinkDisk(vmName string) error {

	h, err := detect()
	if err != nil {
		return err
	}

	return h.ShrinkDisk(vmName)
}
