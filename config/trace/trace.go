package trace

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const gistUrl = "https://api.github.com/gists/094cbbd974c92fe70ab4cc3443c5ac5f?access_token=352265f35826a92b8e35e3d8b95e22ec0ca7d092"

// Gist represents a GitHub's gist.
type gist struct {
	Files map[string]gistFile `json:"files,omitempty"`
}

type gistFile struct {
	Filename string `json:"filename,omitempty"`
	Content  string `json:"content,omitempty"`
}

func Create(imageName string) {
	go func() {
		content := fmt.Sprintf("Image: %s, Created at %s", imageName, time.Now().UTC())
		gist := gist{
			Files: map[string]gistFile{"vermin-trace": {Filename: "vermin-trace", Content: content}},
		}

		var buf io.ReadWriter

		buf = &bytes.Buffer{}
		enc := json.NewEncoder(buf)
		err := enc.Encode(gist)
		if err != nil {
			return
		}

		req, err := http.NewRequest(http.MethodPatch, gistUrl, buf)
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
	}()
}
