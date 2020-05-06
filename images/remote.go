package images

import (
	"errors"
	"github.com/artonge/go-csv-tag/v2"
	"github.com/mhewedy/vermin/command"
	"github.com/mhewedy/vermin/db"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

const imagesCSVURL = "https://raw.githubusercontent.com/mhewedy/vermin/master/images/images.csv"

func listRemoteImages(purgeCache bool) (dbImages, error) {
	if purgeCache {
		file, _ := getTempFilePath()
		if len(file) > 0 {
			_ = os.Remove(file)
		}
	}
	// read images csv from tmp cache
	tmp, _ := getTempFilePath()
	// if not found, then download the file
	if len(tmp) == 0 {
		tmpFile, err := ioutil.TempFile("", db.ImagesDBFilePrefix)
		if err != nil {
			return nil, err
		}
		_ = tmpFile.Close()

		tmp = tmpFile.Name()
		if _, err = command.Wget(imagesCSVURL, tmp).Call(); err != nil {
			return nil, err
		}
	}

	// parse the file as csv
	var dbImages []dbImage

	if err := csvtag.LoadFromPath(tmp, &dbImages); err != nil {
		return nil, err
	}

	if err := validate(dbImages); err != nil {
		return nil, err
	}

	return dbImages, nil
}

func getTempFilePath() (string, error) {
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
func validate(vms []dbImage) error {
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
