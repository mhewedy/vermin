package db

import (
	"log"
	"os"
)

const (
	Username = "vermin"
)

const (
	VMNamePrefix       = "vm_"
	ImagesDBFilePrefix = "vermin_images.csv."
)

var (
	ImagesDir      = getVerminDir() + string(os.PathSeparator) + "images"
	VMsBaseDir     = getVerminDir() + string(os.PathSeparator) + "vms"
	PrivateKeyPath = getVerminDir() + string(os.PathSeparator) + "vermin_rsa"
)

func GetImageFilePath(imageName string) string {
	return ImagesDir + string(os.PathSeparator) + imageName + ".ova"
}

func GetVMPath(vm string) string {
	return VMsBaseDir + string(os.PathSeparator) + vm
}

func getVerminDir() string {
	dir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal("cannot obtain user home dir")
	}
	return dir + string(os.PathSeparator) + ".vermin"
}
