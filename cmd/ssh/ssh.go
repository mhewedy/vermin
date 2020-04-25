package ssh

import (
	"errors"
	"fmt"
	"time"
	"vermin/cmd"
	"vermin/db"
	"vermin/ip"
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

func Shell(vmName string) error {
	ipAddr, err := ip.Find(vmName, false)
	if err != nil {
		return err
	}
	return cmd.ExecuteI("ssh", "-i", db.GetPrivateKeyPath(), "-o", "StrictHostKeyChecking=no",
		db.GetUsername()+"@"+ipAddr)
}

func Execute(vmName string, command string) (string, error) {
	ipAddr, err := ip.Find(vmName, false)
	if err != nil {
		return "", err
	}
	return cmd.Execute("ssh", "-i", db.GetPrivateKeyPath(), "-o", "StrictHostKeyChecking=no",
		db.GetUsername()+"@"+ipAddr, "--", command)
}

//ExecuteI execute ssh in interactive mode
func ExecuteI(vmName string, command string) error {
	ipAddr, err := ip.Find(vmName, false)
	if err != nil {
		return err
	}
	return cmd.ExecuteI("ssh", "-i", db.GetPrivateKeyPath(), "-o", "StrictHostKeyChecking=no",
		db.GetUsername()+"@"+ipAddr, "--", command)
}

// ExecuteIArgs run ssh in interactive mode with args
func ExecuteIArgs(vmName string, args ...string) error {
	ipAddr, err := ip.Find(vmName, false)
	if err != nil {
		return err
	}

	var cargs = []string{"-i", db.GetPrivateKeyPath(), "-o", "StrictHostKeyChecking=no", db.GetUsername() + "@" + ipAddr}
	cargs = append(cargs, args...)

	return cmd.ExecuteI("ssh", cargs...)
}

// EstablishConn make sure connection to the vm is established or return an error if not
func EstablishConn(vmName string) error {
	d := &delay{
		max: 1 * time.Minute,
	}
	for {
		if _, err := Execute(vmName, "ls"); err == nil {
			break
		}
		fmt.Println("Trying to establish connection to", vmName, "...")
		if err := d.sleep(); err != nil {
			break
		}
	}
	return nil
}
