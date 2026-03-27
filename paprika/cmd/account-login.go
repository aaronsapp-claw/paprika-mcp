package cmd

import "github.com/spf13/cobra"

var accountLoginCmd = &cobra.Command{
	Use:   "account",
	Short: "account",
}

func init() {
	rootCmd.AddCommand(accountLoginCmd)
}
