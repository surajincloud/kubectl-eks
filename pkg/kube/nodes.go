package kube

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

func GetNodes(configFlags *genericclioptions.ConfigFlags) ([]corev1.Node, error) {
	clientSet, _ := ClientSet(configFlags)
	nodeList, err := clientSet.CoreV1().Nodes().List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return []corev1.Node{}, err
	}
	return nodeList.Items, nil
}

func GetSA(configFlags *genericclioptions.ConfigFlags) ([]corev1.ServiceAccount, error) {
	clientSet, namespace := ClientSet(configFlags)
	nodeList, err := clientSet.CoreV1().ServiceAccounts(namespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return []corev1.ServiceAccount{}, err
	}
	return nodeList.Items, nil
}
