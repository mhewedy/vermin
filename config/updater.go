package config

import (
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

const versionURL = "https://github.com/mhewedy/vermin/releases/latest"
const updateURL = "https://github.com/mhewedy/vermin#install-vermin"

var client = &http.Client{
	CheckRedirect: func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	},
}

// CheckForUpdates checks for updates at random times
func CheckForUpdates(currentVersion string) {
	rand.Seed(time.Now().UnixNano())
	r := rand.Intn(100)

	if r == 0 {
		req, err := http.NewRequest("HEAD", versionURL, nil)
		if err != nil {
			return
		}

		resp, err := client.Do(req)
		if err != nil {
			return
		}
		defer resp.Body.Close()

		loc, err := resp.Location()
		if err != nil {
			return
		}

		s := strings.Split(loc.Path, "/")
		newVersion := s[len(s)-1]

		if currentVersion != newVersion {
			fmt.Printf("\nNew version avaiable %s, check %s\n\n", newVersion, updateURL)
		}
	}
}
