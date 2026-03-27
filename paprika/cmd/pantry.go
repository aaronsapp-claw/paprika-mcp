package cmd

import "github.com/spf13/cobra"

var pantryCmd = &cobra.Command{
	Use: "pantry",
	Short: "pantry",
}

func init() {
	rootCmd.AddCommand(pantryCmd)
}
