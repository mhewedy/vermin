package vms

import (
	"fmt"
	"github.com/mhewedy/vermin/cmd"
	"github.com/mhewedy/vermin/cmd/scp"
	"github.com/mhewedy/vermin/cmd/ssh"
	"github.com/mhewedy/vermin/db"
	"os"
)

func Tag(vmName string, tag string) error {
	return appendToFile(db.GetVMPath(vmName)+"/"+db.Tags, []byte(tag+"\n"), 0775)
}

func Start(vmName string) error {
	fmt.Println("Starting", vmName, "...")
	if _, err := cmd.Execute("vboxmanage", "startvm", vmName, "--type", "headless"); err != nil {
		return err
	}
	return nil
}

func Stop(vmName string) error {
	fmt.Println("Stopping", vmName, "...")
	if _, err := cmd.Execute("vboxmanage", "controlvm", vmName, "poweroff"); err != nil {
		return err
	}
	return nil
}

func SecureShell(vmName string, command string) error {
	if err := ssh.EstablishConn(vmName); err != nil {
		return err
	}
	if len(command) == 0 {
		return ssh.OpenTerminal(vmName)
	} else {
		return ssh.ExecuteI(vmName, command)
	}
}

func Remove(vmName string) error {
	_ = Stop(vmName)
	fmt.Println("Removing", vmName, "...")
	if _, err := cmd.Execute("vboxmanage", "unregistervm", vmName, "--delete"); err != nil {
		return err
	}
	return os.RemoveAll(db.GetVMPath(vmName))
}

func PortForward(vmName string, ports string) error {
	a, err := getPortForwardArgs(ports)
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

func CopyFiles(vmName string, file string, toVM bool) error {
	if toVM {
		return scp.CopyToVMHomeDir(vmName, file)
	} else {
		return scp.CopyToLocalCWD(vmName, file)
	}
}
