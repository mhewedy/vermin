package virtualbox

import (
	"strings"
)

func findBridgeInfo(keys ...string) ([]string, error) {

	r, err := vboxManage("list", "bridgedifs").Call()
	if err != nil {
		return nil, err
	}

	list := strings.Split(r, "\n\n")[0]
	lines := strings.Split(list, "\n")

	elements := make(map[string]string)

	for _, line := range lines {
		s := strings.SplitN(line, ":", 2)

		key := strings.TrimSpace(s[0])
		value := strings.TrimSpace(s[1])
		elements[key] = value
	}

	var ret = make([]string, len(keys))
	for i, key := range keys {
		ret[i] = elements[key]
	}

	return ret, nil
}
