package vagrant

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/mhewedy/vermin/hypervisor"
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

type suggestionResp struct {
	Boxes []struct {
		Tag string `json:"tag"`
	} `json:"boxes"`
}

func GetImageURL(image string) (string, error) {

	stop := progress.Show("Getting Image info from Vagrant Cloud", false)
	defer stop()

	user, imageName, imageVersion := getImageParts(image)

	bytes, err := sendAndReceive(fmt.Sprintf("https://app.vagrantup.com/%s/boxes/%s", user, imageName))
	if err != nil {
		return "", nil
	}

	var obj vagrantResp
	if err := json.Unmarshal(bytes, &obj); err != nil {
		return "", err
	}

	p, err := filterByVersion(obj, user, imageName, imageVersion)
	if err != nil {
		suggest, _ := getSuggestion(imageName, "* ")
		if suggest != nil && len(suggest) > 0 {
			return "", errors.New(err.Error() + ", do you mean:\n" + strings.Join(suggest, "\n"))
		}
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
		return getCurrentProvider(resp.Versions[0].Providers)
	} else {
		for _, v := range resp.Versions {
			if v.Version == imageVersion {
				return getCurrentProvider(v.Providers)
			}
		}
		return provider{}, fmt.Errorf("vagrant image version not found, "+
			"check correct version number here: https://app.vagrantup.com/%s/boxes/%s", user, imageName)
	}
}

func getCurrentProvider(providers []provider) (provider, error) {

	h, err := hypervisor.GetHypervisorName(true)
	if err != nil {
		return provider{}, err
	}

	for _, p := range providers {
		if p.Name == h {
			return p, nil
		}
	}
	return provider{}, fmt.Errorf("%s: image not found for specified version", h)
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

func getSuggestion(imageName, prepend string) ([]string, error) {

	h, err := hypervisor.GetHypervisorName(false)
	if err != nil {
		return nil, err
	}

	bytes, err := sendAndReceive(fmt.Sprintf("https://app.vagrantup.com/api/v1/search?sort=downloads&provider=%s&q=%s&limit=2",
		h, imageName))
	if err != nil {
		return nil, nil
	}

	var obj suggestionResp
	if err := json.Unmarshal(bytes, &obj); err != nil {
		return nil, err
	}

	size := len(obj.Boxes)
	if size == 0 {
		return []string{}, nil
	}

	result := make([]string, size)
	for i, b := range obj.Boxes {
		result[i] = prepend + b.Tag
	}
	return result, nil
}

func sendAndReceive(url string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return b, nil
}
