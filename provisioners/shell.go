package provisioners

import (
	"github.com/mhewedy/vermin/command/scp"
	"github.com/mhewedy/vermin/command/ssh"
	"path/filepath"
)

type Shell struct{}

func (Shell) Accept(ptype string) bool {
	return "shell" == ptype
}

func (Shell) Exec(vmName string, script string) error {

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
