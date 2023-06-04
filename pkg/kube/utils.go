package kube

import (
	"fmt"
	"os"
)

func GetClusterName(clusterName string) (string, error) {
	if clusterName == "" {
		clusterName = os.Getenv("AWS_EKS_CLUSTER")
		if clusterName == "" {
			return "", fmt.Errorf("please pass cluster name with --cluster-name or with AWS_EKS_CLUSTER environment variable")
		}
		return clusterName, nil
	}
	return clusterName, nil
}

func GetRegion(region string) (string, error) {
	if region == "" {
		region = os.Getenv("AWS_REGION")
		if region == "" {
			return "", fmt.Errorf("please pass region name with --region or with AWS_REGION environment variable")
		}
		return region, nil
	}
	return region, nil
}
