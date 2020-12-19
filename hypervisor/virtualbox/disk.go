package virtualbox

import (
	"fmt"
	"github.com/mhewedy/vermin/db"
	"github.com/mhewedy/vermin/hypervisor/base"
	"github.com/mhewedy/vermin/log"
	"github.com/mhewedy/vermin/progress"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func (*virtualbox) ShrinkDisk(vmName string) error {

	stop := progress.Show("Shrinking disk", false)
	defer stop()

	box, err := getBoxInfo(vmName)
	if err != nil {
		return err
	}

	origDiskPath := filepath.Join(db.GetVMPath(vmName), box.Disk.Location)

	if isVMDK(box.Disk) {

		tmpDir, err := ioutil.TempDir("", "vermin_disk_shrink_"+vmName)
		if err != nil {
			return err
		}
		defer os.RemoveAll(tmpDir)

		vdiPath := filepath.Join(tmpDir, box.Disk.Location+".vdi")
		newVMDKPath := filepath.Join(tmpDir, box.Disk.Location)

		log.Debug("clone vmdk disk into vdi")
		if err := vboxManage("clonehd", origDiskPath, vdiPath, "--format", "vdi").Run(); err != nil {
			return err
		}
		log.Debug("shrink the vdi")
		if err := vboxManage("modifyhd", vdiPath, "--compact").Run(); err != nil {
			return err
		}
		log.Debug("clone the vdi back into a new vmdk")
		if err := vboxManage("clonehd", vdiPath, newVMDKPath, "--format", "vmdk").Run(); err != nil {
			return err
		}
		log.Debug("set the uuid of the new vmdk with the same uuid of the old vmdk")
		if err := vboxManage("internalcommands", "sethduuid", newVMDKPath, box.Disk.UUID).Run(); err != nil {
			return err
		}
		log.Debug("copy the new vmdk into the same location of the old vmdk")
		return copyFile(newVMDKPath, origDiskPath)

	} else if isVDI(box.Disk) {
		return vboxManage("modifyhd", origDiskPath, "--compact").Run()
	} else {
		return fmt.Errorf("unsupported disk format %s", box.Disk.Location)
	}
}

func isVMDK(disk *base.Disk) bool {
	return strings.HasSuffix(strings.ToLower(disk.Location), ".vmdk")
}

func isVDI(disk *base.Disk) bool {
	return strings.HasSuffix(strings.ToLower(disk.Location), ".vdi")
}

// Copy the src file to dst. Any existing file will be overwritten and will not
// copy file attributes.
func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}
	return out.Close()
}
