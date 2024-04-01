package ip

import (
	"sync"

	"github.com/mhewedy/vermin/cmd"
	"github.com/mhewedy/vermin/hypervisor"
	"github.com/mhewedy/vermin/hypervisor/base"
	"github.com/mhewedy/vermin/log"
)

func ping() error {

	subnet, err := hypervisor.GetSubnet()
	if err != nil {
		return err
	}
	log.Debug("subnet: %v", subnet)

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
