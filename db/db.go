package db

import (
	"io/ioutil"
	"log"
	"os"
	"strings"
)

const (
	image    = "image"
	tags     = "tags"
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

// ----
func WriteTag(vmName string, tag string) error {
	return appendToFile(GetVMPath(vmName)+"/"+tags, []byte(tag+"\n"), 0755)
}

func ReadTags(vmName string) (string, error) {
	content, err := readFromFile(vmName, tags)
	if err != nil {
		return "", err
	}
	return strings.TrimSuffix(strings.ReplaceAll(content, "\n", ", "), ", "), nil
}

// ----

func WriteImageData(vmName string, imageName string) error {
	return ioutil.WriteFile(GetVMPath(vmName)+"/"+image, []byte(imageName), 0755)
}

func ReadImageData(vmName string) (string, error) {
	content, err := readFromFile(vmName, image)
	if err != nil {
		return "", err
	}
	return strings.ReplaceAll(content, "\n", ""), nil
}

func getVerminDir() string {
	dir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal("cannot obtain user home dir")
	}
	return dir + string(os.PathSeparator) + ".vermin"
}
