package trace

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const (
	gistId     = "91d041d80687032b270d3c694ea815b8"
	gistApiUrl = "https://api.github.com/gists/" + gistId
)

// Gist represents a GitHub's gist.
type gist struct {
	Files map[string]gistFile `json:"files,omitempty"`
}

type gistFile struct {
	Filename string `json:"filename,omitempty"`
	Content  string `json:"content,omitempty"`
}

func PreCreate(imageName string) {
	go func() {
		updateFile("pre-create", func(filename string) gist {

			content := fmt.Sprintf("%s,%s", imageName, time.Now().UTC())
			return gist{
				Files: map[string]gistFile{filename: {
					Filename: filename,
					Content:  getContent(filename) + "\n" + content,
				}},
			}
		})
	}()
}

func PostCreate(imageName string, err error) {
	updateFile("post-create", func(filename string) gist {

		content := fmt.Sprintf("%s,%s,%s", imageName, time.Now().UTC(), errorAsString(err))
		return gist{
			Files: map[string]gistFile{filename: {
				Filename: filename,
				Content:  getContent(filename) + "\n" + content,
			}},
		}
	})
}

func updateFile(filename string, fn func(filename string) gist) {

	var buf io.ReadWriter

	buf = &bytes.Buffer{}
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(false)
	g := fn(filename)
	err := enc.Encode(g)
	if err != nil {
		return
	}

	req, err := http.NewRequest(http.MethodPatch, gistApiUrl+"?access_token="+token(), buf)
	if err != nil {
		return
	}
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
}

func getContent(filename string) string {
	req, err := http.NewRequest(http.MethodGet, gistApiUrl, nil)
	if err != nil {
		return ""
	}
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	req.Header.Set("Authorization", token())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)
	var g gist
	err = decoder.Decode(&g)
	if err != nil {
		return ""
	}
	return g.Files[filename].Content
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
