package ssh

import (
	"github.com/mhewedy/vermin/cmd"
	"github.com/mhewedy/vermin/ip"
)

func OpenTerminal(vmName string) error {
	ipAddr, err := ip.Find(vmName, false)
	if err != nil {
		return err
	}

	return cmd.Ssh(ipAddr, "-tt").Interact()
}

func Execute(vmName string, command string) (string, error) {
	ipAddr, err := ip.Find(vmName, false)
	if err != nil {
		return "", err
	}

	return cmd.Ssh(ipAddr, "--", command).Call()
}

//ExecInteract execute ssh in interactive mode
func ExecInteract(vmName string, command string) error {
	ipAddr, err := ip.Find(vmName, false)
	if err != nil {
		return err
	}

	return cmd.Ssh(ipAddr, "--", command).Interact()
}

// WithArgs run ssh command with args passed
func WithArgs(vmName string, args []string) error {
	ipAddr, err := ip.Find(vmName, false)
	if err != nil {
		return err
	}

	return cmd.Ssh(ipAddr, args...).Interact()
}
