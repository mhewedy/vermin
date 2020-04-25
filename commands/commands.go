package commands

import (
	"fmt"
	"github.com/mhewedy/vermin/cmd"
	"github.com/mhewedy/vermin/cmd/ssh"
	"github.com/mhewedy/vermin/commands/images"
	"github.com/mhewedy/vermin/commands/info"
	"github.com/mhewedy/vermin/commands/ip"
	"github.com/mhewedy/vermin/commands/port"
	"github.com/mhewedy/vermin/db"
	"io"
	"os"
)

func Ps(all bool) (string, error) {
	vms, err := info.List(all)
	if err != nil {
		return "", err
	}
	return info.Get(vms), nil
}

func Images() (string, error) {
	return images.List()
}

func Ip(vmName string, purge bool) (string, error) {
	return ip.Find(vmName, purge)
}

func Create(imageName string, script string, cpus int, mem int) (string, error) {
	return create(imageName, script, cpus, mem)
}

func Tag(vmName string, tag string) error {
	data := []byte(tag + "\n")
	f, err := os.OpenFile(db.GetVMPath(vmName)+"/"+db.Tags, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0775)
	if err != nil {
		return err
	}
	n, err := f.Write(data)
	if err == nil && n < len(data) {
		err = io.ErrShortWrite
	}
	if err1 := f.Close(); err == nil {
		err = err1
	}
	return err
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
		return ssh.Shell(vmName)
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

func CopyFiles(vmName string, file string, toVM bool) error {
	if toVM {
		return copyToVMHomeDir(vmName, file)
	} else {
		return copyToLocalCWD(vmName, file)
	}
}
