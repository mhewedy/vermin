package main

import (
	"vermin/cmd"
	"vermin/db"
	"vermin/ip"
)

func copyToVM(vmName string, localFile string, vmFile string) error {
	ipAddr, err := ip.Find(vmName, false)
	if err != nil {
		return err
	}
	_, err = cmd.Execute("scp",
		"-i", db.GetPrivateKeyPath(),
		localFile,
		db.GetUsername()+"@"+ipAddr+":"+vmFile,
	)
	return err
}
