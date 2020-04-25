package main

import (
	"fmt"
	"os"
	"vermin/cmd"
	"vermin/cmd/ssh"
	"vermin/db"
	"vermin/info"
	"vermin/port"
)

func ps(all bool) (string, error) {
	vms, err := info.List(all)
	if err != nil {
		return "", err
	}
	return info.Get(vms), nil
}

func start(vmName string) error {
	fmt.Println("Starting", vmName, "...")
	if _, err := cmd.Execute("vboxmanage", "startvm", vmName, "--type", "headless"); err != nil {
		return err
	}
	return nil
}

func stop(vmName string) error {
	fmt.Println("Stopping", vmName, "...")
	if _, err := cmd.Execute("vboxmanage", "controlvm", vmName, "poweroff"); err != nil {
		return err
	}
	return nil
}

func remove(vmName string) error {
	_ = stop(vmName)

	fmt.Println("Removing", vmName, "...")
	if _, err := cmd.Execute("vboxmanage", "unregistervm", vmName, "--delete"); err != nil {
		return err
	}

	return os.RemoveAll(db.GetVMPath(vmName))
}

func secureShell(vmName string, command string) error {
	if err := ssh.EstablishConn(vmName); err != nil {
		return err
	}

	if len(command) == 0 {
		return ssh.Shell(vmName)
	} else {
		return ssh.ExecuteI(vmName, command)
	}
}

func portForward(vmName string, ports string) error {
	a, err := port.MapPortForward(ports)
	if err != nil {
		return err
	}

	if err := ssh.EstablishConn(vmName); err != nil {
		return err
	}

	fmt.Println("Connected. Press CTRL+C anytime to stop")
	if err := ssh.ExecuteIArgs(vmName, append(a, "-N")...); err != nil {
		return err
	}

	return nil
}
