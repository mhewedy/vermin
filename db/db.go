package db

import "os"

const (
	Image = "image"
	Tags  = "tags"
)

func GetVMPath(vm string) string {
	return os.Getenv("HOME") + "/.viper/" + vm + "/"
}
