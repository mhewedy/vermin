package ssh

import (
	"github.com/mhewedy/vermin/command"
	"github.com/mhewedy/vermin/ip"
)

func OpenTerminal(vmName string) error {
	ipAddr, err := ip.Find(vmName, false)
	if err != nil {
		return err
	}

	return command.Ssh(ipAddr, "-tt").Interact()
}

func Execute(vmName string, cmd string) (string, error) {
	ipAddr, err := ip.Find(vmName, false)
	if err != nil {
		return "", err
	}

	return command.Ssh(ipAddr, "--", cmd).Call()
}

//ExecInteract execute ssh in interactive mode
func ExecInteract(vmName string, cmd string) error {
	ipAddr, err := ip.Find(vmName, false)
	if err != nil {
		return err
	}

	return command.Ssh(ipAddr, "--", cmd).Interact()
}

// WithArgs run ssh command with args passed
func WithArgs(vmName string, args []string) error {
	ipAddr, err := ip.Find(vmName, false)
	if err != nil {
		return err
	}

	return command.Ssh(ipAddr, args...).Interact()
}
