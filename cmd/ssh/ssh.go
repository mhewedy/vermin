package ssh

import (
	"errors"
	"fmt"
	"github.com/mhewedy/vermin/cmd"
	"github.com/mhewedy/vermin/db"
	"github.com/mhewedy/vermin/ip"
	"time"
)

type delay struct {
	iter  int
	start time.Time
	max   time.Duration
}

func (b *delay) sleep(seconds int) error {
	elapsed := time.Now().Sub(b.start).Milliseconds()
	if !b.start.IsZero() && elapsed >= b.max.Milliseconds() {
		return errors.New("cannot accomplish task, time elapsed")
	}
	if b.iter == 0 {
		b.start = time.Now()
	}
	b.iter++
	time.Sleep(time.Duration(seconds) * time.Second)
	return nil
}

func OpenTerminal(vmName string) error {
	ipAddr, err := ip.Find(vmName, false)
	if err != nil {
		return err
	}
	return cmd.ExecuteI("ssh", "-tt", "-i", db.GetPrivateKeyPath(), "-o", "StrictHostKeyChecking=no",
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
	q := cmd.PrintProgress("Establishing connection")
	defer func() { *q <- true; fmt.Println() }()

	d := &delay{
		max: 1 * time.Minute,
	}
	var err error
	for {
		if _, err = Execute(vmName, "ls"); err == nil {
			break
		}
		if err = d.sleep(2); err != nil {
			break
		}
	}
	return err
}
