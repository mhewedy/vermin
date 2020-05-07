package ssh

import (
	"errors"
	"github.com/mhewedy/vermin/progress"
	"time"
)

type delay struct {
	iter  int
	start time.Time
	max   time.Duration
}

func (b *delay) sleep(seconds int) error {
	elapsed := time.Now().Sub(b.start).Milliseconds()
	if !b.start.IsZero() && elapsed >= b.max.Milliseconds() {
		return errors.New("Cannot accomplish task, time elapsed")
	}
	if b.iter == 0 {
		b.start = time.Now()
	}
	b.iter++
	time.Sleep(time.Duration(seconds) * time.Second)
	return nil
}

// EstablishConn make sure connection to the vm is established or return an error if not
func EstablishConn(vmName string) error {
	stop := progress.Show("Establishing connection")
	defer stop()

	d := &delay{
		max: 1 * time.Minute,
	}
	var err error
	for {
		if _, err = Execute(vmName, "ls"); err == nil {
			break
		}
		if err = d.sleep(2); err != nil {
			break
		}
	}
	return err
}
