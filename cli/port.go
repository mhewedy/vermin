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
package cli

import (
	"errors"
	"fmt"
	"github.com/mhewedy/vermin/vms"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

// portCmd represents the port command
var portCmd = &cobra.Command{
	Use:   "port",
	Short: "Forward port(s) from a VM to host",
	Long: `Forward port(s) from a VM to host

Usage: vermin port <vm> <vm port>[:local port] [<vm port>[:local port]]

Examples:

Forward vm port 4040 to local port 4040:
$ vermin port vm_01 4040

Forward vm port 4040 to local port 40040:
$ vermin port vm_01 4040:40040

Forward vm port 4040 to local port 40040 and port 8080 to 8080
$ vermin port vm_01 4040:40040 8080

Forward vm port 4040 to local port 40040 and ports in range (8080 to 8088) to range(8080 to 8088):
$ vermin port vm_01 4040:40040 8080-8088

Forward vm port 4040 to local port 40040 and ports in range (8080 to 8088) to range(9080 to 9088):
$ vermin port vm_01 4040:40040 8080-8088:9080-9088
`,
	Run: func(cmd *cobra.Command, args []string) {
		vmName := args[0]
		command := strings.Join(args[1:], " ")

		err := vms.PortForward(vmName, command)
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
			return errors.New("ports required")
		}
		return nil
	},
	ValidArgsFunction: listRunningVms,
}

func init() {
	rootCmd.AddCommand(portCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// portCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	//portCmd.Flags().BoolP("purge", "p", false, "Purge the IP cache")
}
