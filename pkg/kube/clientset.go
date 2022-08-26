package kube

import (
	"fmt"

	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/kubernetes"
)

// ClientSet k8s clientset
func ClientSet(configFlags *genericclioptions.ConfigFlags) (*kubernetes.Clientset, string) {
	namespace := ""
	if configFlags.Namespace != nil {
		namespace = *configFlags.Namespace
	}
	config, err := configFlags.ToRESTConfig()
	if err != nil {
		panic("kube config load error")
	}
	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		fmt.Println("Error generating Kubernetes configuration", err)
	}
	return clientSet, namespace
}
