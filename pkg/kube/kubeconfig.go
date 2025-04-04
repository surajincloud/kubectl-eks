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
	loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
	configOverrides := &clientcmd.ConfigOverrides{}

	clientConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, configOverrides)

	rawConfig, err := clientConfig.RawConfig()
	if err != nil {
		return api.Config{}, err
	}
	return rawConfig, nil
}
