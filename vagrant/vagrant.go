package vagrant

import (
	"github.com/mhewedy/vermin/command"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// IsValidImage check the image name format to be "vagrant/<base>/<image>[:version]",
// example vagrant/hashicorp/bionic64
func IsValidImage(image string) bool {
	s := strings.Split(image, "/")
	return len(s) == 3 && s[0] == "vagrant"
}

func GetImageURL(image string) (string, error) {
	// TODO
	// get json and parse it to return download URL
	return "https://vagrantcloud.com/hashicorp/boxes/bionic64/versions/1.0.282/providers/virtualbox.box", nil
}

func ProcessImage(imagePath string) error {

	imageDir := path.Dir(imagePath)

	// gunzip the downloaded file
	// TODO change from command to golang code
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
