package logout

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/abrimentcloud/kubectl-abriment/config"

	"k8s.io/client-go/tools/clientcmd"
)

// saveYamlFile get kubeconfig file bytes and stores it either in existing custom path or default path.
func RemoveAbrimentFromConfigfile() error {
	path := os.Getenv("KUBECONFIG")

	// If KUBECONFIG environment variable is not set, then set the default path.
	if path == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("error getting user home directory | %v", err)
		}
		path = filepath.Join(home, ".kube", config.KubeconfigFileName)
	}

	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return nil
	}

	kubeconfig, err := clientcmd.LoadFromFile(path)
	if err != nil {
		return fmt.Errorf("error loading existing kubeconfig file | %v", err)
	}

	// check for abriment cluster exsist in the kubeconfig file on machine. if not exists then add the cluster to clusters map.
	delete(kubeconfig.Clusters, config.AbrimentCluster)

	// check for abriment context exsist in the kubeconfig file on machine. if not exists then add the context to context map.
	delete(kubeconfig.Contexts, config.AbrimentContext)

	// check for abriment users exsist in the kubeconfig file on machine. if not exists then add the context to users map.
	delete(kubeconfig.AuthInfos, config.AbrimentUser)

	// writed edited configfile to disk.
	if err := clientcmd.WriteToFile(*kubeconfig, path); err != nil {
		return fmt.Errorf("error modifying existing kubeconfig file | %v", err)
	}

	return nil
}
