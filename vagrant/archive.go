package vagrant

import (
	"archive/tar"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

func tarFiles(w io.Writer, base string, files []os.FileInfo) error {

	tw := tar.NewWriter(w)

	for _, file := range files {
		hdr := &tar.Header{
			Name: file.Name(),
			Mode: 0600,
			Size: int64(file.Size()),
		}
		if err := tw.WriteHeader(hdr); err != nil {
			return err
		}

		b, err := ioutil.ReadFile(filepath.Join(base, file.Name()))
		if err != nil {
			return err
		}
		if _, err := tw.Write(b); err != nil {
			return err
		}
	}
	if err := tw.Close(); err != nil {
		return err
	}

	return nil
}
