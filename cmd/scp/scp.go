package scp

import (
	"errors"
	"github.com/mhewedy/vermin/cmd"
	"github.com/mhewedy/vermin/db"
	"github.com/mhewedy/vermin/ip"
	"strings"
)

const (
	CopySeparator = ":"
)

type vmPath struct {
	name string
	file string
}

func Copy(src string, dest string) error {

	srcVmPath, srcIsVm := toVmPath(src)
	destVmPath, destIsVm := toVmPath(dest)

	if srcIsVm && destIsVm {
		return copyBetweenVMs(srcVmPath, destVmPath)
	} else if srcIsVm {
		return copyFromVM(srcVmPath, dest)
	} else if destIsVm {
		return copyToVM(src, destVmPath)
	} else {
		return errors.New("src/dest one of them should be vm")
	}
}

// convert <name>:<file> string to vmPath{name, file}
func toVmPath(srcDest string) (vmPath, bool) {
	if !strings.Contains(srcDest, CopySeparator) {
		return vmPath{}, false
	}

	vmAndPath := strings.Split(srcDest, CopySeparator)

	return vmPath{vmAndPath[0], vmAndPath[1]}, true
}

func copyFromVM(vmPath vmPath, localFile string) error {
	ipAddr, err := ip.Find(vmPath.name, false)
	if err != nil {
		return err
	}

	_, err = cmd.Scp(db.GetUsername()+"@"+ipAddr+":"+vmPath.file, localFile).Call()
	return err
}

func copyToVM(localFile string, vmPath vmPath) error {
	ipAddr, err := ip.Find(vmPath.name, false)
	if err != nil {
		return err
	}

	_, err = cmd.Scp(localFile, db.GetUsername()+"@"+ipAddr+":"+vmPath.file).Call()
	return err
}

func copyBetweenVMs(srcVmPath vmPath, destVmPath vmPath) error {

	srcIPAddr, err := ip.Find(srcVmPath.name, false)
	if err != nil {
		return err
	}

	destIPAddr, err := ip.Find(destVmPath.name, false)
	if err != nil {
		return err
	}

	_, err = cmd.Scp("-3",
		db.GetUsername()+"@"+srcIPAddr+":"+srcVmPath.file,
		db.GetUsername()+"@"+destIPAddr+":"+destVmPath.file).Call()
	return err
}
