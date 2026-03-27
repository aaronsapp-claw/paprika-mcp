// commandspec:custom:start
package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/aarons22/paprika-tools/paprika/internal"
	"github.com/aarons22/paprika-tools/paprika/internal/client"
	"github.com/spf13/cobra"
)

var (
	loginEmail    string
	loginPassword string
)

var accountLoginLoginCmd = &cobra.Command{
	Use:   "login",
	Short: "Authenticate and save a Bearer token to ~/.config/paprika/config.yaml",
	Long: `Authenticate with the Paprika API and save the resulting Bearer token to
~/.config/paprika/config.yaml (mode 0600). Subsequent commands will read
the token from that file automatically.`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		if loginEmail == "" || loginPassword == "" {
			return fmt.Errorf("--email and --password are required")
		}

		c := client.NewClient("https://www.paprikaapp.com/api/v1", "")

		body := map[string]string{"email": loginEmail, "password": loginPassword}
		resp, err := c.Do("POST", "/account/login/", nil, nil, body)
		if err != nil {
			return err
		}

		// Extract token from {"result": {"token": "..."}}
		var parsed struct {
			Result struct {
				Token string `json:"token"`
			} `json:"result"`
		}
		if err := json.Unmarshal(resp, &parsed); err != nil || parsed.Result.Token == "" {
			return fmt.Errorf("unexpected response (no token): %s", string(resp))
		}

		cfg := &internal.Config{Token: parsed.Result.Token}
		path, err := internal.Save("paprika", cfg)
		if err != nil {
			return fmt.Errorf("saving token: %w", err)
		}

		fmt.Fprintf(os.Stderr, "Token saved to %s\n", path)
		return nil
	},
}

func init() {
	accountLoginCmd.AddCommand(accountLoginLoginCmd)
	accountLoginLoginCmd.Flags().StringVar(&loginEmail, "email", "", "Paprika account email (required)")
	accountLoginLoginCmd.Flags().StringVar(&loginPassword, "password", "", "Paprika account password (required)")
}

// commandspec:custom:end
