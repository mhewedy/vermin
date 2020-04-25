package commands

import (
	"path"
	"vermin/cmd"
	"vermin/commands/ip"
	"vermin/db"
)

func copyToVMHomeDir(vmName string, localFile string) error {
	return copyToVM(vmName, localFile, "~/"+path.Base(localFile))
}
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

func copyToLocalCWD(vmName string, vmFile string) error {
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
