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

// modifyCmd represents the modify command
var modifyCmd = &cobra.Command{
	Use:   "modify",
	Short: "Modify a VM HW specs (cpu, memory)",
	Long:  "Modify a VM HW specs (cpu, memory)",
	Example: `

To change the VM to use 2 cores and 512MB memory
$ vermin modify vm_01 --cpus 2 --mem 512
`,
	Run: func(cmd *cobra.Command, args []string) {
		vmName := args[0]
		var script string
		if len(args) > 1 {
			script = args[1]
			checkFilePath(script)
		}
		cpus, _ := cmd.Flags().GetInt("cpus")
		mem, _ := cmd.Flags().GetInt("mem")

		if err := vms.Modify(vmName, cpus, mem); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("vm required")
		}

		cpus, _ := cmd.Flags().GetInt("cpus")
		mem, _ := cmd.Flags().GetInt("mem")

		if cpus == 0 && mem == 0 {
			return errors.New("should specify cpus and/or mem specs")
		}
		return nil
	},
	ValidArgsFunction: listImages,
}

func init() {
	rootCmd.AddCommand(modifyCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// modifyCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	modifyCmd.Flags().IntP("cpus", "c", 0, "Number of cpu cores")
	modifyCmd.Flags().IntP("mem", "m", 0, "Memory size in mega bytes")
}
