package kube

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"

	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"
)

func getClusterFromKubeconfig() (string, error) {
	config, err := ReadKubeconfig()
	if err != nil {
		return "", err
	}
	// Getting the current context
	currentContext := config.CurrentContext
	context := config.Contexts[currentContext]

	return context.Cluster, nil
}

func ReadKubeconfig() (api.Config, error) {
	// Check if KUBECONFIG environment variable is set
	kubeconfigPath := os.Getenv("KUBECONFIG")

	// If KUBECONFIG is not set, use the default path
	if kubeconfigPath == "" {
		usr, err := user.Current()
		if err != nil {
			fmt.Printf("Failed to get user information: %v\n", err)
			return api.Config{}, err
		}
		kubeconfigPath = filepath.Join(usr.HomeDir, ".kube", "config")
	}

	// Loading the kubeconfig file
	config, err := clientcmd.LoadFromFile(kubeconfigPath)
	if err != nil {
		fmt.Printf("Failed to load kubeconfig: %v\n", err)
		return api.Config{}, err
	}
	return *config, nil
}
