package auth

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/adrg/xdg"
	"github.com/spf13/cobra"
	"github.com/way-platform/ileap-go"
)

// Credentials for the rFMS CLI.
type Credentials struct {
	BaseURL           string                  `json:"baseUrl"`
	ClientCredentials ileap.ClientCredentials `json:"clientCredentials,omitzero"`
}

func resolveCredentialsFilepath() (string, error) {
	return xdg.ConfigFile("ileap-go/auth.json")
}

// ReadCredentials reads the rFMS CLI credentials.
func ReadCredentials() (*Credentials, error) {
	credentialsFilepath, err := resolveCredentialsFilepath()
	if err != nil {
		return nil, err
	}
	if _, err := os.Stat(credentialsFilepath); err != nil {
		return nil, err
	}
	data, err := os.ReadFile(credentialsFilepath)
	if err != nil {
		return nil, err
	}
	var credentials Credentials
	if err := json.Unmarshal(data, &credentials); err != nil {
		return nil, err
	}
	return &credentials, nil
}

// NewClient creates a new iLEAP client using the CLI credentials.
func NewClient() (*ileap.Client, error) {
	auth, err := ReadCredentials()
	if err != nil {
		return nil, err
	}
	return ileap.NewClient(
		ileap.WithBaseURL(auth.BaseURL),
		ileap.WithReuseTokenAuth(auth.ClientCredentials),
	), nil
}

func writeCredentials(credentials *Credentials) error {
	credentialsFilepath, err := resolveCredentialsFilepath()
	if err != nil {
		return err
	}
	data, err := json.MarshalIndent(credentials, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(credentialsFilepath, data, 0o600)
}

func removeCredentials() error {
	credentialsFilepath, err := resolveCredentialsFilepath()
	if err != nil {
		return err
	}
	return os.RemoveAll(credentialsFilepath)
}

// NewCommand returns a new [cobra.Command] for rFMS CLI authentication.
func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "auth",
		Short: "Authenticate with an iLEAP API",
	}
	cmd.AddCommand(newLoginCommand())
	cmd.AddCommand(newLogoutCommand())
	return cmd
}

func newLoginCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "login",
		Short: "Login to an iLEAP API",
	}
	clientID := cmd.Flags().String("client-id", "", "client ID to use for authentication")
	cmd.MarkFlagRequired("client-id")
	clientSecret := cmd.Flags().String("client-secret", "", "client secret to use for authentication")
	cmd.MarkFlagRequired("client-secret")
	baseURL := cmd.Flags().String("base-url", "", "base URL to use for authentication")
	cmd.MarkFlagRequired("base-url")
	cmd.RunE = func(cmd *cobra.Command, _ []string) error {
		if !strings.HasPrefix(*baseURL, "http://") && !strings.HasPrefix(*baseURL, "https://") {
			return fmt.Errorf("--base-url must start with http:// or https://")
		}
		authenticator := ileap.NewOAuth2TokenAuthenticator(*clientID, *clientSecret, *baseURL)
		clientCredentials, err := authenticator.Authenticate(cmd.Context())
		if err != nil {
			return err
		}
		auth := &Credentials{
			BaseURL:           *baseURL,
			ClientCredentials: clientCredentials,
		}
		if err := writeCredentials(auth); err != nil {
			return err
		}
		cmd.Printf("Logged in to %s.\n", *baseURL)
		return nil
	}
	return cmd
}

func newLogoutCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "logout",
		Short: "Logout from the current authenticated iLEAP API",
		RunE: func(cmd *cobra.Command, _ []string) error {
			if err := removeCredentials(); err != nil {
				return err
			}
			cmd.Println("Logged out.")
			return nil
		},
	}
}
