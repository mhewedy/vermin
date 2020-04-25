package main

import (
	"fmt"
	"github.com/mhewedy/vermin/commands"
	"github.com/mhewedy/vermin/db"
	"github.com/spf13/cobra"
	"log"
	"os"
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

	/*rootCmd := &cobra.Command{}

	cmd := newCommand()
	rootCmd.AddCommand(cmd)

	cmd.AddCommand(newNestedCommand())

	if err := rootCmd.Execute(); err != nil {
		println(err.Error())
	}*/
}

func psCommand() *cobra.Command {
	cmd := &cobra.Command{
		Run: func(cmd *cobra.Command, args []string) {
			ps, err := commands.Ps(false)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(ps)
		},
		Use:   `ps`,
		Short: "List running VMs",
		Long:  "List running VMs, use -a to list all VMs",
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
