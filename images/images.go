package images

import (
	"errors"
	"fmt"
	"github.com/mhewedy/vermin/db"
	"io"
	"io/ioutil"
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
	var rimage *rimage
	for i := range remote {
		r := remote[i]
		if r.Name == image {
			rimage = &r
			break
		}
	}

	if rimage == nil {
		display, _ := List()
		return errors.New(fmt.Sprintf("invalid image name: '%s'.", image) +
			" Valid images are:\n" + strings.Join(display, "\n"))
	}

	return download(rimage)
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

func download(r *rimage) error {
	fmt.Printf("Image '%s' could not be found. Attempting to find and install...\n", r.Name)
	fmt.Printf("Downloading: %s", r.URL)

	// download to a temp file
	tmpFile, err := ioutil.TempFile("", strings.ReplaceAll(r.Name, "/", "_"))
	if err != nil {
		return err
	}
	_ = tmpFile.Close()
	defer os.Remove(tmpFile.Name())

	if _, err := wget(r.URL, tmpFile.Name()); err != nil {
		return err
	}

	// copy the downloaded file to images directory
	if err := os.MkdirAll(db.GetImagesDir()+"/"+strings.Split(r.Name, "/")[0], 0755); err != nil {
		return err
	}
	if err := copyFile(tmpFile.Name(), db.GetImagesDir()+"/"+r.Name+".ova"); err != nil {
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

// Copy the src file to dst. Any existing file will be overwritten and will not
// copy file attributes.
func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}
	return out.Close()
}
