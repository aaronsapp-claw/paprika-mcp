package cmd

import (
	"fmt"
	"os"
	"github.com/spf13/cobra"
	"github.com/aarons22/paprika-mcp/paprika/internal/client"
	"github.com/aarons22/paprika-mcp/paprika/internal/output"
)

var statusGetSyncStatusCmd = &cobra.Command{
	Use: "getSyncStatus",
	Short: "Get sync status counters",
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		baseURL, _ := cmd.Root().PersistentFlags().GetString("base-url")
		token := os.Getenv("PAPRIKA_TOKEN")
		c := client.NewClient(baseURL, token)
		pathParams := map[string]string{}
		queryParams := map[string]string{}
		resp, err := c.Do("GET", "/status/", pathParams, queryParams, nil)
		if err != nil {
			return err
		}
		jsonMode, _ := cmd.Root().PersistentFlags().GetBool("json")
		noColor, _ := cmd.Root().PersistentFlags().GetBool("no-color")
		if jsonMode {
			fmt.Printf("%s\n", string(resp))
		} else {
			if err := output.PrintTable(resp, noColor); err != nil {
				fmt.Println(string(resp))
			}
		}
		return nil
	},
}

func init() {
	statusCmd.AddCommand(statusGetSyncStatusCmd)
}
