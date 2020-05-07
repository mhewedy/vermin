/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"errors"
	"fmt"
	"github.com/mhewedy/vermin/vms"
	"github.com/spf13/cobra"
	"os"
)

// mountCmd represents the mount command
var mountCmd = &cobra.Command{
	Use:   "mount",
	Short: "Mount host path inside the VM",
	Long: `Mount host path to /vermin directory inside the VM

To mount the ~/Downloads directory to /vermin inside the VM:
$ vermin mount vm_01 ~/Downloads

Note: 
1. You can mount as many times as you want, and each time old mounts are removed.
2. Mounts are transient.
`,
	Run: func(cmd *cobra.Command, args []string) {
		vmName := args[0]
		hostPath := args[1]
		checkFilePath(hostPath)

		err := vms.Mount(vmName, hostPath)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 2 {
			return errors.New("Invalid number of arguments passed,\nUse vermin mount -h to show help.")
		}
		return nil
	},
	ValidArgsFunction: listRunningVms,
}

func init() {
	rootCmd.AddCommand(mountCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// mountCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	//mountCmd.Flags().IntP("cpus", "c", 1, "Number of cpu cores")
	//mountCmd.Flags().IntP("mem", "m", 1024, "Memory size in mega bytes")
}
