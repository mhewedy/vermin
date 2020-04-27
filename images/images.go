package images

import (
	"errors"
	"fmt"
	"github.com/mhewedy/vermin/cmd"
	"github.com/mhewedy/vermin/db"
	"os"
	"strings"
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
		display, _ := List()
		return errors.New(fmt.Sprintf("invalid image name: '%s'.", image) +
			" Valid images are:\n" + strings.Join(display, "\n"))
	}

	return download(vm)
}

type image struct {
	name   string
	cached bool
}

func List() ([]string, error) {
	list, err := list()
	if err != nil {
		return nil, err
	}
	var result = make([]string, len(list))
	for i := range list {
		result[i] = list[i].name
	}
	return result, nil
}

func Display() (string, error) {

	list, err := list()
	if err != nil {
		return "", err
	}
	var result string

	for i := range list {
		if list[i].cached {
			result += list[i].name + "\t\t(cached)\n"
		} else {
			result += list[i].name + "\n"
		}
	}
	return result, nil
}

func list() ([]image, error) {
	var result []image

	cached, err := listCachedImages()
	if err != nil {
		return nil, err
	}
	for i := range cached {
		result = append(result, image{cached[i], true})
	}

	remote, err := listRemoteImagesNames()
	if err != nil {
		return nil, err
	}
	for i := range remote {
		r := remote[i]
		if !contains(cached, r) {
			result = append(result, image{r, false})
		}
	}
	return result, nil
}

func download(vm *vm) error {
	fmt.Printf("Downloading image %s ", vm.Name)

	sp := strings.Split(vm.Name, "/")
	vmBasePath := db.GetImagesDir() + "/" + sp[0]

	if err := os.MkdirAll(vmBasePath, 0755); err != nil {
		return err
	}

	if _, err := cmd.ExecuteP("wget", "-O", vmBasePath+"/"+sp[1]+".ova", vm.URL); err != nil {
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
