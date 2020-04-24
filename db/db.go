package db

import "os"

const (
	Image = "image"
	Tags  = "tags"
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
