package provisioners

import (
	"fmt"
	"github.com/mhewedy/vermin/cmd/scp"
	"github.com/mhewedy/vermin/cmd/ssh"
	"path/filepath"
)

type Shell struct{}

func (Shell) Accept(ptype string) bool {
	return "shell" == ptype
}

func (Shell) Exec(vmName string, script string) error {

	vmFile := "/tmp/" + filepath.Base(script)
	toVM := fmt.Sprintf("%s%s%s", vmName, scp.CopySeparator, vmFile)

	if err := scp.Copy(script, toVM); err != nil {
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
