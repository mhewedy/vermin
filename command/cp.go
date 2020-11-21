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
package command

import (
	"errors"
	"fmt"
	"github.com/mhewedy/vermin/cmd/scp"
	"github.com/mhewedy/vermin/vms"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

// cpCmd represents the cp command
var cpCmd = &cobra.Command{
	Use:   "cp",
	Short: "Copy files/folders between a VM and the local filesystem or between two VMs",
	Long: `Copy files/folders between a VM and the local filesystem or between two VMs

vermin cp <source> <destination>

where: <source> and <destination> take the form of <vm_name>:/path/to/file/or/directory
In case of the <source> or <destination> is the local host, then the "<vm_name>:" part is removed.

Note: You need to have appropriate permissions on the files you need to copy.
`,
	Example: `
Copy file.txt from host to user's home directory inside the vm
$ vermin cp ~/file.txt vm_01:~ 

Copy file.txt from user's home directory inside the vm to the current directory on the host
$ vermin cp vm_01:project/file.txt . 

Copy /etc/os-release from the inside the vm to ~/temp on the host
$ vermin cp vm_02:/etc/os-release ~/temp

Copy file.txt from vm_01 home dir to the vm_02 /tmp dir
$ vermin cp vm_01:~/file.txt vm_02:/tmp 

`,
	Run: func(cmd *cobra.Command, args []string) {
		source := args[0]
		destination := args[1]

		if err := vms.CopyFiles(source, destination); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("source required")
		}
		if len(args) < 2 {
			return errors.New("destination required")
		}

		if !strings.Contains(args[0], scp.CopySeparator) &&
			!strings.Contains(args[1], scp.CopySeparator) {
			return errors.New(`either source or destination should be a VM path in the form of "vm_name:/path/to/copy"`)
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
	//cpCmd.Flags().StringP("local", "l", "", "Local file/folder to copy to VM home directory")
	//cpCmd.Flags().StringP("remote", "r", "", "VM file/folder to copy to host at current directory")
}
