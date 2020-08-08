package vms

import (
	"github.com/mhewedy/vermin/command"
	"github.com/mhewedy/vermin/command/scp"
	"github.com/mhewedy/vermin/command/ssh"
	"github.com/mhewedy/vermin/ip"
	"path/filepath"
)

func ProvisionShell(vmName string, script string) error {

	vmFile := "/tmp/" + filepath.Base(script)
	if err := scp.CopyToVM(vmName, script, vmFile); err != nil {
		return err
	}
	if _, err := ssh.Execute(vmName, "chmod +x "+vmFile); err != nil {
		return err
	}
	if err := ssh.ExecInteract(vmName, vmFile); err != nil {
		return err
	}

	return nil
}

func ProvisionAnsible(vmName string, script string) error {

	ipAddr, err := ip.Find(vmName, false)
	if err != nil {
		return err
	}

	return command.AnsiblePlaybook(ipAddr, script).Interact()
}
