package images

import (
	"fmt"
	"github.com/schollz/progressbar/v3"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

func Download(image string) error {
	// check image against cached
	cached, err := listCachedImages()
	if err != nil {
		return err
	}

	if contains(cached, image) {
		return nil
	}

	remote, err := listRemoteImages(false)
	if err != nil {
		return err
	}

	dbImage, err := remote.findByName(image)
	if err != nil {
		return err
	}

	return download(dbImage)
}

func download(r *dbImage) error {
	fmt.Printf("Image '%s' could not be found. Attempting to find and install \n", r.Name)

	// download to a temp file
	tmpFile, err := ioutil.TempFile("", strings.ReplaceAll(r.Name, "/", "_"))
	if err != nil {
		return err
	}
	defer os.Remove(tmpFile.Name())

	resp, err := http.Get(r.URL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	bar := buildDownloadBar(resp)

	if _, err = io.Copy(io.MultiWriter(tmpFile, bar), resp.Body); err != nil {
		return err
	}

	return writeNewImage(tmpFile, r.Name)
}

func buildDownloadBar(resp *http.Response) *progressbar.ProgressBar {
	bar := progressbar.DefaultBytes(
		resp.ContentLength,
		"Downloading",
	)
	progressbar.OptionSetTheme(progressbar.Theme{
		Saucer: "#", SaucerPadding: ".", BarStart: "[", BarEnd: "]"})(bar)
	return bar
}
