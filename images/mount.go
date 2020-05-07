package images

import (
	"fmt"
	"github.com/mhewedy/vermin/command"
	"github.com/mhewedy/vermin/db"
	"github.com/mhewedy/vermin/db/info"
	"strconv"
	"strings"
	"time"
)

// should check for vm is running
// should check for valid hostPath
func Mount(vmName string, hostPath string) error {

	image, err := db.ReadImageData(vmName)
	if err != nil {
		return err
	}

	if err = checkCanMount(image); err != nil {
		return err
	}

	if err = removeMounts(vmName); err != nil {
		return err
	}

	if _, err = command.VBoxManage("sharedfolder",
		"add",
		vmName,
		"--name", strconv.FormatInt(time.Now().Unix(), 10),
		"--hostpath", hostPath,
		"--transient",
		"--automount",
		"--auto-mount-point", "/vermin").Call(); err != nil {
		return err
	}

	return nil
}

func checkCanMount(image string) error {
	remote, _ := listRemoteImages(false)

	dbImage, err := remote.findByName(image)
	if err != nil {
		return err
	}

	if !dbImage.Mount {
		mounted := remote.findByMount(true).names()
		return fmt.Errorf("image '%s' cannot be mounted, "+
			"images can be mounted are:\n%s", image, strings.Join(mounted, "\n"))
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
