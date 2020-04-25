package main

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
	"vermin/cmd"
	"vermin/cmd/ssh"
	"vermin/db"
	"vermin/images"
	"vermin/info"
)

func create(imageName string, script string, cpus int, mem int) (string, error) {
	if err := images.Create(imageName); err != nil {
		return "", err
	}
	vmName, err := nextName()
	if err != nil {
		return "", err
	}
	// set defaults
	if cpus == 0 {
		cpus = 1
	}
	if mem == 0 {
		mem = 1024
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

	if err := provision(vmName, script); err != nil {
		return "", err
	}

	return vmName, nil
}

func provision(vmName string, script string) error {
	if len(script) == 0 {
		return nil
	}

	fmt.Println("Provisioning", vmName, "...")
	time.Sleep(1 * time.Second)

	if err := start(vmName); err != nil {
		return err
	}
	if err := establishConn(vmName); err != nil {
		return err
	}

	vmFile := "/tmp/" + filepath.Base(script)

	if err := copyToVM(vmName, script, vmFile); err != nil {
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

	l, err := info.List(true)
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
