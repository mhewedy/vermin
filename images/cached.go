package images

import (
	"errors"
	"github.com/mhewedy/vermin/db"
	"github.com/mhewedy/vermin/progress"
	"os"
	"path/filepath"
	"strings"
)

const ova = ".ova"

func Remove(imageName string) error {
	existingImgs, _ := List()
	if !contains(existingImgs, imageName) {
		return errors.New("Image not found")
	}

	if err := os.RemoveAll(db.GetImageFilePath(imageName)); err != nil {
		return err
	}

	progress.Immediate("Deleting image", imageName)
	return nil
}

func listCachedImages() ([]string, error) {
	baseDir := db.ImagesDir + string(os.PathSeparator)

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
