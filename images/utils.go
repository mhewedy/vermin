package images

import (
	"errors"
	"strings"
)

func contains(a []string, s string) bool {
	for i := range a {
		if a[i] == s {
			return true
		}
	}
	return false
}

func validateName(imageName string) error {
	if len(strings.Split(imageName, "/")) != 2 {
		return errors.New("Name doesn't follow pattern <user>/<box>")
	}
	return nil
}
