package main

import (
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
	//fmt.Println(ps(false))

	//err := start("vm_01")
	//if err != nil {
	//	log.Fatal(err)
	//}

	//fmt.Println(ssh("vm_02"))

	//fmt.Println(create("ubuntu/bionic", "/Users/mhewedy/Work/Code/viper/samples/provision/sample.sh", 0, 0))
	//vmName, err := create("ubuntu/bionic", "", 0, 0)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//if err := start(vmName); err != nil {
	//	log.Fatal(err)
	//}

	if err := secureShell("vm_02", ""); err != nil {
		log.Fatal(err)
	}

	//fmt.Println(remove("vm_11"))

	//fmt.Println(
	//	provision("vm_12", "/Users/mhewedy/Work/Code/viper/samples/provision/sample.sh"),
	//)

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
