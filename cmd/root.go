package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/abrimentcloud/kubectl-abriment/config"
	"github.com/abrimentcloud/kubectl-abriment/pkg/login"
	"github.com/abrimentcloud/kubectl-abriment/pkg/logout"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var rootCmd = &cobra.Command{
	Use:     "kubectl-abriment",
	Aliases: []string{"i"},
	Short:   "Run in interactive mode",
	Long:    "Run the plugin in interactive mode with guided prompts",
	Run:     interactive,
}

var helpCmd = &cobra.Command{
	Use:   "help",
	Short: "Display help information",
	Long:  "Display detailed help information about the kubectl-abriment plugin and its usage",
	Run:   help,
}

func Execute() {
	loginCmd.Flags().StringVarP(&user, "username", "u", "", "provide your username")
	loginCmd.Flags().StringVarP(&pass, "password", "p", "", "provide your password")
	loginCmd.Flags().StringVarP(&token, "token", "t", "", "provide your token")
	loginCmd.Flags().StringVar(&dryrun, "dry-run", "", "options: client.")

	rootCmd.AddCommand(helpCmd)
	rootCmd.AddCommand(loginCmd)
	rootCmd.AddCommand(logoutCmd)

	cobra.CheckErr(rootCmd.Execute())
}

func help(cmd *cobra.Command, args []string) {
	fmt.Print(`
kubectl-abriment - Kubernetes Plugin for Abriment

DESCRIPTION:
    This plugin allows you to authenticate with the Abriment backend service and
    automatically configure your kubeconfig file with the necessary cluster,
    context, and user information.

USAGE:
    kubectl abriment
	kubectl abriment login [flags]
    kubectl abriment logout
    kubectl abriment help

AUTHENTICATION OPTIONS:
    You can authenticate using either:
    1. Username and password
    2. Token-based authentication

LOGIN FLAGS:
    -u, --username string    Your username for authentication
    -p, --password string    Your password for authentication
    -t, --token string       Your authentication token
        --dry-run string     Options: client (prints config without saving)

COMMANDS:
    help             Show this help message
    login            Login with providing credentials and udpates the config file
    logout           Logout modifies the config file and deletes the abriment configs

EXAMPLES:
    # Login with username and password
    kubectl abriment login -u myuser -p mypassword

    # Login with token
    kubectl abriment login -t mytoken

    # Preview config without saving (dry-run)
    kubectl abriment login -u myuser -p mypassword --dry-run client

    # Logout
    kubectl abriment logout

    # Interactive mode
    kubectl abriment

ENVIRONMENT VARIABLES:
    LOGIN_ENDPOINT     Backend login endpoint (default: https://backend.abriment.com/dashboard/api/login/)
    CONFIG_ENDPOINT    Backend config endpoint (default: https://backend.abriment.com/api/v1/paas/kubeconfig/)
    KUBECONFIG        Custom path for kubeconfig file

BEHAVIOR:
    1. Authenticates with the backend service using provided credentials
    2. Retrieves the kubeconfig from the backend
    3. If kubeconfig exists locally, merges the new configuration
    4. If kubeconfig doesn't exist, creates a new one
    5. Adds/updates cluster, context, and user information for Abriment

NOTES:
    - The plugin will merge configurations if an existing kubeconfig is found
    - It preserves existing configurations while adding Abriment-specific resources
    - Default kubeconfig location: ~/.kube/config (unless KUBECONFIG is set)
`)
}

