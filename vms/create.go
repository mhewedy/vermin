package vms

import (
	"bufio"
	"fmt"
	"github.com/mhewedy/vermin/cmd"
	"github.com/mhewedy/vermin/cmd/scp"
	"github.com/mhewedy/vermin/cmd/ssh"
	"github.com/mhewedy/vermin/db"
	"github.com/mhewedy/vermin/images"
	"io/ioutil"
	"path/filepath"
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
	// execute command
	fmt.Printf("Creating %s from image %s ", vmName, imageName)
	if _, err = cmd.ExecuteP("vboxmanage",
		"import", db.GetImageFilePath(imageName),
		"--vsys", "0",
		"--vmname", vmName,
		"--basefolder", db.GetVMsBaseDir(),
		"--cpus", fmt.Sprintf("%d", cpus),
		"--memory", fmt.Sprintf("%d", mem),
	); err != nil {
		return "", err
	}
	if err = ioutil.WriteFile(db.GetVMPath(vmName)+"/"+db.Image, []byte(imageName), 0775); err != nil {
		return "", err
	}

	if err := SetNetworkAdapter(vmName); err != nil {
		return "", err
	}

	if err := provision(vmName, script); err != nil {
		return "", err
	}

	return vmName, nil
}

func SetNetworkAdapter(vmName string) error {
	fmt.Println("Setting Network adapter ...")
	r, err := cmd.Execute("vboxmanage", "list", "bridgedifs")
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

	if _, err = cmd.Execute("vboxmanage", "modifyvm", vmName, "--nic1", "bridged"); err != nil {
		return nil
	}

	if _, err := cmd.Execute("vboxmanage", "modifyvm", vmName, "--bridgeadapter1", fmt.Sprintf(`"%s"`, adapter)); err != nil {
		return nil
	}

	return nil
}

func provision(vmName string, script string) error {
	if len(script) == 0 {
		return nil
	}
	fmt.Println("Provisioning", vmName, "...")
	if err := Start(vmName); err != nil {
		return err
	}
	if err := ssh.EstablishConn(vmName); err != nil {
		return err
	}

	vmFile := "/tmp/" + filepath.Base(script)
	if err := scp.CopyToVM(vmName, script, vmFile); err != nil {
		return err
	}
	if _, err := ssh.Execute(vmName, "chmod +x "+vmFile); err != nil {
		return err
	}
	if err := ssh.ExecuteI(vmName, vmFile); err != nil {
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
