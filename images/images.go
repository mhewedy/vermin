package images

import (
	"fmt"
	"github.com/mhewedy/vermin/db"
	"os"
)

var (
	format = "%-40s%-10s\n"
	header = fmt.Sprintf(format, "IMAGE NAME", "DISK")
)

func List() ([]string, error) {
	return listCachedImages()
}

func Display() (string, error) {
	list, err := listCachedImages()
	if err != nil {
		return "", err
	}

	result := header

	for i := range list {
		stat, _ := os.Stat(db.GetImageFilePath(list[i]))
		result += fmt.Sprintf(format, list[i],
			fmt.Sprintf("%.1fGB", float64(stat.Size())/(1042*1024*1024.0)))
	}
	return result, nil
}
