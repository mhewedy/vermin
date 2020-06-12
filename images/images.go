package images

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

	var result string
	for i := range list {
		if list[i].cached {
			result += list[i].name + "\t\t(cached)\n"
		} else {
			result += list[i].name + "\n"
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
