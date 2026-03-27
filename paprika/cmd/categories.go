package cmd

import "github.com/spf13/cobra"

var categoriesCmd = &cobra.Command{
	Use: "categories",
	Short: "categories",
}

func init() {
	rootCmd.AddCommand(categoriesCmd)
}
