package cmd

import (
	"fmt"

	"github.com/abrimentcloud/kubectl-abriment/pkg/logout"

	"github.com/spf13/cobra"
)

var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Logout from the Kubernetes cluster",
	Long:  "Remove abriment configs from config file",
	Run: func(cmd *cobra.Command, args []string) {

		if err := logout.RemoveAbrimentFromConfigfile(); err != nil {
			fmt.Println("logout failed.")
			return
		}

		fmt.Println("Logged out successfully!")
	},
}
