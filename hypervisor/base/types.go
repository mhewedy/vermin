package base

type Hypervisor interface {
	Start(vmName string) error

	Commit(vmName, imageName string) error

	Info(vmName string) ([]string, error)

	Create(imageName, vmName string, cpus int, mem int) error

	List(all bool, excludeFunc func(vmName string) bool) ([]string, error)

	Stop(vmName string) error

	Remove(vmName string) error

	Modify(vmName string, cpus int, mem int) error

	ShowGUI(vmName string) error

	AddMount(vmName, ipAddr, hostPath, guestPath string) error

	RemoveMounts(vmName, ipAddr string) error

	ListMounts(vmName, ipAddr string) ([]MountPath, error)

	SetNetworkAdapterAsBridge(vmName string) error

	GetBoxInfo(vmName string) (*Box, error)
}

type MountPath struct {
	HostPath  string
	GuestPath string
}

type Box struct {
	CPU      string
	Mem      string
	DiskSize string
	MACAddr  string
}
