package vagrant

import (
	"errors"
	"github.com/mhewedy/vermin/log"
	"github.com/mhewedy/vermin/progress"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
)

var (
	// errImageIsNotGzipped return when unable to gunzip an image,
	// mostly it is only tarred and could be used directly and no need to gunzip and tar
	errImageIsNotGzipped = errors.New("image is not gzipped, may be tarred only")
)

// ProcessImage process vagrant image after downloaded
func ProcessImage(imagePath string) error {

	stop := progress.Show("Convert Vagrant image into Vermin format", false)
	defer stop()

	imageDir := path.Dir(imagePath)

	if err := gunzipVagrantBox(imagePath, imageDir, false); err != nil {
		if err != errImageIsNotGzipped {
			return err
		}

		log.Info("cannot ungzip image %s, try untaring only...", imagePath)
		if err = gunzipVagrantBox(imagePath, imageDir, true); err != nil {
			return err
		}
	}

	// remove the downloaded file
	if err := os.Remove(imagePath); err != nil {
		return err
	}

	if err := createOVAFile(imagePath, imageDir); err != nil {
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

func gunzipVagrantBox(imagePath, imageDir string, tarOnly bool) error {

	gzipFile, err := os.Open(imagePath)
	if err != nil {
		return err
	}
	defer gzipFile.Close()

	return gunzip(gzipFile, imageDir, tarOnly)
}

func createOVAFile(imagePath, imageDir string) error {
	// get ovf, vmdk FileInfo
	infos, err := ioutil.ReadDir(imageDir)
	if err != nil {
		return err
	}
	var ovaFileInfo os.FileInfo
	vmdkFileInfos := make([]os.FileInfo, 0)
	for _, info := range infos {
		if strings.HasSuffix(info.Name(), ".ovf") {
			ovaFileInfo = info
		}
		if strings.HasSuffix(info.Name(), ".vmdk") {
			vmdkFileInfos = append(vmdkFileInfos, info)
		}
	}

	// create ova by TARing (ovf and vmdk)
	file, err := os.Create(imagePath)
	if err != nil {
		return err
	}
	defer file.Close()

	return tarFiles(file, imageDir, append(vmdkFileInfos, ovaFileInfo))
}
