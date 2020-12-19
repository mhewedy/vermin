package vms

import (
	"fmt"
	"github.com/mhewedy/vermin/cmd/ssh"
	"github.com/mhewedy/vermin/hypervisor"
	"github.com/mhewedy/vermin/progress"
)

func Shrink(vmName string) error {

	fmt.Println("The VM will restarted as part of the disk shrinking process.\n" +
		"Please note that, this is a time-consuming process and requires a free disk space on your disk.")

	if err := ssh.EstablishConn(vmName); err != nil {
		return err
	}

	if err := zerofyDisk(vmName); err != nil {
		return err
	}

	if err := Stop(vmName); err != nil {
		return err
	}

	if err := hypervisor.ShrinkDisk(vmName); err != nil {
		return err
	}

	return Start(vmName)
}

func zerofyDisk(vmName string) error {
	stop := progress.Show("Filling free disk space with zeros", false)
	defer stop()
	// sometimes the an error returned, however the command succeed
	_, _ = ssh.Execute(vmName, "sh -c 'cat /dev/zero > zero.fill; sync; sleep 1; sync; rm -f zero.fill'")
	return nil
}
