package provisioners

import (
	"github.com/mhewedy/vermin/cmd"
	"github.com/mhewedy/vermin/ip"
)

type Ansible struct{}

func (Ansible) Accept(ptype string) bool {
	return ptype == "ansible"
}

func (Ansible) Exec(vmName string, script string) error {

	ipAddr, err := ip.GetIpAddress(vmName)
	if err != nil {
		return err
	}

	return cmd.AnsiblePlaybook(ipAddr, script).Interact()
}
