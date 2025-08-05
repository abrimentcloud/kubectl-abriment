package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/abrimentcloud/kubectl-abriment/config"
	"github.com/abrimentcloud/kubectl-abriment/pkg/login"

	"github.com/spf13/cobra"
)

var (
	user, pass, token, dryrun string
)

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Get paas kubeconfig file",
	Long:  "Login to the provider and get config file after login",
	Run: func(cmd *cobra.Command, args []string) {

		if token == "" && (user == "" || pass == "") {
			fmt.Println("You should either provide token or username and password")
			fmt.Println("Run \"kubectl login help\" for more information.")
			return
		}

		cfg, err := config.ParseCfg()
		if err != nil {
			fmt.Println(err)
			return
		}

		var body login.LoginBody
		if token != "" {
			body = login.LoginBody{
				UnsocpedToken: token,
			}
		}

		if user != "" {
			body = login.LoginBody{
				Username: user,
				Password: pass,
			}
		}

		bb, err := json.Marshal(body)
		if err != nil {
			fmt.Println("Invalid username and password!")
			return
		}
		bodyBytes := bytes.NewReader(bb)
		res, err := login.Login(bodyBytes, cfg.LoginEndpoint)
		if err != nil {
			fmt.Println("Invalid token!")
			return
		}

		yamlBytes, err := login.GetYamlConfig(res.Data.Token.ID, cfg)
		if err != nil {
			fmt.Println("Error getting config file!")
			return
		}

		var dryRunChoice bool
		if dryrun == "client" {
			dryRunChoice = true
		}

		if err := login.SaveConfigToConfigfile(yamlBytes, dryRunChoice); err != nil {
			fmt.Println(err)
			return
		}

		if dryRunChoice {
			return
		}

		fmt.Println("Configuration saved successfully!")
	},
}
