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
	"vermin/port"
)

type delay struct {
	iter  int
	start time.Time
	max   time.Duration
}

func (b *delay) sleep() error {
	elapsed := time.Now().Sub(b.start).Milliseconds()
	if !b.start.IsZero() && elapsed >= b.max.Milliseconds() {
		return errors.New("time elapsed")
	}
	if b.iter == 0 {
		b.start = time.Now()
	}
	b.iter++
	time.Sleep(time.Duration(2*b.iter) * time.Second)
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

func secureShell(vmName string, command string) error {
	if err := establishConn(vmName); err != nil {
		return err
	}

	if len(command) == 0 {
		return ssh.Shell(vmName)
	} else {
		return ssh.ExecuteI(vmName, command)
	}
}

func portForward(vmName string, ports string) error {
	a, err := port.GetPortForwardArgs(ports)
	if err != nil {
		return err
	}

	if err := establishConn(vmName); err != nil {
		return err
	}

	var args = make([]string, len(a)*2+1)
	c := 1
	for i := range args {
		if i%2 == 0 {
			args[i] = "-L"
		} else {
			args[i] = a[i-c]
			c++
		}
	}
	args[len(args)-1] = "-N"

	fmt.Println("Connected. Press CTRL+C anytime to stop")
	if err := ssh.ExecuteIArgs(vmName, args...); err != nil {
		return err
	}

	return nil
}

// establishConn make sure connection to the vm is established or return an error if not
func establishConn(vmName string) error {
	d := &delay{
		max: 1 * time.Minute,
	}
	for {
		if _, err := ssh.Execute(vmName, "ls"); err == nil {
			break
		}
		fmt.Println("Trying to establish connection to", vmName, "...")
		if err := d.sleep(); err != nil {
			break
		}
	}
	return nil
}
