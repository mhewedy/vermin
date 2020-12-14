package vms

import (
	"errors"
	"fmt"
	"github.com/mhewedy/vermin/cmd/scp"
	"github.com/mhewedy/vermin/cmd/ssh"
	"github.com/mhewedy/vermin/db"
	"github.com/mhewedy/vermin/hypervisor"
	"github.com/mhewedy/vermin/images"
	"github.com/mhewedy/vermin/ip"
	"github.com/mhewedy/vermin/progress"
	"os"
	"strings"
)

func Tag(vmName string, tag string, remove bool) error {
	if err := checkVM(vmName); err != nil {
		return err
	}
	if remove {
		return db.RemoveTag(vmName, tag)
	} else {
		return db.AddTag(vmName, tag)
	}
}

func Start(vmName string) error {
	if err := checkVM(vmName); err != nil {
		return err
	}

	if isRunningVM(vmName) {
		return fmt.Errorf(`VM already running, use "vermin ssh %s" to ssh into the VM.`, vmName)
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

	if !isRunningVM(vmName) {
		return fmt.Errorf(`VM already stopped, use "vermin start %s" to start the VM.`, vmName)
	}

	progress.Immediate("Stopping", vmName)
	if err := hypervisor.Stop(vmName); err != nil {
		return err
	}
	return nil
}

func Restart(vmName string) error {
	if err := Stop(vmName); err != nil {
		return nil
	}
	return start(vmName)
}

func SecureShell(vmName string) error {
	if err := checkRunningVM(vmName); err != nil {
		return err
	}

	if err := ssh.EstablishConn(vmName); err != nil {
		return err
	}

	return ssh.OpenTerminal(vmName)
}

func Exec(vmName string, cmd string) error {
	if err := checkRunningVM(vmName); err != nil {
		return err
	}

	if err := ssh.EstablishConn(vmName); err != nil {
		return err
	}

	return ssh.ExecInteract(vmName, cmd)
}

func Remove(vmName string, force bool) error {
	if !force && isRunningVM(vmName) {
		return errors.New("Cannot stop running VM, use -f flag to force remove")
	}

	if err := checkVM(vmName); err != nil {
		return err
	}
	_ = Stop(vmName)

	if err := hypervisor.Remove(vmName); err != nil {
		return err
	}
	return os.RemoveAll(db.GetVMPath(vmName))
}

func PortForward(vmName string, ports string) error {
	if err := checkRunningVM(vmName); err != nil {
		return err
	}

	a, err := getPortForwardArgs(vmName, ports)
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

func CopyFiles(src string, dest string) error {

	if err := checkRunningVMForCopy(src); err != nil {
		return err
	}
	if err := checkRunningVMForCopy(dest); err != nil {
		return err
	}

	return scp.Copy(src, dest)
}

func checkRunningVMForCopy(srcDest string) error {
	if strings.Contains(srcDest, scp.CopySeparator) {
		vmName := strings.Split(srcDest, scp.CopySeparator)[0]
		if err := checkRunningVM(vmName); err != nil {
			return err
		}
	}
	return nil
}

func IP(vmName string, purge bool, global bool) (string, error) {
	if !global {
		if err := checkRunningVM(vmName); err != nil {
			return "", err
		}
	}
	return ip.Find(vmName, purge)
}

func Modify(vmName string, cpus int, mem int) error {
	if isRunningVM(vmName) {
		return fmt.Errorf(`Cannot Modify running VM, use "vermin stop %s" to stop the VM first.`, vmName)
	}

	return hypervisor.Modify(vmName, cpus, mem)
}

func GUI(vmName string, nocheck bool) error {
	if !nocheck {
		if err := checkRunningVM(vmName); err != nil {
			return err
		}

		if err := ssh.EstablishConn(vmName); err != nil {
			return err
		}
	}
	return hypervisor.ShowGUI(vmName)
}

func Commit(vmName, imageName string, override bool) error {
	if err := checkVM(vmName); err != nil {
		return err
	}

	if isRunningVM(vmName) {
		return fmt.Errorf(`VM is running, use "vermin stop %s" to stop the VM before commiting image from it.`, vmName)
	}

	return images.Commit(vmName, imageName, override)
}
