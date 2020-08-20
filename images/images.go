package images

import (
	"fmt"
	"github.com/mhewedy/vermin/db"
	"os"
)

var (
	format = "%-40s%-20s%-10s\n"
	header = fmt.Sprintf(format, "IMAGE NAME", "CACHED", "DISK")
)

func List() ([]string, error) {
	list, err := list(false)
	if err != nil {
		return nil, err
	}

	var result = make([]string, len(list))
	for i := range list {
		result[i] = list[i].name
	}
	return result, nil
}

func Display(purgeCache bool) (string, error) {
	list, err := list(purgeCache)
	if err != nil {
		return "", err
	}

	result := header

	for i := range list {
		if list[i].cached {
			stat, _ := os.Stat(db.GetImageFilePath(list[i].name))
			result += fmt.Sprintf(format, list[i].name, "true",
				fmt.Sprintf("%.1fGB", float64(stat.Size())/(1042*1024*1024.0)))
		} else {
			result += fmt.Sprintf(format, list[i].name, "", "")
		}
	}
	return result, nil
}

type image struct {
	name   string
	cached bool
}

func list(purgeCache bool) ([]image, error) {
	var result []image

	cached, err := listCachedImages()
	if err != nil {
		return nil, err
	}
	for i := range cached {
		result = append(result, image{cached[i], true})
	}

	remote, err := listRemoteImages(purgeCache)
	if err != nil {
		return nil, err
	}

	for i := range remote {
		r := remote[i]
		if !contains(cached, r.Name) {
			result = append(result, image{r.Name, false})
		}
	}
	return result, nil
}
