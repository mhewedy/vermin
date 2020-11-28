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
	"github.com/mhewedy/vermin/images"
	"github.com/spf13/cobra"
	"os"
)

// removeImageCmd represents the remove image command
var removeImageCmd = &cobra.Command{
	Use:   "rmi",
	Short: "Remove one or more Image",
	Long:  `Remove one or more Image`,
	Run: func(cmd *cobra.Command, args []string) {
		for _, vmName := range args {
			//force, _ := cmd.Flags().GetBool("force")

			err := images.Remove(vmName)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}
	},
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("image required")
		}
		return nil
	},
	ValidArgsFunction: listImages,
}

func init() {
	rootCmd.AddCommand(removeImageCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// removeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	//removeCmd.Flags().BoolP("force", "f", false, "force remove running VM")
}
