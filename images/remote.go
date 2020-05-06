package images

import (
	"errors"
	"github.com/artonge/go-csv-tag"
	"github.com/mhewedy/vermin/command"
	"github.com/mhewedy/vermin/db"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

const imagesDBURL = "https://raw.githubusercontent.com/mhewedy/vermin/master/images/images.csv"

type rimage struct {
	Name string `csv:"name"`
	URL  string `csv:"url"`
}

func listRemoteImagesNames(purgeCache bool) ([]string, error) {
	vms, err := listRemoteImages(purgeCache)
	if err != nil {
		return nil, err
	}

	images := make([]string, 0)

	for i := range vms {
		images = append(images, vms[i].Name)
	}

	return images, nil
}

func listRemoteImages(purgeCache bool) ([]rimage, error) {
	if purgeCache {
		file, _ := getCSVTempFilePath()
		if len(file) > 0 {
			_ = os.Remove(file)
		}
	}
	// read images csv from tmp cache
	tmp, _ := getCSVTempFilePath()
	// if not found, then download the file
	if len(tmp) == 0 {
		tmpFile, err := ioutil.TempFile("", db.ImagesDBFilePrefix)
		if err != nil {
			return nil, err
		}
		_ = tmpFile.Close()

		tmp = tmpFile.Name()
		if _, err = command.Wget(imagesDBURL, tmp).Call(); err != nil {
			return nil, err
		}
	}

	// parse the file as csv
	var vms []rimage
	err := csvtag.Load(csvtag.Config{
		Path: tmp,
		Dest: &vms,
	})

	if err != nil {
		return nil, err
	}

	err = validate(vms)
	if err != nil {
		return nil, err
	}

	return vms, nil
}

func getCSVTempFilePath() (string, error) {
	file, err := filepath.Glob(os.TempDir() + "/" + db.ImagesDBFilePrefix + "*")
	if err != nil {
		return "", err
	}
	if len(file) == 0 {
		return "", nil
	}
	return file[0], nil
}

// Check name follows <distro>/<version> name
// Check Unique name
func validate(vms []rimage) error {
	names := make(map[string]bool)

	for i := range vms {
		// check name
		if len(strings.Split(vms[i].Name, "/")) != 2 {
			return errors.New("name doesn't follow pattern <distro>/<version>")
		}
		// check duplicate
		_, found := names[vms[i].Name]
		if found {
			return errors.New("remote list cannot contains duplicates")
		}

		names[vms[i].Name] = true
	}
	return nil
}
