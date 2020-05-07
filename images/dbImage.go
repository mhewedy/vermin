package images

import (
	"fmt"
	"strings"
)

type dbImage struct {
	Name  string `csv:"name"`
	URL   string `csv:"url"`
	Mount bool   `csv:"mount"`
}

type dbImages []dbImage

func (dbImages dbImages) findByName(name string) (*dbImage, error) {
	var dbImage *dbImage

	for i := range dbImages {
		r := dbImages[i]
		if r.Name == name {
			dbImage = &r
			break
		}
	}

	if dbImage == nil {
		return nil, fmt.Errorf("Invalid image name: '%s', valid images are:\n%s", name, strings.Join(dbImages.names(), "\n"))
	}

	return dbImage, nil
}

func (dbImages dbImages) findByMount(mount bool) dbImages {
	var ret []dbImage

	for i := range dbImages {
		r := dbImages[i]
		if r.Mount == mount {
			ret = append(ret, r)
		}
	}

	return ret
}

func (dbImages dbImages) names() []string {
	var names []string
	for i := range dbImages {
		r := dbImages[i]
		names = append(names, r.Name)
	}
	return names
}
