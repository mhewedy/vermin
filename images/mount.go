package images

import (
	"fmt"
	"strings"
)

func CheckCanMount(image string) error {
	remote, _ := listRemoteImages(false)

	dbImage, err := remote.findByName(image)
	if err != nil {
		return err
	}

	if !dbImage.Mount {
		mounted := remote.findByMount(true).names()
		return fmt.Errorf("image '%s' cannot be mounted, "+
			"images can be mounted are:\n%s", image, strings.Join(mounted, "\n"))
	}

	return nil
}
