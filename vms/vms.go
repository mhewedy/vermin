package vms

import (
	"errors"
	"fmt"
	"github.com/mhewedy/vermin/command"
	"github.com/mhewedy/vermin/command/scp"
	"github.com/mhewedy/vermin/command/ssh"
	"github.com/mhewedy/vermin/db"
	"github.com/mhewedy/vermin/ip"
	"os"
)

func Tag(vmName string, tag string) error {
	if err := checkVM(vmName); err != nil {
		return err
	}
	return db.WriteTag(vmName, tag)
}

func Start(vmName string) error {
	if err := checkVM(vmName); err != nil {
		return err
	}

	if err := start(vmName); err != nil {
		return err
	}
	return nil
}

func Stop(vmName string) error {
	if err := checkVM(vmName); err != nil {
		return err
	}

	fmt.Println("Stopping", vmName)
	if _, err := command.VBoxManage("controlvm", vmName, "poweroff").Call(); err != nil {
		return err
	}
	return nil
}

func SecureShell(vmName string, cmd string) error {
	if err := checkRunningVM(vmName); err != nil {
		return err
	}

	if err := ssh.EstablishConn(vmName); err != nil {
		return err
	}

	if len(cmd) == 0 {
		return ssh.OpenTerminal(vmName)
	} else {
		return ssh.Interact(vmName, cmd)
	}
}

func Remove(vmName string, force bool) error {
	if !force && isRunningVM(vmName) {
		return errors.New("Cannot stop running VM, use -f flag to force remove")
	}

	if err := checkVM(vmName); err != nil {
		return err
	}
	_ = Stop(vmName)
	fmt.Println("Removing", vmName)
	if _, err := command.VBoxManage("unregistervm", vmName, "--delete").Call(); err != nil {
		return err
	}
	return os.RemoveAll(db.GetVMPath(vmName))
}

func PortForward(vmName string, ports string) error {
	if err := checkRunningVM(vmName); err != nil {
		return err
	}

	a, err := getPortForwardArgs(ports)
	if err != nil {
		return err
	}

	if err := ssh.EstablishConn(vmName); err != nil {
		return err
	}

	fmt.Println("Connected. Press CTRL+C anytime to stop")
	if err := ssh.WithArgs(vmName, append(a, "-N")); err != nil {
		return err
	}

	return nil
}

func CopyFiles(vmName string, file string, toVM bool) error {
	if err := checkRunningVM(vmName); err != nil {
		return err
	}

	if toVM {
		return scp.CopyToVMHomeDir(vmName, file)
	} else {
		return scp.CopyToLocalCWD(vmName, file)
	}
}

func IP(vmName string, purge bool, global bool) (string, error) {
	if !global {
		if err := checkRunningVM(vmName); err != nil {
			return "", err
		}
	}
	return ip.Find(vmName, purge)
}
