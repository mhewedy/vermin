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
	"github.com/mhewedy/vermin/vms"
	"github.com/spf13/cobra"
)

// commitCmd represents the commit command
var commitCmd = &cobra.Command{
	Use:   "commit",
	Short: "Commit a VM into a new Image",
	Long:  `Commit a VM into a new Image to be used later as a template to create VMs from.`,
	Example: `
Commit vm_01 that have installed elastic inside as an image so to be used later to obtain a VM that contains elastic installed.
$ vermin commit vm_01 elk/elastic
`,
	Run: func(cmd *cobra.Command, args []string) {

		vmName := args[0]
		imageName := args[1]
		override, _ := cmd.Flags().GetBool("override")

		err := vms.Commit(vmName, imageName, override)
		exitOnError(err)
	},
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("vm required")
		}
		if len(args) < 2 {
			return errors.New("new image name required, and should be in format base/name, example k8s/worker")
		}
		return nil
	},
	ValidArgsFunction: listStoppedVms,
}

func init() {
	rootCmd.AddCommand(commitCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// commitCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	commitCmd.Flags().BoolP("override", "", false, "override any existing image")
}
