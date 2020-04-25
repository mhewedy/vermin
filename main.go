package main

import (
	"github.com/mhewedy/vermin/cli"
	"github.com/mhewedy/vermin/db"
	"log"
	"os"
)

func init() {
	if err := os.MkdirAll(db.GetImagesDir(), 0755); err != nil {
		log.Fatal(err)
	}

	if err := os.MkdirAll(db.GetVMsBaseDir(), 0755); err != nil {
		log.Fatal(err)
	}
}

func main() {
	cli.Execute()
}
