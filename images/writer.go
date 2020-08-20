package images

import (
	"github.com/mhewedy/vermin/db"
	"github.com/mhewedy/vermin/vagrant"
	"io"
	"os"
	"strings"
)

func writeNewImage(tmpFile *os.File, imageName string) error {
	// copy the downloaded file to images directory
	parts := strings.Split(imageName, "/")
	if err := os.MkdirAll(db.ImagesDir+"/"+strings.Join(parts[0:len(parts)-1], "/"), 0755); err != nil {
		return err
	}

	imagePath := db.ImagesDir + "/" + imageName + ova
	if err := copyFile(tmpFile.Name(), imagePath); err != nil {
		return err
	}

	if vagrant.IsValidImage(imageName) {
		if err := vagrant.ProcessImage(imagePath); err != nil {
			return err
		}
	}

	return nil
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
