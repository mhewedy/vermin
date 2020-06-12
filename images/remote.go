package images

import (
	"errors"
	"github.com/artonge/go-csv-tag/v2"
	"github.com/mhewedy/vermin/db"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
)

const imagesCSVURL = "https://raw.githubusercontent.com/mhewedy/vermin/master/images/images.csv"

func listRemoteImages(purgeCache bool) (dbImages, error) {
	if purgeCache {
		file, _ := getCSVTempFilePath()
		if len(file) > 0 {
			_ = os.Remove(file)
		}
	}
	// read images csv from csvFile cache
	csvFile, _ := getCSVTempFilePath()

	// if not found, then download the file
	if len(csvFile) == 0 {
		tmpFile, err := ioutil.TempFile("", db.ImagesDBFilePrefix)
		if err != nil {
			return nil, err
		}
		defer tmpFile.Close()

		csvFile = tmpFile.Name()

		resp, err := http.Get(imagesCSVURL)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		if _, err := io.Copy(tmpFile, resp.Body); err != nil {
			return nil, err
		}
	}

	// parse the file as csv
	var dbImages []dbImage

	if err := csvtag.LoadFromPath(csvFile, &dbImages); err != nil {
		return nil, err
	}

	if err := validate(dbImages); err != nil {
		return nil, err
	}

	return dbImages, nil
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
func validate(vms []dbImage) error {
	names := make(map[string]bool)

	for i := range vms {
		// check name

		if err := validateName(vms[i].Name); err != nil {
			return err
		}
		// check duplicate
		_, found := names[vms[i].Name]
		if found {
			return errors.New("Remote list cannot contains duplicates")
		}

		names[vms[i].Name] = true
	}
	return nil
}
