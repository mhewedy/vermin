package main

import (
	"fmt"
	"vermin/cmd"
	"vermin/db"
	"vermin/ip"
)

func start(vmName string) error {
	fmt.Println("starting", vmName, "...")
	if _, err := cmd.Execute("vboxmanage", "startvm", vmName, "--type", "headless"); err != nil {
		return err
	}
	return nil
}

func stop(vmName string) error {
	fmt.Println("stopping", vmName, "...")
	if _, err := cmd.Execute("vboxmanage", "controlvm", vmName, "poweroff"); err != nil {
		return err
	}
	return nil
}

func establishConn(vmName string) error {

	return nil
}

func ssh(vmName string, command ...string) (string, error) {

	ipAddr, err := ip.Find(vmName, false)
	if err != nil {
		return "", err
	}

	args := []string{
		"-i",
		db.GetPrivateKeyPath(),
		db.GetUsername() + "@" + ipAddr,
	}

	if len(command) > 0 {
		args = append(args, "--")
		args = append(args, command...)
		return cmd.Execute("ssh", args...)
	} else {
		return "", cmd.ExecuteI("ssh", args...)
	}
}
