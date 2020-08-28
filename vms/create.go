package vms

import (
	"bufio"
	"fmt"
	"github.com/mhewedy/vermin/command"
	"github.com/mhewedy/vermin/command/ssh"
	"github.com/mhewedy/vermin/config/trace"
	"github.com/mhewedy/vermin/db"
	"github.com/mhewedy/vermin/images"
	"github.com/mhewedy/vermin/progress"
	"github.com/mhewedy/vermin/provisioners"
	"os"
	"runtime"
	"sort"
	"strconv"
	"strings"
)

type ProvisionScript struct {
	Script string
	Func   provisioners.Func
}

func Create(imageName string, ps ProvisionScript, cpus int, mem int) (string, error) {
	trace.Create(imageName)

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

	// execute import cmd
	importCmd := command.VBoxManage(
		"import", db.GetImageFilePath(imageName),
		"--vsys", "0",
		"--vmname", vmName,
		"--basefolder", db.VMsBaseDir,
		"--cpus", fmt.Sprintf("%d", cpus),
		"--memory", fmt.Sprintf("%d", mem),
	)
	if _, err = importCmd.CallWithProgress(fmt.Sprintf("Creating %s from image %s", vmName, imageName)); err != nil {
		return "", err
	}

	if err := setNetworkAdapter(vmName); err != nil {
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

func setNetworkAdapter(vmName string) error {
	progress.Immediate("Setting bridged network adapter")
	r, err := command.VBoxManage("list", "bridgedifs").Call()
	if err != nil {
		return err
	}

	l, _, err := bufio.NewReader(strings.NewReader(r)).ReadLine()
	if err != nil {
		return err
	}
	adapter := strings.ReplaceAll(string(l), "Name:", "")
	adapter = strings.TrimSpace(adapter)

	if _, err = command.VBoxManage("modifyvm", vmName, "--nic1", "bridged").Call(); err != nil {
		return nil
	}

	if runtime.GOOS == "windows" {
		adapter = fmt.Sprintf(`"%s"`, adapter)
	}
	if _, err := command.VBoxManage("modifyvm", vmName, "--bridgeadapter1", adapter).Call(); err != nil {
		return nil
	}

	return nil
}

func start(vmName string) error {
	progress.Immediate("Starting", vmName)
	if _, err := command.VBoxManage("startvm", vmName, "--type", "headless").Call(); err != nil {
		return err
	}

	if err := ssh.EstablishConn(vmName); err != nil {
		return err
	}
	return nil
}

func buildVMName() (string, error) {
	var curr int

	l, err := List(true)
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
