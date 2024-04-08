package images

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/mhewedy/vermin/images/vagrant"
	"github.com/schollz/progressbar/v3"
)

type image struct {
	name string
	url  string
}

func Download(imageName string) error {
	// check image against cached
	cached, err := listCachedImages()
	if err != nil {
		return err
	}

	if contains(cached, imageName) {
		return nil
	}

	fmt.Printf("Image '%s' could not be found. Attempting to find and install \n", imageName)

	url, err := vagrant.GetImageURL("vagrant/" + imageName)
	if err != nil {
		return err
	}
	dbImg := &image{name: imageName, url: url}

	return download(dbImg)
}

func download(r *image) error {

	// download to a temp file
	tmpFile, err := ioutil.TempFile("", strings.ReplaceAll(r.name, "/", "_"))
	if err != nil {
		return err
	}
	defer os.Remove(tmpFile.Name())

	resp, err := http.Get(r.url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	bar := buildDownloadBar(resp)

	if _, err = io.Copy(io.MultiWriter(tmpFile, bar), resp.Body); err != nil {
		return err
	}

	return writeNewImage(tmpFile, r.name)
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
