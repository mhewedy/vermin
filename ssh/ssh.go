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
	return cmd.ExecuteI("ssh", "-i", db.GetPrivateKeyPath(), db.GetUsername()+"@"+ipAddr)
}

func Execute(vmName string, command string) (string, error) {
	ipAddr, err := ip.Find(vmName, false)
	if err != nil {
		return "", err
	}
	return cmd.Execute("ssh", "-i", db.GetPrivateKeyPath(), db.GetUsername()+"@"+ipAddr, "--", command)
}

//ExecuteO execute ssh commands and set cmd stdout to os.Stdout
func ExecuteO(vmName string, command string) error {
	ipAddr, err := ip.Find(vmName, false)
	if err != nil {
		return err
	}
	return cmd.ExecuteO("ssh", "-i", db.GetPrivateKeyPath(), db.GetUsername()+"@"+ipAddr, "--", command)
}
