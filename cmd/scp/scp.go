package scp

import (
	"github.com/mhewedy/vermin/cmd"
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
	_, err = cmd.Execute("scp",
		"-i", db.GetPrivateKeyPath(),
		localFile,
		db.GetUsername()+"@"+ipAddr+":"+vmFile,
	)
	return err
}

func CopyToLocalCWD(vmName string, vmFile string) error {
	ipAddr, err := ip.Find(vmName, false)
	if err != nil {
		return err
	}
	_, err = cmd.Execute("scp",
		"-i", db.GetPrivateKeyPath(),
		db.GetUsername()+"@"+ipAddr+":"+vmFile,
		".",
	)
	return err
}
