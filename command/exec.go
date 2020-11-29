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
	"github.com/mhewedy/vermin/vms"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

// execCmd represents the exec command
var execCmd = &cobra.Command{
	Use:   "exec",
	Short: "Run a command in a running VM",
	Long:  `Run a command in a running VM`,
	Example: `
Execute remote command:
$ vermin exec vm_09 cat /etc/passwd
`,
	Run: func(cmd *cobra.Command, args []string) {
		vmName := normalizeVmName(args[0])
		command := strings.Join(args[1:], " ")

		err := vms.Exec(vmName, command)
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
			return errors.New("command required")
		}
		return nil
	},
	ValidArgsFunction: listRunningVms,
}

func init() {
	rootCmd.AddCommand(execCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// execCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	//execCmd.Flags().BoolP("purge", "p", false, "Purge the IP cache")
}
