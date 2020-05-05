package ssh

import (
	"github.com/mhewedy/vermin/command"
	"github.com/mhewedy/vermin/db"
	"github.com/mhewedy/vermin/ip"
)

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
