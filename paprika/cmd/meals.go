package cmd

import "github.com/spf13/cobra"

var mealsCmd = &cobra.Command{
	Use: "meals",
	Short: "meals",
}

func init() {
	rootCmd.AddCommand(mealsCmd)
}
