package ssh

import (
	"errors"
	"github.com/mhewedy/vermin/command"
	"github.com/mhewedy/vermin/db"
	"github.com/mhewedy/vermin/ip"
	"github.com/mhewedy/vermin/progress"
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

	return command.Ssh("-tt", "-i", db.GetPrivateKeyPath(),
		"-o", "StrictHostKeyChecking=no", db.GetUsername()+"@"+ipAddr,
	).Interact()
}

func Execute(vmName string, cmd string) (string, error) {
	ipAddr, err := ip.Find(vmName, false)
	if err != nil {
		return "", err
	}

	return command.Ssh("-i", db.GetPrivateKeyPath(), "-o", "StrictHostKeyChecking=no",
		db.GetUsername()+"@"+ipAddr, "--", cmd).Call()
}

//Interact execute ssh in interactive mode
func ExecuteI(vmName string, cmd string) error {
	ipAddr, err := ip.Find(vmName, false)
	if err != nil {
		return err
	}

	return command.Ssh("-i", db.GetPrivateKeyPath(),
		"-o", "StrictHostKeyChecking=no", db.GetUsername()+"@"+ipAddr, "--", cmd,
	).Interact()
}

// ExecuteIArgs run ssh in interactive mode with args
func ExecuteIArgs(vmName string, args ...string) error {
	ipAddr, err := ip.Find(vmName, false)
	if err != nil {
		return err
	}

	var cargs = []string{"-i", db.GetPrivateKeyPath(), "-o", "StrictHostKeyChecking=no", db.GetUsername() + "@" + ipAddr}
	cargs = append(cargs, args...)

	return command.Ssh(cargs...).Interact()
}

// EstablishConn make sure connection to the vm is established or return an error if not
func EstablishConn(vmName string) error {
	stop := progress.Show("Establishing connection")
	defer stop()

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
