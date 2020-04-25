package main

import (
	"errors"
	"fmt"
	"os"
	"time"
	"vermin/cmd"
	"vermin/cmd/ssh"
	"vermin/db"
	"vermin/info"
)

type delay struct {
	iter  int
	start time.Time
	max   time.Duration
}

func (b *delay) sleep(seconds int) error {
	elapsed := time.Now().Sub(b.start).Milliseconds()
	if !b.start.IsZero() && elapsed >= b.max.Milliseconds() {
		return errors.New("time elapsed")
	}
	if b.iter == 0 {
		b.start = time.Now()
	}
	time.Sleep(time.Duration(seconds) * time.Second)
	b.iter++
	return nil
}

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

// establishConn make sure connection to the vm is established or return an error if not
func establishConn(vmName string) error {
	d := &delay{
		max: 3 * time.Minute,
	}
	for {
		if _, err := ssh.Execute(vmName, "ls"); err == nil {
			break
		}
		if err := d.sleep(3); err != nil {
			break
		}
	}
	return nil
}
