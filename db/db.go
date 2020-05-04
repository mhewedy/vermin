package db

import (
	"io/ioutil"
	"log"
	"os"
	"strings"
)

const (
	image = "image"
	tags  = "tags"
)

const (
	NamePrefix      = "vm_"
	ImageFilePrefix = "vermin_images.csv."
)

func GetImagesDir() string {
	return getVerminDir() + string(os.PathSeparator) + "images"
}

func GetVMsBaseDir() string {
	return getVerminDir() + string(os.PathSeparator) + "vms"
}

func GetImageFilePath(imageName string) string {
	return GetImagesDir() + string(os.PathSeparator) + imageName + ".ova"
}

func GetVMPath(vm string) string {
	return GetVMsBaseDir() + string(os.PathSeparator) + vm
}

func GetPrivateKeyPath() string {
	return getVerminDir() + string(os.PathSeparator) + "vermin_rsa"
}

func GetUsername() string {
	return "vermin"
}

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
