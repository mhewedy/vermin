package virtualbox

import (
	"bufio"
	"errors"
	"fmt"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/mhewedy/vermin/cmd"
	"github.com/mhewedy/vermin/db"
	"github.com/mhewedy/vermin/hypervisor/base"
)

const (
	ipProperty = "/VirtualBox/GuestInfo/Net/0/V4/IP"
)

var propertyMap = map[string]string{
	"ip": ipProperty,
}

var Instance = &virtualbox{}

// should be private by the end of the day
func vboxManage(args ...string) *cmd.Cmd {
	return &cmd.Cmd{
		Command: "vboxmanage",
		Args:    args,
	}
}

type virtualbox struct {
}

func (*virtualbox) Start(vmName string) error {
	return vboxManage("startvm", vmName, "--type", "headless").Run()
}

func (*virtualbox) Commit(vmName, imageName string) error {
	export := vboxManage("export", vmName, "--ovf20", "-o", imageName)
	_, err := export.CallWithProgress(fmt.Sprintf("Committing %s into image %s", vmName, imageName))
	return err
}

func (*virtualbox) Create(imageName, vmName string, cpus int, mem int) error {
	importCmd := vboxManage(
		"import", db.GetImageFilePath(imageName),
		"--vsys", "0",
		"--vmname", vmName,
		"--basefolder", db.VMsBaseDir,
		"--cpus", fmt.Sprintf("%d", cpus),
		"--memory", fmt.Sprintf("%d", mem),
	)

	if runtime.GOOS == "windows" {
		// Wrap image file path and VMs base directory in double quotes
		importCmd.Args[1] = `"` + importCmd.Args[1] + `"`
		importCmd.Args[7] = `"` + importCmd.Args[7] + `"`
	}

	if _, err := importCmd.CallWithProgress(fmt.Sprintf("Creating %s from image %s", vmName, imageName)); err != nil {
		return err
	}

	return nil
}

func (*virtualbox) List(all bool, excludeFunc func(string) bool) ([]string, error) {
	var args = [2]string{"list"}
	if all {
		args[1] = "vms"
	} else {
		args[1] = "runningvms"
	}

	r, err := vboxManage(args[:]...).Call()
	if err != nil {
		return nil, err
	}

	var vms []string
	fields := strings.Fields(r)

	for i := range fields {
		if i%2 == 0 {
			vmName := strings.ReplaceAll(fields[i], `"`, "")
			if !excludeFunc(vmName) {
				vms = append(vms, vmName)
			}
		}
	}

	return vms, nil
}

func (*virtualbox) Stop(vmName string) error {
	if _, err := vboxManage("controlvm", vmName, "poweroff").Call(); err != nil {
		return err
	}

	return nil
}

func (*virtualbox) Remove(vmName string) error {
	msg := fmt.Sprintf("Removing %s", vmName)
	if _, err := vboxManage("unregistervm", vmName, "--delete").CallWithProgress(msg); err != nil {
		return err
	}

	return nil
}

func (*virtualbox) Modify(vmName string, cpus int, mem int) error {
	var params = []string{"modifyvm", vmName}
	if cpus > 0 {
		params = append(params, "--cpus", fmt.Sprintf("%d", cpus))
	}
	if mem > 0 {
		params = append(params, "--memory", fmt.Sprintf("%d", mem))
	}

	if _, err := vboxManage(params...).Call(); err != nil {
		return err
	}

	return nil
}

func (*virtualbox) GetVMProperty(vmName, property string) (*string, error) {
	prop, ok := propertyMap[property]
	if !ok {
		return nil, fmt.Errorf("property %s not found", property)
	}

	guestProperty, err := vboxManage("guestproperty",
		"get",
		vmName,
		prop).Call()

	if err != nil {
		return nil, err
	}

	fields := strings.Fields(guestProperty)
	ipAddress := ""
	if len(fields) >= 2 {
		ipAddress = fields[1]
	}
	return &ipAddress, nil
}

