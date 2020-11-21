package ssh

import (
	"errors"
	"github.com/mhewedy/vermin/progress"
	"time"
)

type delay struct {
	errMsg string
	iter   int
	start  time.Time
	max    time.Duration
}

func (b *delay) sleep() error {
	elapsed := time.Now().Sub(b.start).Milliseconds()
	if !b.start.IsZero() && elapsed >= b.max.Milliseconds() {
		return errors.New(b.errMsg)
	}
	if b.iter == 0 {
		b.start = time.Now()
	}
	b.iter++
	time.Sleep(100 * time.Millisecond)
	return nil
}

// EstablishConn make sure connection to the vm is established or return an error if not
func EstablishConn(vmName string) error {
	stop := progress.Show("Establishing connection", true)
	defer stop()

	d := &delay{
		errMsg: "Cannot establish connection.",
		max:    1 * time.Minute,
	}
	var err error
	for {
		if _, err = Execute(vmName, "ls"); err == nil {
			break
		}
		if err = d.sleep(); err != nil {
			break
		}
	}
	return err
}
