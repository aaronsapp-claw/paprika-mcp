package cmd

import "github.com/spf13/cobra"

var grocerylistsCmd = &cobra.Command{
	Use: "grocerylists",
	Short: "grocerylists",
}

func init() {
	rootCmd.AddCommand(grocerylistsCmd)
}
