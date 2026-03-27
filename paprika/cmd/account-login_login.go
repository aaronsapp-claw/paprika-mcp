package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"github.com/spf13/cobra"
	"github.com/aarons22/paprika-mcp/paprika/internal/client"
	"github.com/aarons22/paprika-mcp/paprika/internal/output"
)

var (
	accountLoginLoginCmdBody string
	accountLoginLoginCmdBodyFile string
)

var accountLoginLoginCmd = &cobra.Command{
	Use: "login",
	Short: "Authenticate and obtain a Bearer token",
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		baseURL, _ := cmd.Root().PersistentFlags().GetString("base-url")
		token := os.Getenv("PAPRIKA_TOKEN")
		c := client.NewClient(baseURL, token)
		pathParams := map[string]string{}
		queryParams := map[string]string{}
		if accountLoginLoginCmdBodyFile != "" {
			fileData, err := os.ReadFile(accountLoginLoginCmdBodyFile)
			if err != nil {
				return fmt.Errorf("reading body-file: %w", err)
			}
			if !json.Valid(fileData) {
				return fmt.Errorf("body-file does not contain valid JSON")
			}
			accountLoginLoginCmdBody = string(fileData)
		}
		if accountLoginLoginCmdBody != "" {
			if !json.Valid([]byte(accountLoginLoginCmdBody)) {
				return fmt.Errorf("--body does not contain valid JSON")
			}
			var bodyObj interface{}
			_ = json.Unmarshal([]byte(accountLoginLoginCmdBody), &bodyObj)
			resp, err := c.Do("POST", "/account/login/", pathParams, queryParams, bodyObj)
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
		}
		resp, err := c.Do("POST", "/account/login/", pathParams, queryParams, nil)
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
	accountLoginCmd.AddCommand(accountLoginLoginCmd)
	accountLoginLoginCmd.Flags().StringVar(&accountLoginLoginCmdBody, "body", "", "Raw JSON body (overrides individual flags)")
	accountLoginLoginCmd.Flags().StringVar(&accountLoginLoginCmdBodyFile, "body-file", "", "Path to JSON file to use as request body")
}
