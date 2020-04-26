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
	"github.com/spf13/cobra"
	"os"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create VM from an image",
	Long: `Create VM from an image

To list all available images:
$ vermin images

To create VM with default settings:
$ vermin create <image>

To create VM with default settings and provide a script to provision the VM:
$ vermin create <image> </path/to/shell/script.sh>
`,
	Run: func(cmd *cobra.Command, args []string) {
		imageName := args[0]
		var script string
		if len(args) > 1 {
			script = args[1]
		}
		cpus, _ := cmd.Flags().GetInt("cpus")
		mem, _ := cmd.Flags().GetInt("mem")

		vmName, err := vms.Create(imageName, script, cpus, mem)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Printf("Image created successfuly.\nUse the command: 'vermin start %s' to start the vm."+
			"\nThen use the command 'vermin ssh %s' to use the vm.\n", vmName, vmName)
	},
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("image required\nUse the command 'vermin images' to list all images available")
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
}
