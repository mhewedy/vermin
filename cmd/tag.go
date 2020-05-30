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

// tagCmd represents the tag command
var tagCmd = &cobra.Command{
	Use:   "tag",
	Short: "Add or remove tag to a VM",
	Long: `Add or remove tag to a VM
You can tag a VM as many times as you want
`,
	Run: func(cmd *cobra.Command, args []string) {
		vmName := args[0]
		tag := args[1]

		remove, _ := cmd.Flags().GetBool("remove")

		err := vms.Tag(vmName, tag, remove)
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
			return errors.New("tag required")
		}
		return nil
	},
	ValidArgsFunction: listAllVms,
}

func init() {
	rootCmd.AddCommand(tagCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// tagCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	tagCmd.Flags().BoolP("remove", "r", false, "remove tag")
}
