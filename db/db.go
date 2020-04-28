package db

import (
	"log"
	"os"
)

const (
	Image      = "image"
	Tags       = "tags"
	NamePrefix = "vm_"
	ImageFile  = "vermin_images.csv."
)

func GetVMPath(vm string) string {
	return GetVMsBaseDir() + string(os.PathSeparator) + vm
}

func GetHomeDir() string {
	dir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal("cannot obtain user home dir")
	}
	return dir + string(os.PathSeparator) + ".vermin"
}

func GetImagesDir() string {
	return GetHomeDir() + string(os.PathSeparator) + "images"
}

func GetImageFilePath(imageName string) string {
	return GetImagesDir() + string(os.PathSeparator) + imageName + ".ova"
}

func GetVMsBaseDir() string {
	return GetHomeDir() + string(os.PathSeparator) + "vms"
}

func GetPrivateKeyPath() string {
	return GetHomeDir() + string(os.PathSeparator) + "vermin_rsa"
}

func GetUsername() string {
	return "vermin"
}
