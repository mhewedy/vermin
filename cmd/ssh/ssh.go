package ssh

import (
	"vermin/cmd"
	"vermin/db"
	"vermin/ip"
)

func Shell(vmName string) error {
	ipAddr, err := ip.Find(vmName, false)
	if err != nil {
		return err
	}
	return cmd.ExecuteI("ssh", "-i", db.GetPrivateKeyPath(), "-o", "StrictHostKeyChecking=no",
		db.GetUsername()+"@"+ipAddr)
}

func Execute(vmName string, command string) (string, error) {
	ipAddr, err := ip.Find(vmName, false)
	if err != nil {
		return "", err
	}
	return cmd.Execute("ssh", "-i", db.GetPrivateKeyPath(), "-o", "StrictHostKeyChecking=no",
		db.GetUsername()+"@"+ipAddr, "--", command)
}

//ExecuteI execute ssh commands and set cmd stdout to os.Stdout
func ExecuteI(vmName string, command string) error {
	ipAddr, err := ip.Find(vmName, false)
	if err != nil {
		return err
	}
	return cmd.ExecuteI("ssh", "-i", db.GetPrivateKeyPath(), "-o", "StrictHostKeyChecking=no",
		db.GetUsername()+"@"+ipAddr, "--", command)
}
