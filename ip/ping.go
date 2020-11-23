package ip

import (
	"github.com/mhewedy/vermin/cmd"
	"github.com/mhewedy/vermin/debug"
	"sync"
)

func ping() {

	cidrs := getCIDRs()
	debug.Log("ciders: %s", cidrs)

	var wg sync.WaitGroup
	wg.Add(len(cidrs))

	for _, c := range cidrs {

		if c.len > 65535 { // skipping cidrs with len > x.x.x.x/16 including 127.0.0.0/8
			wg.Done()
			continue
		}

		debug.Log("pining cider %s", c)
		go func(c cidr) {
			pingCIDR(c)
			wg.Done()
		}(c)
	}

	wg.Wait()
}

func pingCIDR(c cidr) {

	var wg sync.WaitGroup
	wg.Add(c.len)

	for c.hasNext() {
		c = c.next()

		go func(cc cidr) {
			_ = cmd.Ping(cc.IP()).Run()
			wg.Done()
		}(c)
	}

	wg.Wait()

}
