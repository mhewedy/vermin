package provisioners

import (
	"github.com/mhewedy/vermin/command"
	"github.com/mhewedy/vermin/ip"
)

type Ansible struct{}

func (Ansible) Accept(ptype string) bool {
	return ptype == "ansible"
}

func (Ansible) Exec(vmName string, script string) error {

	ipAddr, err := ip.Find(vmName, false)
	if err != nil {
		return err
	}

	return command.AnsiblePlaybook(vmName, ipAddr, script).Interact()
}
