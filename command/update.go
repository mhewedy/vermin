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
	"github.com/spf13/cobra"
	"os"
)

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:     "update",
	Aliases: []string{"modify"},
	Short:   "Update configuration of a VM",
	Long:    "Update configuration of a VM",
	Example: `

You can either change the cpu and/or mem or you can use the update command to shrink the disk.

To change the VM to use 2 cores and 512MB memory
$ vermin update vm_01 --cpus 2 --mem 512

To shrink the disk size on the host machine
$ vermin update vm_01 --shrink-disk
`,
	Run: func(cmd *cobra.Command, args []string) {
		vmName := normalizeVmName(args[0])

		shrink, _ := cmd.Flags().GetBool("shrink-disk")

		if shrink {
			if err := vms.Shrink(vmName); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		} else {
			cpus, _ := cmd.Flags().GetInt("cpus")
			mem, _ := cmd.Flags().GetInt("mem")

			if err := vms.Modify(vmName, cpus, mem); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}
	},
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("vm required")
		}

		shrink, _ := cmd.Flags().GetBool("shrink-disk")
		if !shrink {
			cpus, _ := cmd.Flags().GetInt("cpus")
			mem, _ := cmd.Flags().GetInt("mem")

			if cpus == 0 && mem == 0 {
				return errors.New("should specify cpus and/or mem specs")
			}
		}
		return nil
	},
	ValidArgsFunction: listImages,
}

func init() {
	rootCmd.AddCommand(updateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// updateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	updateCmd.Flags().IntP("cpus", "c", 0, "Number of cpu cores")
	updateCmd.Flags().IntP("mem", "m", 0, "Memory size in mega bytes")
	updateCmd.Flags().BoolP("shrink-disk", "", false, "Shrink the disk to reduce the size on the host machine")
}
