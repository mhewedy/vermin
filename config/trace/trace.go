package trace

import (
	"encoding/base64"
	"github.com/mhewedy/go-gistlog"
	"strings"
	"time"
)

var log = gistlog.NewLog("91d041d80687032b270d3c694ea815b8", token())

func PreCreate(imageName string) {

	log.InsertAsync("pre-create", []string{
		imageName,
		time.Now().UTC().String(),
	})
}

func PostCreate(imageName string, err error) {

	_ = log.Insert("post-create", []string{
		imageName,
		time.Now().UTC().String(),
		errorAsString(err),
	})
}

func errorAsString(err error) string {
	var errStr string
	if err != nil {
		errStr = strings.Join(strings.Split(err.Error(), "\n"), " ")
	}
	return errStr
}

func token() string {
	b, err := base64.StdEncoding.DecodeString("OTk5NTIzYmUyNmIyYjk0NGJlYzg5OGZlNjNmYmVhOTM4Y2NiNTE0Mg==")
	if err != nil {
		return ""
	}
	return string(b)
}
