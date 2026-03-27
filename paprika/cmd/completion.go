package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var completionCmd = &cobra.Command{
	Use:       "completion [bash|zsh|fish|powershell]",
	Short:     "Generate shell completion scripts",
	Args:      cobra.ExactArgs(1),
	ValidArgs: []string{"bash", "zsh", "fish", "powershell"},
	RunE: func(cmd *cobra.Command, args []string) error {
		root := cmd.Root()
		switch args[0] {
		case "bash":
			return root.GenBashCompletion(os.Stdout)
		case "zsh":
			return root.GenZshCompletion(os.Stdout)
		case "fish":
			return root.GenFishCompletion(os.Stdout, true)
		case "powershell":
			return root.GenPowerShellCompletion(os.Stdout)
		default:
			return fmt.Errorf("unsupported shell %q: must be one of bash, zsh, fish, powershell", args[0])
		}
	},
}

func init() {
	rootCmd.AddCommand(completionCmd)
}
