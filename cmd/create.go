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

const paramShell = "shell"
const paramAnsible = "ansible"

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new VM",
	Long:  "Create a new VM",
	Example: `
To list all available images:
$ vermin images

To create VM with default settings:
$ vermin create <image>

To create VM with default settings and provide a script to provision the VM:
$ vermin create <image> </path/to/shell/script.sh>
`,
	Run: func(cmd *cobra.Command, args []string) {
		imageName := args[0]

		var ps vms.ProvisionScript

		if len(args) > 1 {
			ps.Script = args[1]
			checkFilePath(ps.Script)

			t, _ := cmd.Flags().GetString("type")
			switch t {
			case paramShell:
				ps.Func = vms.ProvisionShell
			case paramAnsible:
				ps.Func = vms.ProvisionAnsible
			}

		}
		cpus, _ := cmd.Flags().GetInt("cpus")
		mem, _ := cmd.Flags().GetInt("mem")

		vmName, err := vms.Create(imageName, ps, cpus, mem)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		fmt.Printf("\nVM is ready, to connect to it use:\n$ vermin ssh %s\n", vmName)
	},
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("Image required\nUse the command 'vermin images' to list all images available")
		}

		typeStr, _ := cmd.Flags().GetString("type")
		if !(strings.EqualFold(typeStr, paramShell) || strings.EqualFold(typeStr, paramAnsible)) {
			return errors.New("provision script type should be either shell or ansible")
		}

		return nil
	},
	ValidArgsFunction: listImages,
}

func init() {
	rootCmd.AddCommand(createCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	createCmd.Flags().IntP("cpus", "c", 1, "Number of cpu cores")
	createCmd.Flags().IntP("mem", "m", 1024, "Memory size in mega bytes")
	createCmd.Flags().StringP("type", "t", "shell", "the type of provision script, can be shell or ansible")
}
