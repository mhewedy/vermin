package images

import (
	"github.com/mhewedy/vermin/db"
	"io"
	"os"
	"strings"
)

func writeNewImage(tmpFile *os.File, imageName string) error {
	// copy the downloaded file to images directory
	if err := os.MkdirAll(db.ImagesDir+"/"+strings.Split(imageName, "/")[0], 0755); err != nil {
		return err
	}

	if err := copyFile(tmpFile.Name(), db.ImagesDir+"/"+imageName+ova); err != nil {
		return err
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
