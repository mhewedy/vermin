package vagrant

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/mhewedy/vermin/progress"
	"io/ioutil"
	"net/http"
	"strings"
)

type vagrantResp struct {
	Versions []version `json:"versions"`
}
type provider struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}
type version struct {
	Version   string     `json:"version"`
	Status    string     `json:"status"`
	Providers []provider `json:"providers"`
}

func GetImageURL(image string) (string, error) {

	stop := progress.Show("Getting Image information from Vagrant Cloud", false)
	defer stop()

	user, imageName, imageVersion := getImageParts(image)

	req, err := http.NewRequest("GET", fmt.Sprintf("https://app.vagrantup.com/%s/boxes/%s", user, imageName), nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var obj vagrantResp
	if err := json.Unmarshal(b, &obj); err != nil {
		return "", err
	}

	p, err := filterByVersion(obj, user, imageName, imageVersion)
	if err != nil {
		return "", err
	}

	return p.URL, nil
}

func filterByVersion(resp vagrantResp, user, imageName, imageVersion string) (provider, error) {

	if len(resp.Versions) == 0 {
		return provider{}, fmt.Errorf("Cannot found vagrant image "+
			"at https://app.vagrantup.com/%s/boxes/%s", user, imageName)
	}

	if imageVersion == "" {
		// get latest version
		return getVirtualBoxProvider(resp.Versions[0].Providers)
	} else {
		for _, v := range resp.Versions {
			if v.Version == imageVersion {
				return getVirtualBoxProvider(v.Providers)
			}
		}
		return provider{}, fmt.Errorf("vagrant image version not found, "+
			"check correct version number here: https://app.vagrantup.com/%s/boxes/%s", user, imageName)
	}
}

func getVirtualBoxProvider(providers []provider) (provider, error) {
	for _, p := range providers {
		if p.Name == "virtualbox" {
			return p, nil
		}
	}
	return provider{}, errors.New("VirtualBox image not found for specified version")
}

// getImageParts syntax vagrant/USER/BOX:VERSION e.g.: vagrant/ubuntu/trusty64:20190429.0.1
// returns USER, BOX, VERSION (ubuntu, trusty64, 20190429.0.1)
func getImageParts(image string) (user string, imageName string, imageVersion string) {

	parts := strings.Split(image, "/")
	user = parts[1]
	box := parts[2]

	boxParts := strings.Split(box, ":")
	imageName = boxParts[0]

	if len(boxParts) > 1 {
		imageVersion = boxParts[1]
	}
	return user, imageName, imageVersion
}
