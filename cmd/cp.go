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

// cpCmd represents the cp command
var cpCmd = &cobra.Command{
	Use:   "cp",
	Short: "Copy files/folders between a VM and the local filesystem",
	Long:  "Copy files/folders between a VM and the local filesystem",
	Example: `
Copy file.txt from host to user's home directory inside the vm
$ vermin cp vm_01 -l ~/file.txt

Copy file.txt from user's home directory inside the vm to the current directory on the host
$ vermin cp vm_01 -r project/file.txt

Copy /etc/os-release from the inside the vm to the current directory on the host
$ vermin cp vm_02 -r /etc/os-release
`,
	Run: func(cmd *cobra.Command, args []string) {
		vmName := args[0]

		var err error

		l, _ := cmd.Flags().GetString("local")
		r, _ := cmd.Flags().GetString("remote")

		if len(l) > 0 {
			err = vms.CopyFiles(vmName, l, true)
		} else if len(r) > 0 {
			err = vms.CopyFiles(vmName, r, false)
		} else {
			fmt.Println("missing file parameter, use -h for help.")
			os.Exit(1)
		}

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("vm required")
		}
		return nil
	},
	ValidArgsFunction: listRunningVms,
}

func init() {
	rootCmd.AddCommand(cpCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// cpCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	cpCmd.Flags().StringP("local", "l", "", "Local file/folder to copy to VM home directory")
	cpCmd.Flags().StringP("remote", "r", "", "VM file/folder to copy to host at current directory")
}
