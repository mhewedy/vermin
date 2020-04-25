package cli

import (
	"errors"
	"github.com/spf13/cobra"
	"os"
)

var completionCmd = &cobra.Command{
	Use:   "completion",
	Short: "Generates completion scripts (bash and zsh)",
	Long: `To load completion run
. <(vermin completion bash)
or 
. <(vermin completion zsh)

You might consider adding it to your shell config file (.bashrc or .zshrc)
`,
	Run: func(cmd *cobra.Command, args []string) {
		switch args[0] {
		case "bash":
			_ = rootCmd.GenBashCompletion(os.Stdout)
		case "zsh":
			_ = rootCmd.GenZshCompletion(os.Stdout)
		}
	},
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 || !(args[0] == "bash" || args[0] == "zsh") {
			return errors.New("should choose between bash or zsh")
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(completionCmd)
}
