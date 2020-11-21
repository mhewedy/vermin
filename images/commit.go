package images

import (
	"errors"
	"fmt"
	"github.com/mhewedy/vermin/hypervisor"
	"io/ioutil"
	"os"
	"strings"
)

func Commit(vmName, imageName string, override bool) error {

	if err := validateName(imageName); err != nil {
		return err
	}

	if !override {
		existingImgs, _ := List()
		if contains(existingImgs, imageName) {
			return errors.New("Image with same name already exists, either choose a new name or use the --override flag")
		}
	}

	tmpDir, err := ioutil.TempDir("", "vermin_commit")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmpDir)

	ovaFile := tmpDir + strings.ReplaceAll(imageName, "/", "_") + ova

	if err = hypervisor.Commit(vmName, ovaFile); err != nil {
		return err
	}

	tmpFile, err := os.Open(ovaFile)
	if err != nil {
		return err
	}
	defer tmpFile.Close()

	if err := writeNewImage(tmpFile, imageName); err != nil {
		return err
	}

	fmt.Printf("\nImage is ready, to create a VM from it use:\n$ vermin create %s\n", imageName)
	return nil
}
