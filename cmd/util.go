package cmd

import (
	"fmt"
	"os"
)

func checkFilePath(path string) {
	if _, err := os.Stat(path); err != nil {
		fmt.Println("file not found", path)
		os.Exit(1)
	}
}
