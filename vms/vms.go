package vms

import (
	"errors"
	"fmt"
	"github.com/mhewedy/vermin/command"
	"github.com/mhewedy/vermin/command/scp"
	"github.com/mhewedy/vermin/command/ssh"
	"github.com/mhewedy/vermin/db"
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

	progress.Immediate("Stopping", vmName)
	if _, err := command.VBoxManage("controlvm", vmName, "poweroff").Call(); err != nil {
		return err
	}
	return nil
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

	msg := fmt.Sprintf("Removing %s", vmName)
	if _, err := command.VBoxManage("unregistervm", vmName, "--delete").CallWithProgress(msg); err != nil {
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

	var params = []string{"modifyvm", vmName}
	if cpus > 0 {
		params = append(params, "--cpus", fmt.Sprintf("%d", cpus))
	}
	if mem > 0 {
		params = append(params, "--memory", fmt.Sprintf("%d", mem))
	}

	_, err := command.VBoxManage(params...).Call()
	return err
}

func GUI(vmName string) error {
	if err := checkRunningVM(vmName); err != nil {
		return err
	}

	if err := ssh.EstablishConn(vmName); err != nil {
		return err
	}

	return command.VBoxManage("startvm", "--type", "separate", vmName).Run()
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
