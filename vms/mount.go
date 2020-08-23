package vms

import (
	"errors"
	"fmt"
	"github.com/mhewedy/vermin/command"
	"github.com/mhewedy/vermin/command/ssh"
	"github.com/mhewedy/vermin/db"
	"github.com/mhewedy/vermin/db/props"
	"github.com/mhewedy/vermin/images"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func Mount(vmName, hostPath, guestPath string) error {
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

	if err = removeMounts(vmName); err != nil {
		return err
	}

	absHostPath, err := filepath.Abs(hostPath)
	if err != nil {
		return err
	}

	shareName := strconv.FormatInt(time.Now().Unix(), 10)

	if _, err = command.VBoxManage("sharedfolder",
		"add",
		vmName,
		"--name", shareName,
		"--hostpath", absHostPath,
		"--transient",
		"--automount",
		"--auto-mount-point",
		guestPath).Call(); err != nil {
		return err
	}

	mountCmd := fmt.Sprintf("sudo mkdir -p %s; ", guestPath) +
		fmt.Sprintf("sudo mount -t vboxsf -o uid=1000,gid=1000 %s %s", shareName, guestPath)

	if _, err = ssh.Execute(vmName, mountCmd); err != nil {
		return err
	}

	return nil
}

func removeMounts(vmName string) error {
	transientMounts, err := props.FindByPrefix(vmName, "SharedFolderNameTransientMapping")
	if err != nil {
		return err
	}

	for i := range transientMounts {
		mountName := transientMounts[i]
		if _, err = command.VBoxManage("sharedfolder", "remove", vmName, "--name", mountName, "--transient").Call(); err != nil {
			return err
		}
		if guestPath, err := getMountGuestPath(vmName, mountName); err == nil {
			_, _ = ssh.Execute(vmName, "sudo umount "+guestPath)
		}
	}

	return nil
}

// return guestPath of mount
func getMountGuestPath(vmName, mountName string) (string, error) {
	mountData, _ := ssh.Execute(vmName, "sudo mount")
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
