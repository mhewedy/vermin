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
	return GetVMsBaseDir() + "/" + vm
}

func GetHomeDir() string {
	dir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal("cannot obtain user home dir")
	}
	return dir + "/.vermin"
}

func GetImagesDir() string {
	return GetHomeDir() + "/images"
}

func GetImageFilePath(imageName string) string {
	return GetImagesDir() + "/" + imageName + ".ova"
}

func GetVMsBaseDir() string {
	return GetHomeDir() + "/vms"
}

func GetPrivateKeyPath() string {
	return GetHomeDir() + "/vermin_rsa"
}

func GetUsername() string {
	return "vermin"
}
