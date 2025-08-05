package login

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/abrimentcloud/kubectl-abriment/config"
	"github.com/abrimentcloud/kubectl-abriment/response"

	"k8s.io/client-go/tools/clientcmd"
)

type LoginBody struct {
	UnsocpedToken string `json:"unscoped_token,omitempty"`
	Username      string `json:"username,omitempty"`
	Password      string `json:"password,omitempty"`
	Project       string `json:"project,omitempty"`
}

// send requests to login endpoint with provided credentials and returns the response.
func Login(body io.Reader, url string) (response.Response, error) {
	cli := http.Client{}

	req, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		return response.Response{}, fmt.Errorf("creating post request failed | %v", err)
	}
	req.Header.Add("Content-Type", "application/json")

	resp, err := cli.Do(req)
	if err != nil {
		return response.Response{}, fmt.Errorf("post request failed | %v", err)
	}
	defer resp.Body.Close()

	res := response.Response{}
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return response.Response{}, fmt.Errorf("post response decode failed | %v", err)
	}

	if !res.Success {
		return response.Response{}, errors.New(res.Message)
	}

	return res, nil
}

// getYamlConfig requests paas kubeconfig endpoint and returns the yaml file content.
func GetYamlConfig(token string, cfg *config.Config) ([]byte, error) {
	cli := http.Client{}

	req, err := http.NewRequest(http.MethodGet, cfg.ConfigEndpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("creating get request failed | %v", err)
	}
	req.Header.Add("X-Auth-Token", token)

	resp, err := cli.Do(req)
	if err != nil {
		return nil, fmt.Errorf("get request failed | %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("get response (200) decode failed | %v", err)
		}

		return body, nil
	}

	res := response.Response{}
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, fmt.Errorf("get response decode failed | %v", err)
	}

	return nil, errors.New(res.Message)
}

// saveYamlFile get kubeconfig file bytes and stores it either in existing custom path or default path.
func SaveConfigToConfigfile(yamlBytes []byte, dryrun bool) error {
	existingPath := os.Getenv("KUBECONFIG")

	existingPath, _ = strings.CutSuffix(existingPath, config.KubeconfigFileName)
	path := existingPath

	// If KUBECONFIG environment variable is not set, then set a default path.
	if existingPath == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("error getting user home directory | %v", err)
		}
		path = filepath.Join(home, ".kube")
	}

	if err := os.MkdirAll(path, 0755); err != nil {
		return fmt.Errorf("error creating directory for kubeconfig | %v", err)
	}

	filename := path

	// checks if KUBECONFIG custom path conatins the kubeconfig file name or not (It can be just directory and not the filename itself).
	if !strings.Contains(path, config.KubeconfigFileName) {
		filename = filepath.Join(path, config.KubeconfigFileName)
	}

	newkubeconfig, err := clientcmd.Load(yamlBytes)
	if err != nil {
		return fmt.Errorf("error loading new kubeconfig file | %v", err)
	}

	_, err = os.Stat(filename)
	if os.IsNotExist(err) {
		if err := clientcmd.WriteToFile(*newkubeconfig, filename); err != nil {
			return fmt.Errorf("error writing new kubeconfig file | %v", err)
		}
		return nil
	}

	existingkubeconfig, err := clientcmd.LoadFromFile(filename)
	if err != nil {
		return fmt.Errorf("error loading existing kubeconfig file | %v", err)
	}

	// add abriment-cluster to kubeconfig clusters.
	existingkubeconfig.Clusters[config.AbrimentCluster] = newkubeconfig.Clusters[config.AbrimentCluster]

	// add abriment-context to kubeconfig contexts.
	existingkubeconfig.Contexts[config.AbrimentContext] = newkubeconfig.Contexts[config.AbrimentContext]

	// add abriment-user to kubeconfig users.
	existingkubeconfig.AuthInfos[config.AbrimentUser] = newkubeconfig.AuthInfos[config.AbrimentUser]

	if dryrun {
		dryBytes, err := clientcmd.Write(*existingkubeconfig)
		if err != nil {
			return fmt.Errorf("error write configfile to bytes | %v", err)
		}

		fmt.Println("\nKubeconfig Preview:")
		fmt.Println("=" + strings.Repeat("=", 50))
		fmt.Println(string(dryBytes))
		fmt.Println("=" + strings.Repeat("=", 50))
		return nil
	}

	// writed edited configfile to disk.
	if err := clientcmd.WriteToFile(*existingkubeconfig, filename); err != nil {
		return fmt.Errorf("error modifying existing kubeconfig file | %v", err)
	}

	return nil
}
