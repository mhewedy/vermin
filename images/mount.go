package images

func CanMount(image string) bool {

	remote, _ := listRemoteImages(false)

	for i := range remote {
		r := remote[i]
		if r.Name == image {
			return r.Mount
		}
	}

	return false
}
