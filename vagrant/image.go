package vagrant

import (
	"github.com/mhewedy/vermin/command"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func ProcessImage(imagePath string) error {

	imageDir := path.Dir(imagePath)

	// gunzip the downloaded file
	// TODO change from using tar command to golang code
	if err := command.Tar("xzf", imagePath, "-C", imageDir).Run(); err != nil {
		return err
	}

	// remove the downloaded file
	if err := os.Remove(imagePath); err != nil {
		return err
	}

	// get ovf, vmdk FileInfo
	infos, err := ioutil.ReadDir(imageDir)
	if err != nil {
		return err
	}
	var ovaFileInfo, vmdkFileInfo os.FileInfo
	for _, info := range infos {
		if strings.HasSuffix(info.Name(), ".ovf") {
			ovaFileInfo = info
		}
		if strings.HasSuffix(info.Name(), ".vmdk") {
			vmdkFileInfo = info
		}
	}

	// create ova by TARing (ovf and vmdk)
	file, err := os.Create(imagePath)
	if err != nil {
		return err
	}
	defer file.Close()

	if err := tarFiles(file, imageDir, []os.FileInfo{ovaFileInfo, vmdkFileInfo}); err != nil {
		return err
	}

	// remove all files except ova
	return filepath.Walk(imageDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && !strings.HasSuffix(info.Name(), ".ova") {
			if err := os.Remove(path); err != nil {
				return err
			}
		}
		return nil
	})
}
