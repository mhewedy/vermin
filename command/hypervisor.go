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
	"fmt"
	"github.com/mhewedy/vermin/hypervisor"
	"github.com/spf13/cobra"
)

// hypervisorCmd  represents the hypervisor command
var hypervisorCmd = &cobra.Command{
	Use:   "hypervisor",
	Short: "print the name of the detected hypervisor",
	Long:  `print the name of the detected hypervisor.`,
	Run: func(cmd *cobra.Command, args []string) {

		h, err := hypervisor.GetHypervisorName(false)
		exitOnError(err)
		fmt.Println(h)
	},
}

func init() {
	rootCmd.AddCommand(hypervisorCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// commitCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	//commitCmd.Flags().BoolP("override", "", false, "override any existing image")
}
