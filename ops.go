package main

import (
	"errors"
	"fmt"
	"math"
	"time"
	"vermin/cmd"
	"vermin/ssh"
)

type backoff struct {
	iter  int
	start time.Time
	base  time.Duration
	max   time.Duration
}

func (b *backoff) sleep() error {
	elapsed := time.Now().Sub(b.start).Milliseconds()
	if !b.start.IsZero() && elapsed >= b.max.Milliseconds() {
		return errors.New("time elapsed")
	}
	if b.iter == 0 {
		b.start = time.Now()
	}
	time.Sleep(time.Duration(math.Pow(float64(b.iter), 2)) * b.base)
	b.iter++
	return nil
}

func start(vmName string) error {
	fmt.Println("starting", vmName, "...")
	if _, err := cmd.Execute("vboxmanage", "startvm", vmName, "--type", "headless"); err != nil {
		return err
	}
	return nil
}

func stop(vmName string) error {
	fmt.Println("stopping", vmName, "...")
	if _, err := cmd.Execute("vboxmanage", "controlvm", vmName, "poweroff"); err != nil {
		return err
	}
	return nil
}

// establishConn make sure connection to the vm is established or return an error if not
func establishConn(vmName string) error {
	bo := &backoff{
		base: 500 * time.Millisecond,
		max:  5 * time.Minute,
	}
	for {
		if _, err := ssh.Execute(vmName, "ls"); err == nil {
			break
		}
		if err := bo.sleep(); err != nil {
			break
		}
	}
	return nil
}
