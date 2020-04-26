package db

import "os"

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
	return os.Getenv("HOME") + "/.vermin"
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
