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
	"github.com/mhewedy/vermin/images"
	"os"

	"github.com/spf13/cobra"
)

// imagesCmd represents the images command
var imagesCmd = &cobra.Command{
	Use:     "images",
	Aliases: []string{"image"},
	Short:   "List remote and cached images",
	Long: `List remote and cached images
	
Images are cached after the first time it is downloaded.
You can find images from Vagrant at: https://app.vagrantup.com/search
`,
	Example: `
Use the image in creating a VM:
$ vermin create <image>
`,
	Run: func(cmd *cobra.Command, args []string) {
		i, err := images.Display()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Print(i)
	},
}

func init() {
	rootCmd.AddCommand(imagesCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// imagesCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	//imagesCmd.Flags().BoolP("purge", "p", false, "Purge images list cache")
}
