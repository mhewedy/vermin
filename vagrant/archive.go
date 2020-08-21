package vagrant

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
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

func gunzip(baseDir string, gzipStream io.Reader) error {

	uncompressedStream, err := gzip.NewReader(gzipStream)
	if err != nil {
		return errImageIsNotGzipped
	}

	tarReader := tar.NewReader(uncompressedStream)

	for true {
		header, err := tarReader.Next()

		if err == io.EOF {
			break
		}

		if err != nil {
			return err
		}

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.Mkdir(filepath.Join(baseDir, header.Name), 0755); err != nil {
				return err
			}
		case tar.TypeReg:
			outFile, err := os.Create(filepath.Join(baseDir, header.Name))
			if err != nil {
				return err
			}
			defer outFile.Close()
			if _, err := io.Copy(outFile, tarReader); err != nil {
				return err
			}

		default:
			return fmt.Errorf("ExtractTarGz: uknown type: %s in %s", header.Typeflag, header.Name)
		}

	}

	return nil
}
