package cmd

import "github.com/spf13/cobra"

var statusCmd = &cobra.Command{
	Use: "status",
	Short: "status",
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
