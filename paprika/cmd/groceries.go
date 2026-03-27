package cmd

import "github.com/spf13/cobra"

var groceriesCmd = &cobra.Command{
	Use: "groceries",
	Short: "groceries",
}

func init() {
	rootCmd.AddCommand(groceriesCmd)
}
