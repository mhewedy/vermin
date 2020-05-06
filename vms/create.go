package vms

import (
	"bufio"
	"fmt"
	"github.com/mhewedy/vermin/command"
	"github.com/mhewedy/vermin/command/scp"
	"github.com/mhewedy/vermin/command/ssh"
	"github.com/mhewedy/vermin/db"
	"github.com/mhewedy/vermin/images"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strconv"
	"strings"
)

func Create(imageName string, script string, cpus int, mem int) (string, error) {
	if err := images.Create(imageName); err != nil {
		return "", err
	}

	vmName, err := nextName()
	if err != nil {
		return "", err
	}

	if err = os.MkdirAll(db.GetVMPath(vmName), 0755); err != nil {
		return "", err
	}

	if err = db.WriteImageData(vmName, imageName); err != nil {
		return "", err
	}

	// execute import cmd
	importCmd := command.VBoxManage(
		"import", db.GetImageFilePath(imageName),
		"--vsys", "0",
		"--vmname", vmName,
		"--basefolder", db.GetVMsBaseDir(),
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

	if len(script) > 0 {
		if err := provision(vmName, script); err != nil {
			return "", err
		}
	}

	return vmName, nil
}

func setNetworkAdapter(vmName string) error {
	fmt.Println("Setting bridged network adapter ...")
	r, err := command.VBoxManage("list", "bridgedifs").Call()
	if err != nil {
		return err
	}

	reader := bufio.NewReader(strings.NewReader(r))
	l, _, err := reader.ReadLine()
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

func provision(vmName string, script string) error {
	fmt.Println("Provisioning", vmName, "...")

	vmFile := "/tmp/" + filepath.Base(script)
	if err := scp.CopyToVM(vmName, script, vmFile); err != nil {
		return err
	}
	if _, err := ssh.Execute(vmName, "chmod +x "+vmFile); err != nil {
		return err
	}
	if err := ssh.Interact(vmName, vmFile); err != nil {
		return err
	}

	return nil
}

func start(vmName string) error {
	fmt.Println("Starting", vmName, "...")
	if _, err := command.VBoxManage("startvm", vmName, "--type", "headless").Call(); err != nil {
		return err
	}

	if err := ssh.EstablishConn(vmName); err != nil {
		return err
	}
	return nil
}

func nextName() (string, error) {
	var curr int

	l, err := List(true)
	if err != nil {
		return "", err
	}

	if len(l) == 0 {
		curr = 0
	} else {
		sort.Slice(l, func(i, j int) bool {
			ii, _ := strconv.Atoi(strings.ReplaceAll(l[i], db.NamePrefix, ""))
			jj, _ := strconv.Atoi(strings.ReplaceAll(l[j], db.NamePrefix, ""))
			return ii <= jj
		})
		curr, _ = strconv.Atoi(strings.ReplaceAll(l[len(l)-1], db.NamePrefix, ""))
	}

	return fmt.Sprintf(db.NamePrefix+"%02d", curr+1), nil
}