func (*virtualbox) ShowGUI(vmName string) error {
	return vboxManage("startvm", "--type", "separate", vmName).Run()
}

func (*virtualbox) AddMount(vmName, ipAddr, hostPath, guestPath string) error {

	shareName := strconv.FormatInt(time.Now().Unix(), 10)

	if _, err := vboxManage("sharedfolder",
		"add",
		vmName,
		"--name", shareName,
		"--hostpath", hostPath,
		"--transient",
		"--automount",
		"--auto-mount-point",
		guestPath).Call(); err != nil {
		return err
	}

	mountCmd := fmt.Sprintf("sudo mkdir -p %s; ", guestPath) +
		fmt.Sprintf("sudo mount -t vboxsf -o uid=1000,gid=1000 %s %s", shareName, guestPath)

	if err := cmd.Ssh(ipAddr, "--", mountCmd).Run(); err != nil {
		return err
	}

	return nil
}

func (*virtualbox) RemoveMounts(vmName, ipAddr string) error {
	transientMounts, err := findByPrefix(vmName, "SharedFolderNameTransientMapping")
	if err != nil {
		return err
	}

	for i := range transientMounts {
		mountName := transientMounts[i]
		if _, err = vboxManage("sharedfolder", "remove", vmName, "--name", mountName, "--transient").Call(); err != nil {
			return err
		}
		if guestPath, err := getMountGuestPath(ipAddr, mountName); err == nil {
			_ = cmd.Ssh(ipAddr, "--", "sudo umount "+guestPath).Run()
		}
	}

	return nil
}

func (*virtualbox) ListMounts(vmName, ipAddr string) ([]base.MountPath, error) {
	result := make([]base.MountPath, 0)

	hostPaths, err := findByPrefix(vmName, "SharedFolderPathTransientMapping")
	if err != nil {
		return nil, err
	}

	transientMounts, err := findByPrefix(vmName, "SharedFolderNameTransientMapping")
	if err != nil {
		return nil, err
	}

	for i := range transientMounts {
		if guestPath, err := getMountGuestPath(ipAddr, transientMounts[i]); err == nil {
			result = append(result, base.MountPath{
				HostPath:  hostPaths[i],
				GuestPath: guestPath,
			})
		}
	}

	return result, nil
}

// return guestPath of mount
func getMountGuestPath(ipAddr, mountName string) (string, error) {
	mountData, _ := cmd.Ssh(ipAddr, "--", "sudo mount").Call()
	mounts := strings.Split(mountData, "\n")

	for _, line := range mounts {
		if strings.HasPrefix(line, mountName) {
			fields := strings.Fields(line)
			if len(fields) > 3 {
				return fields[2], nil
			}
		}
	}
	return "", errors.New("cannot get mount path")
}

func (*virtualbox) SetNetworkAdapterAsBridge(vmName string) error {
	r, err := vboxManage("list", "bridgedifs").Call()
	if err != nil {
		return err
	}

	l, _, err := bufio.NewReader(strings.NewReader(r)).ReadLine()
	if err != nil {
		return err
	}
	adapter := strings.ReplaceAll(string(l), "Name:", "")
	adapter = strings.TrimSpace(adapter)

	if _, err = vboxManage("modifyvm", vmName, "--nic1", "bridged").Call(); err != nil {
		return nil
	}

	if runtime.GOOS == "windows" {
		adapter = fmt.Sprintf(`"%s"`, adapter)
	}
	if _, err := vboxManage("modifyvm", vmName, "--bridgeadapter1", adapter).Call(); err != nil {
		return nil
	}

	return nil
}

func (*virtualbox) GetBoxInfo(vmName string) (*base.Box, error) {
	return getBoxInfo(vmName)
}

func (*virtualbox) GetSubnet() (*base.Subnet, error) {

	bridgeInfo, err := findBridgeInfo("IPAddress", "NetworkMask")
	if err != nil {
		return nil, err
	}

	return base.NewSubnet(bridgeInfo[0], bridgeInfo[1])
}
