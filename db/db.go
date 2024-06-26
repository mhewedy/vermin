package db

import (
	"log"
	"os"
	"path/filepath"
	"runtime"
)

const (
	VMNamePrefix      = "vm_"
	VagrantPrivateKey = "vagrant_insecure_private_key"
)

var (
	ImagesDir  = filepath.Join(getVerminDir(), "images/vagrant")
	VMsBaseDir = filepath.Join(getVerminDir(), "vms")
	BaseDir    = getVerminDir()
)

func GetImageFilePath(imageName string) string {
	return filepath.Join(ImagesDir, imageName+".ova")
}

func GetVMPath(vm string) string {
	return filepath.Join(VMsBaseDir, vm)
}

func getVerminDir() string {
	dir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal("cannot obtain user home dir")
	}
	return filepath.Join(dir, ".vermin")
}

func GetUsername() string {
	return "vagrant"
}

func GetPrivateKeyPath() string {
	if runtime.GOOS == "windows" {
		return `"` + filepath.Join(getVerminDir(), VagrantPrivateKey) + `"`
	}
	return filepath.Join(getVerminDir(), VagrantPrivateKey)
}
