package command

import (
	"fmt"
	"github.com/mhewedy/vermin/config"
	"github.com/spf13/cobra"
	"os"
)

func checkFilePath(path string) {
	if _, err := os.Stat(path); err != nil {
		fmt.Println("file not found", path)
		os.Exit(1)
	}
}

func preRun(cmd *cobra.Command, args []string) {
	config.CheckForUpdates(version)
}

func exitOnError(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
