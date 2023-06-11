package cmd

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/eks"
	kubeconfig "github.com/siderolabs/go-kubeconfig"
	"github.com/spf13/cobra"
	awspkg "github.com/surajincloud/kubectl-eks/pkg/aws"
	"github.com/surajincloud/kubectl-eks/pkg/kube"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"
)

// kubeconfigCmd represents the kubeconfig command
var kubeconfigCmd = &cobra.Command{
	Use:   "kubeconfig",
	Short: "Get Kubeconfig for given cluster",
	Long:  "Get Kubeconfig for given cluster",
	RunE:  kubeconfigCommand,
}

// Constants for kubeconfig
const (
	kubeconfigFilePath = "kubeconfig"
)

func kubeconfigCommand(cmd *cobra.Command, args []string) error {

	ctx := context.Background()

	// read flag values
	clusterName, _ := cmd.Flags().GetString("cluster-name")
	out, _ := cmd.Flags().GetBool("out")
	region, _ := cmd.Flags().GetString("region")

	// get Clustername
	clusterName, err := kube.GetClusterName(clusterName)
	if err != nil {
		log.Fatal(err)
	}
	// aws config
	cfg, err := awspkg.GetAWSConfig(ctx, region)
	if err != nil {
		log.Fatal(err)
	}

	// define eks service
	eksSvc := eks.NewFromConfig(cfg)
	descluster, _ := eksSvc.DescribeCluster(ctx, &eks.DescribeClusterInput{
		Name: aws.String(clusterName),
	})

	// Extract necessary cluster details
	serverURL := descluster.Cluster.Endpoint

	// cluser certificate
	certData := descluster.Cluster.CertificateAuthority.Data
	decodedCertData, err := base64.StdEncoding.DecodeString(*certData)
	if err != nil {
		fmt.Printf("Failed to decode certificate authority data: %v\n", err)
		os.Exit(1)
	}
	clusterCAData := string(decodedCertData)
	// username
	userName := clusterName + "-user"
	// build kubeconfig
	config := &api.Config{
		APIVersion: "v1",
		Kind:       "Config",
		Clusters: map[string]*api.Cluster{
			clusterName: {
				Server:                   *serverURL,
				CertificateAuthorityData: []byte(clusterCAData),
			},
		},
		Contexts: map[string]*api.Context{
			clusterName: {
				Cluster:   clusterName,
				AuthInfo:  userName,
				Namespace: "default",
			},
		},
		AuthInfos: map[string]*api.AuthInfo{
			userName: {
				Exec: &api.ExecConfig{
					APIVersion:      "client.authentication.k8s.io/v1beta1",
					Command:         "aws",
					Args:            []string{"eks", "get-token", "--cluster-name", clusterName, "--region", region},
					InteractiveMode: api.IfAvailableExecInteractiveMode,
				},
			},
		},
		CurrentContext: clusterName,
	}

	if !out {
		existingPath, _ := kubeconfig.DefaultPath()
		a, err := kubeconfig.Load(existingPath)
		if err != nil {
			fmt.Println("error")
		}

		err = a.Merge(config, kubeconfig.MergeOptions{ActivateContext: true})
		if err != nil {
			fmt.Println("error merging the kubeconfig")
		}

		err = a.Write(existingPath)
		if err != nil {
			fmt.Println("error writing the kubeconfig")
		}
		fmt.Println("existing kubeconfig file is merged")
	} else {
		kc, err := clientcmd.Write(*config)
		// err = clientcmd.WriteToFile(*config, kubeconfigFilePath)
		if err != nil {
			fmt.Printf("Failed to write kubeconfig file: %v\n", err)
			os.Exit(1)
		}
		fmt.Println(string(kc))
		// fmt.Printf("kubeconfig file created at %s\n", kubeconfigFilePath)
	}
	return nil
}
func init() {
	rootCmd.AddCommand(kubeconfigCmd)
	kubeconfigCmd.PersistentFlags().String("cluster-name", "", "Cluster name")
	kubeconfigCmd.PersistentFlags().String("region", "", "region")
	kubeconfigCmd.PersistentFlags().Bool("out", false, "Print kubeconfig")
}
