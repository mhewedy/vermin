package main

import (
	"fmt"
	"github.com/spf13/cobra"
)

func main() {

	//fmt.Println(ip.Find("vm_01", true))
	//

	//fmt.Println(info.List([]string{"vm_01", "vm_02"}))

	fmt.Println(ps(true))

	//fmt.Println(images.List())

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
