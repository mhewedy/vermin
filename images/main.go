package images

import (
	"errors"
	"fmt"
	"strings"
	"vermin/cmd"
	"vermin/db"
)

func Create(image string) error {

	// check image against cached
	cached, err := listCachedImages()
	if err != nil {
		return err
	}

	for i := range cached {
		if cached[i] == image {
			return nil // already cached
		}
	}

	remote, err := listRemoteImages()
	if err != nil {
		return err
	}

	// check image against remote
	var vm *vm
	for i := range remote {
		r := remote[i]
		if r.Name == image {
			vm = &r
			break
		}
	}

	if vm == nil {
		return errors.New("invalid image name: " + image)
	}

	return download(vm)
}

func List() (string, error) {

	var result string

	cached, err := listCachedImages()
	if err != nil {
		return "", err
	}

	for i := range cached {
		result += cached[i] + "\t(cached)\n"
	}

	remote, err := listRemoteImagesNames()
	if err != nil {
		return "", err
	}
	for i := range remote {
		r := remote[i]
		if !contains(cached, r) {
			result += r + "\n"
		}
	}

	return result, nil
}

func download(vm *vm) error {
	fmt.Println("Downloading image", vm.Name, "from", vm.URL)
	fmt.Print("It might take a while depending on your internet connection")

	sp := strings.Split(vm.Name, "/")
	vmBasePath := db.GetImagesDir() + sp[0]

	_, err := cmd.Execute("mkdir", "-p", vmBasePath)
	if err != nil {
		return err
	}

	err = cmd.ExecuteAndShowProgress("wget", "-O", vmBasePath+"/"+sp[1]+".ova", vm.URL)
	if err != nil {
		return err
	}
	return nil
}

func contains(a []string, s string) bool {
	for i := range a {
		if a[i] == s {
			return true
		}
	}
	return false
}
