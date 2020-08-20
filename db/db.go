package db

import (
	"log"
	"os"
	"path/filepath"
	"strings"
)

const (
	VMNamePrefix       = "vm_"
	ImagesDBFilePrefix = "vermin_images.csv."
)

var (
	ImagesDir  = filepath.Join(getVerminDir(), "images")
	VMsBaseDir = filepath.Join(getVerminDir(), "vms")
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

func GetUsername(vmName string) string {
	vmdb, _ := Load(vmName)
	if IsVagrantImage(vmdb.Image) {
		return "vagrant"
	} else {
		return "vermin"
	}
}

func GetPrivateKeyPath(vmName string) string {
	vmdb, _ := Load(vmName)
	if IsVagrantImage(vmdb.Image) {
		return filepath.Join(getVerminDir(), "vagrant_insecure_private_key")
	} else {
		return filepath.Join(getVerminDir(), "vermin_rsa")
	}
}

// IsValidImage check the image name format to be "vagrant/<base>/<image>[:version]",
// example vagrant/hashicorp/bionic64
func IsVagrantImage(image string) bool {
	s := strings.Split(image, "/")
	return len(s) == 3 && s[0] == "vagrant"
}
