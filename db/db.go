package db

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
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
	imagePath := fmt.Sprintf("`%s`", filepath.Join(ImagesDir, imageName+".ova"))
	return imagePath
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
	return filepath.Join(getVerminDir(), VagrantPrivateKey)
}
