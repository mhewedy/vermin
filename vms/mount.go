package vms

import (
	"fmt"
	"github.com/mhewedy/vermin/command"
	"github.com/mhewedy/vermin/command/ssh"
	"github.com/mhewedy/vermin/db"
	"github.com/mhewedy/vermin/db/props"
	"github.com/mhewedy/vermin/images"
	"path/filepath"
	"strconv"
	"time"
)

func Mount(vmName string, hostPath string) error {
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

	guestFolder := "/vermin"
	shareName := strconv.FormatInt(time.Now().Unix(), 10)

	if _, err = command.VBoxManage("sharedfolder",
		"add",
		vmName,
		"--name", shareName,
		"--hostpath", absHostPath,
		"--transient",
		"--automount",
		"--auto-mount-point",
		guestFolder).Call(); err != nil {
		return err
	}

	mountCmd := fmt.Sprintf("sudo mkdir -p %s; ", guestFolder) +
		fmt.Sprintf("sudo mount -t vboxsf -o uid=1000,gid=1000 %s %s", shareName, guestFolder)

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
		mount := transientMounts[i]
		if _, err = command.VBoxManage("sharedfolder", "remove", vmName, "--name", mount, "--transient").Call(); err != nil {
			return err
		}
	}

	return nil
}
