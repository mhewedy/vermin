package trace

import (
	"encoding/base64"
	"fmt"
	"github.com/matishsiao/goInfo"
	"github.com/mhewedy/go-gistlog"
	"regexp"
	"time"
)

var (
	newline = regexp.MustCompile(`\r?\n`)
	log     = gistlog.NewLog("91d041d80687032b270d3c694ea815b8", token())
)

func PostCreate(imageName, version string, err error) {

	year, month, _ := time.Now().Date()

	_ = log.Insert(fmt.Sprintf("create_%d-%d", year, month), []string{
		imageName,
		time.Now().UTC().String(),
		version,
		goInfo.GetInfo().String(),
		errorAsString(err),
	})
}

func errorAsString(err error) string {
	var errStr string
	if err != nil {
		errStr = newline.ReplaceAllString(err.Error(), " ")
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
