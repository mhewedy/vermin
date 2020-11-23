package ip

import (
	"github.com/mhewedy/vermin/cmd"
	"github.com/mhewedy/vermin/debug"
	"github.com/mhewedy/vermin/hypervisor"
	"github.com/mhewedy/vermin/hypervisor/base"
	"sync"
)

func ping() error {

	subnet, err := hypervisor.GetSubnet()
	if err != nil {
		return err
	}
	debug.Log("subnet: %v", subnet)

	var wg sync.WaitGroup
	wg.Add(subnet.Len)

	for subnet.HasNext() {
		subnet = subnet.Next()

		go func(s *base.Subnet) {
			_ = cmd.Ping(s.IP()).Run()
			wg.Done()
		}(subnet)
	}

	wg.Wait()

	return nil
}
