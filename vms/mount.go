package vms

import (
	"github.com/mhewedy/vermin/command"
	"github.com/mhewedy/vermin/db"
	"github.com/mhewedy/vermin/db/info"
	"github.com/mhewedy/vermin/images"
	"path/filepath"
	"strconv"
	"time"
)

func Mount(vmName string, hostPath string) error {
	if err := checkRunningVM(vmName); err != nil {
		return err
	}

	image, err := db.ReadImageData(vmName)
	if err != nil {
		return err
	}

	if err = images.CheckCanMount(image); err != nil {
		return err
	}

	if err = removeMounts(vmName); err != nil {
		return err
	}

	absHostPath, err := filepath.Abs(hostPath)
	if err != nil {
		return err
	}

	if _, err = command.VBoxManage("sharedfolder",
		"add",
		vmName,
		"--name", strconv.FormatInt(time.Now().Unix(), 10),
		"--hostpath", absHostPath,
		"--transient",
		"--automount",
		"--auto-mount-point", "/vermin").Call(); err != nil {
		return err
	}

	return nil
}

func removeMounts(vmName string) error {
	transientMounts, err := info.FindByPrefix(vmName, "SharedFolderNameTransientMapping")
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
