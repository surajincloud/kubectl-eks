package kube

import (
	"fmt"
	"os"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/apimachinery/pkg/util/duration"
)

func GetClusterName(clusterName string) (string, error) {
	if clusterName == "" {
		clusterName = os.Getenv("AWS_EKS_CLUSTER")
		if clusterName == "" {
			clusterName, err := getClusterFromKubeconfig()
			if err != nil {
				return "", err
			}
			return clusterName, nil
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

func GetAge(creationStamp metav1.Time) string {

	currentTime := time.Now()
	diff := currentTime.Sub(creationStamp.Time)
	return duration.HumanDuration(diff)
}
