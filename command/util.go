package command

import (
	"fmt"
	"github.com/mhewedy/vermin/config"
	"github.com/spf13/cobra"
	"os"
	"regexp"
	"strings"
)

var vmNameRegexDash = regexp.MustCompile("^vm-[0-9]+$")
var vmNameRegexNoVM = regexp.MustCompile("^[0-9]+$")
var vmNameRegexNoSpace = regexp.MustCompile("^vm[0-9]+$")

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

func normalizeVmName(vmName string) string {

	if vmNameRegexDash.MatchString(vmName) {
		return strings.ReplaceAll(vmName, "-", "_")
	}

	if vmNameRegexNoVM.MatchString(vmName) {
		return fmt.Sprintf("vm_%02s", vmName)
	}

	if vmNameRegexNoSpace.MatchString(vmName) {
		return strings.ReplaceAll(vmName, "vm", "vm_")
	}

	return vmName
}
