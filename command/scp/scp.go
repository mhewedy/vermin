package scp

import (
	"github.com/mhewedy/vermin/command"
	"github.com/mhewedy/vermin/db"
	"github.com/mhewedy/vermin/ip"
	"path"
)

func CopyToVMHomeDir(vmName string, localFile string) error {
	return CopyToVM(vmName, localFile, "~/"+path.Base(localFile))
}

func CopyToVM(vmName string, localFile string, vmFile string) error {
	ipAddr, err := ip.Find(vmName, false)
	if err != nil {
		return err
	}

	_, err = command.Scp("-i", db.PrivateKeyPath, localFile, db.Username+"@"+ipAddr+":"+vmFile).Call()

	return err
}

func CopyToLocalCWD(vmName string, vmFile string) error {
	ipAddr, err := ip.Find(vmName, false)
	if err != nil {
		return err
	}

	_, err = command.Scp("-i", db.PrivateKeyPath, db.Username+"@"+ipAddr+":"+vmFile, "./").Call()

	return err
}