func interactive(cmd *cobra.Command, args []string) {

	// Show current configuration
	cfg, err := config.ParseCfg()
	if err != nil {
		fmt.Printf("❌ Error loading configuration: %v\n", err)
		return
	}

	fmt.Printf("Backend Configuration:\n")
	fmt.Printf("   Login Endpoint:  %s\n", cfg.LoginEndpoint)
	fmt.Printf("   Config Endpoint: %s\n\n", cfg.ConfigEndpoint)

	// Ask for authentication method
	var operation string
	opPrompt := &survey.Select{
		Message: "Which operation would you like to do?",
		Options: []string{
			"Login",
			"Logout",
		},
	}
	survey.AskOne(opPrompt, &operation)

	if operation == "Logout" {
		if err := logout.RemoveAbrimentFromConfigfile(); err != nil {
			fmt.Println("logout failed.")
			return
		}

		fmt.Println("Logged out successfully!")
		return
	}

	if operation != "Login" {
		fmt.Println("Invalid Operation!")
		return
	}

	// Ask for authentication method
	var authMethod string
	authPrompt := &survey.Select{
		Message: "How would you like to authenticate?",
		Options: []string{
			"Username and Password",
			"Token",
		},
	}
	survey.AskOne(authPrompt, &authMethod)

	var body login.LoginBody

	if authMethod == "Username and Password" {
		// Get username
		var username string
		usernamePrompt := &survey.Input{
			Message: "Enter your username:",
		}
		survey.AskOne(usernamePrompt, &username)

		// Get password (hidden input)
		fmt.Print("Enter your password: ")
		passwordBytes, err := term.ReadPassword(int(syscall.Stdin))
		if err != nil {
			fmt.Printf("\n❌ Error reading password: %v\n", err)
			return
		}
		password := string(passwordBytes)
		fmt.Println() // New line after password input

		body = login.LoginBody{
			Username: username,
			Password: password,
		}
	} else if authMethod == "Token" {
		// Get token
		var tokenInput string
		tokenPrompt := &survey.Input{
			Message: "Enter your authentication token:",
		}
		survey.AskOne(tokenPrompt, &tokenInput)

		body = login.LoginBody{
			UnsocpedToken: tokenInput,
		}
	} else {
		fmt.Println("Invalid Authentication Method!")
		return
	}

	// Marshal request body
	bb, err := json.Marshal(body)
	if err != nil {
		fmt.Printf("❌ Error preparing request: %v\n", err)
		return
	}

	// Login request
	bodyBytes := bytes.NewReader(bb)
	res, err := login.Login(bodyBytes, cfg.LoginEndpoint)
	if err != nil {
		fmt.Printf("❌ Authentication failed: %v\n", err)
		return
	}

	// Get YAML config
	yamlBytes, err := login.GetYamlConfig(res.Data.Token.ID, cfg)
	if err != nil {
		fmt.Printf("❌ Error retrieving config: %v\n", err)
		return
	}

	fmt.Println("Authentication successful!")

	// Ask about dry-run
	var dryRunChoice bool
	dryRunPrompt := &survey.Confirm{
		Message: "Would you like to preview the configuration without saving it? (dry-run)",
		Default: false,
	}
	survey.AskOne(dryRunPrompt, &dryRunChoice)

	if dryRunChoice {
		fmt.Println("\nKubeconfig Preview:")
		fmt.Println("=" + strings.Repeat("=", 50))
		fmt.Println(string(yamlBytes))
		fmt.Println("=" + strings.Repeat("=", 50))
		return
	}

	// Save config
	if err := login.SaveConfigToConfigfile(yamlBytes); err != nil {
		fmt.Printf("❌ Error saving config: %v\n", err)
		return
	}

	fmt.Println("Configuration saved successfully!")

	// Show final status
	kubeconfigPath := os.Getenv("KUBECONFIG")
	if kubeconfigPath == "" {
		home, _ := os.UserHomeDir()
		kubeconfigPath = filepath.Join(home, ".kube", "config")
	}

	fmt.Printf("\n kubeconfig has been updated.\n")
	fmt.Printf("   Location: %s\n", kubeconfigPath)
	fmt.Printf("   Cluster: %s\n", config.AbrimentCluster)
	fmt.Printf("   Context: %s\n", config.AbrimentContext)
	fmt.Printf("   User: %s\n\n", config.AbrimentUser)
}
