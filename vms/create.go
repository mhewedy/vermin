package vms

import (
	"fmt"
	"github.com/mhewedy/vermin/cmd/ssh"
	"github.com/mhewedy/vermin/db"
	"github.com/mhewedy/vermin/hypervisor"
	"github.com/mhewedy/vermin/images"
	"github.com/mhewedy/vermin/progress"
	"github.com/mhewedy/vermin/provisioners"
	"os"
	"sort"
	"strconv"
	"strings"
)

type ProvisionScript struct {
	Script string
	Func   provisioners.Func
}

func Create(imageName string, ps ProvisionScript, cpus int, mem int) (string, error) {

	if err := images.Download(imageName); err != nil {
		return "", err
	}

	vmName, err := buildVMName()
	if err != nil {
		return "", err
	}

	if err = os.RemoveAll(db.GetVMPath(vmName)); err != nil {
		return "", err
	}
	if err = os.MkdirAll(db.GetVMPath(vmName), 0755); err != nil {
		return "", err
	}

	if err = db.SetImage(vmName, imageName); err != nil {
		return "", err
	}

	if err = hypervisor.Create(imageName, vmName, cpus, mem); err != nil {
		return "", err
	}

	progress.Immediate("Setting bridged network adapter")
	if err := hypervisor.SetNetworkAdapterAsBridge(vmName); err != nil {
		return "", err
	}

	if err := start(vmName); err != nil {
		return "", err
	}

	if len(ps.Script) > 0 {
		fmt.Println("Provisioning", vmName)
		if err := ps.Func(vmName, ps.Script); err != nil {
			return "", err
		}
	}

	return vmName, nil
}

func start(vmName string) error {
	progress.Immediate("Starting", vmName)
	if err := hypervisor.Start(vmName); err != nil {
		return err
	}

	if err := ssh.EstablishConn(vmName); err != nil {
		return err
	}
	return nil
}

func buildVMName() (string, error) {
	var curr int

	l, err := hypervisor.List(true)
	if err != nil {
		return "", err
	}

	if len(l) == 0 {
		curr = 0
	} else {
		sort.Slice(l, func(i, j int) bool {
			ii, _ := strconv.Atoi(strings.ReplaceAll(l[i], db.VMNamePrefix, ""))
			jj, _ := strconv.Atoi(strings.ReplaceAll(l[j], db.VMNamePrefix, ""))
			return ii <= jj
		})
		curr, _ = strconv.Atoi(strings.ReplaceAll(l[len(l)-1], db.VMNamePrefix, ""))
	}

	return fmt.Sprintf(db.VMNamePrefix+"%02d", curr+1), nil
}
