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
	"strings"
)

// mountCmd represents the mount command
var mountCmd = &cobra.Command{
	Use:   "mount",
	Short: "Mount local filesystem inside the VM",
	Long: `Mount local filesystem to a directory inside the VM, 
if the guest directory is not specified, then /vermin is used

Usage: vermin mount <vm> local dir:[<vm dir>]

default <vm dir> is /vermin

Note: 
1. You can mount as many times as you want, and each time overrides old mounts.
2. Mounts are transient.
`,
	Example: `
To mount the ~/Downloads directory to /vermin inside the VM:
$ vermin mount vm_01 ~/Downloads

To mount the ~/MyHtmlProject directory to /var/www/html inside the VM:
$ vermin mount vm_01 ~/MyHtmlProject:/var/www/html
`,
	Run: func(cmd *cobra.Command, args []string) {

		vmName := args[0]
		path := args[1]
		remove, _ := cmd.Flags().GetBool("remove")

		p := strings.Split(path, ":")
		hostPath := p[0]
		checkFilePath(hostPath)

		guestPath := "/vermin"
		if len(p) > 1 {
			guestPath = p[1]
		}

		err := vms.Mount(vmName, hostPath, guestPath, remove)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("vm required")
		}
		if len(args) < 2 {
			return errors.New("path required")
		}
		return nil
	},
	ValidArgsFunction: listRunningVms,
}

var mountLsCmd = &cobra.Command{
	Use:   "ls",
	Short: "list mounted directories",
	Long:  "list mounted directories",
	Run: func(cmd *cobra.Command, args []string) {

		vmName := args[0]

		ps, err := vms.ListMounts(vmName)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Print(ps)
	},
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("image required")
		}
		return nil
	},
	ValidArgsFunction: listRunningVms,
}

func init() {
	rootCmd.AddCommand(mountCmd)
	mountCmd.AddCommand(mountLsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// mountCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	mountCmd.Flags().BoolP("remove", "r", false, "remove mounts before doing the new mount")
	//mountCmd.Flags().IntP("mem", "m", 1024, "Memory size in mega bytes")
}
