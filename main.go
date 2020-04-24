package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"os"
	"vermin/db"
)

func init() {
	if err := os.MkdirAll(db.GetImagesDir(), 0755); err != nil {
		log.Fatal(err)
	}

	if err := os.MkdirAll(db.GetVMsBaseDir(), 0755); err != nil {
		log.Fatal(err)
	}
}

func main() {

	//fmt.Println(images.List())
	fmt.Println(ps(true))

	//fmt.Println(create("centos/8", "", 0, 0))

	/*rootCmd := &cobra.Command{}

	cmd := newCommand()
	rootCmd.AddCommand(cmd)

	cmd.AddCommand(newNestedCommand())

	if err := rootCmd.Execute(); err != nil {
		println(err.Error())
	}*/
}

func newCommand() *cobra.Command {
	cmd := &cobra.Command{
		Run: func(cmd *cobra.Command, args []string) {
			println(`Foo`)
		},
		Use:   `foo`,
		Short: "Command foo",
		Long:  "This is a command",
	}

	return cmd
}

func newNestedCommand() *cobra.Command {
	cmd := &cobra.Command{
		Run: func(cmd *cobra.Command, args []string) {
			println(`Bar`)
		},
		Use:   `bar`,
		Short: "Command bar",
		Long:  "This is a nested command",
	}

	return cmd
}
