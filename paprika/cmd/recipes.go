package cmd

import "github.com/spf13/cobra"

var recipesCmd = &cobra.Command{
	Use: "recipes",
	Short: "recipes",
}

func init() {
	rootCmd.AddCommand(recipesCmd)
}
