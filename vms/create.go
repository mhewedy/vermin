package vms

import (
	"fmt"
	"os"
	"runtime"
	"sort"
	"strconv"
	"strings"

	"github.com/mhewedy/vermin/cmd/ssh"
	"github.com/mhewedy/vermin/db"
	"github.com/mhewedy/vermin/hypervisor"
	"github.com/mhewedy/vermin/images"
	"github.com/mhewedy/vermin/progress"
	"github.com/mhewedy/vermin/provisioners"
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

	if err := changeHostname(vmName); err != nil {
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

func changeHostname(vmName string) error {

	stop := progress.Show("Setting hostname", false)
	defer stop()

	oldHostname, err := ssh.Execute(vmName, "hostname")
	if err != nil {
		return err
	}

	oldHostname = strings.TrimSuffix(oldHostname, "\n")
	oldHostname = strings.ReplaceAll(oldHostname, ".localdomain", "")

	newHostname := strings.Split(vmName, "_")[1]

	var cmds []string

	cmds = append(cmds, fmt.Sprintf("sudo hostname %s", newHostname))

	if runtime.GOOS == "windows" {
		cmds = append(cmds,
			"'echo "+newHostname+" | sudo tee /tmp/new_hostname > /dev/null'",
			"sudo mv /tmp/new_hostname /etc/hostname",
		)
	} else {
		cmds = append(cmds,
			fmt.Sprintf("sudo sh -c 'echo %s > /etc/hostname'", newHostname),
		)
	}

	cmds = append(cmds, fmt.Sprintf("sudo sed -i 's/%s/%s/g' /etc/hosts", oldHostname, newHostname))

	for _, cmd := range cmds {
		if _, err := ssh.Execute(vmName, cmd); err != nil {
			return err
		}
	}

	return nil
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
