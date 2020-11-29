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

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start one or more stopped VMs",
	Long:  `Start one or more stopped VMs`,
	Run: func(cmd *cobra.Command, args []string) {
		for _, vmName := range args {
			err := vms.Start(normalizeVmName(vmName))
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}
	},
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("vm required")
		}
		return nil
	},
	ValidArgsFunction: listStoppedVms,
}

func init() {
	rootCmd.AddCommand(startCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// startCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	//startCmd.Flags().BoolP("purge", "p", false, "Purge the IP cache")
}
