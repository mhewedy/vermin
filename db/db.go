package db

import "os"

const (
	Image = "image"
	Tags  = "tags"
)

func GetVMPath(vm string) string {
	return GetHomeDir() + "/" + vm + "/"
}

func GetHomeDir() string {
	return os.Getenv("HOME") + "/.viper"
}

func GetImagesDir() string {
	return GetHomeDir() + "/boxes/"
}
