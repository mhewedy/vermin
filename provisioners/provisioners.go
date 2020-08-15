package provisioners

import "fmt"

type Func func(vmName string, script string) error

type Provisioner interface {
	Accept(ptype string) bool
	Exec(vmName string, script string) error
}

func Load(pType string) (Func, error) {

	// TODO in future if we need to implement plugin arch,
	//  we might construct this slice dynamically by checking some directory of binaries of some name
	// e.g. vermin-provisioner-<name>
	ps := []Provisioner{
		Shell{},
		Ansible{},
	}

	for _, p := range ps {
		if p.Accept(pType) {
			return p.Exec, nil
		}
	}

	return nil, fmt.Errorf("Cannot load provisioner of type: %s", pType)
}
