package main

import (
	"github.com/mhewedy/vermin/db"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func init() {
	if err := os.MkdirAll(db.ImagesDir, 0755); err != nil {
		log.Fatal(err)
	}

	if err := os.MkdirAll(db.VMsBaseDir, 0755); err != nil {
		log.Fatal(err)
	}

	loadPrivateKey(
		"https://raw.githubusercontent.com/hashicorp/vagrant/master/keys/vagrant",
		db.VagrantPrivateKey,
	)
}

func loadPrivateKey(url, targetFile string) {

	keyPath := filepath.Join(db.BaseDir, targetFile)

	if _, err := os.Stat(keyPath); !os.IsNotExist(err) {
		return
	}

	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	if err = ioutil.WriteFile(keyPath, b, 0600); err != nil {
		log.Fatal(err)
	}
}
