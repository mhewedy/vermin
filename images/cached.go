package images

import (
	"github.com/mhewedy/vermin/db"
	"os"
	"path/filepath"
	"strings"
)

const ova = ".ova"

func listCachedImages() ([]string, error) {
	baseDir := db.GetImagesDir() + string(os.PathSeparator)

	images := make([]string, 0)

	err := filepath.Walk(baseDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && strings.HasSuffix(path, ova) {
			name := strings.ReplaceAll(path, baseDir, "")
			name = strings.ReplaceAll(name, ova, "")
			name = strings.ReplaceAll(name, "\\", "/") // for windows
			images = append(images, name)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return images, nil
}
