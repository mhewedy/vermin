package main

import (
	"fmt"
	"io/ioutil"
	"sort"
	"strconv"
	"strings"
	"vermin/cmd"
	"vermin/db"
	"vermin/images"
	"vermin/info"
)

func create(imageName string, provision string, cpus int, mem int) error {
	if err := images.Create(imageName); err != nil {
		return err
	}
	vmName, err := getNextVMName()
	if err != nil {
		return err
	}
	// set defaults
	if cpus == 0 {
		cpus = 1
	}
	if mem == 0 {
		mem = 1024
	}
	// execute command
	fmt.Printf("creating %s from image %s", vmName, imageName)
	if _, err = cmd.ExecuteP("vboxmanage",
		"import", db.GetImageFilePath(imageName),
		"--vsys", "0",
		"--vmname", vmName,
		"--basefolder", db.GetVMsBaseDir(),
		"--cpus", fmt.Sprintf("%d", cpus),
		"--memory", fmt.Sprintf("%d", mem),
	); err != nil {
		return err
	}
	if err = ioutil.WriteFile(db.GetVMPath(vmName)+"/"+db.Image, []byte(imageName), 0775); err != nil {
		return err
	}
	fmt.Printf("vm created: %s\n", vmName)

	return nil
}

func provision(vmName string, script string) error {

	return nil
}

func getNextVMName() (string, error) {
	var curr int

	l, err := info.List(true)
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
