package cmd

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/eks"
	kubeconfig "github.com/siderolabs/go-kubeconfig"
	"github.com/spf13/cobra"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"
)

// kubeconfigCmd represents the kubeconfig command
var kubeconfigCmd = &cobra.Command{
	Use:   "kubeconfig",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: kubeconfigCommand,
}

// Constants for kubeconfig
const (
	kubeconfigFilePath = "kubeconfig"
)

func kubeconfigCommand(cmd *cobra.Command, args []string) error {

	ctx := context.Background()

	// read flag values
	clusterName, _ := cmd.Flags().GetString("cluster-name")
	region, _ := cmd.Flags().GetString("region")
	merge, _ := cmd.Flags().GetBool("merge")

	if clusterName == "" {
		fmt.Println("please pass cluster name with --cluster-name")
		os.Exit(0)
	}
	if region == "" {
		fmt.Println("please pass region name with --region")
		os.Exit(0)
	}

	// aws config
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
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

	if merge {
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

		err = clientcmd.WriteToFile(*config, kubeconfigFilePath)
		if err != nil {
			fmt.Printf("Failed to write kubeconfig file: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("kubeconfig file created at %s\n", kubeconfigFilePath)
	}
	return nil
}
func init() {
	rootCmd.AddCommand(kubeconfigCmd)
	kubeconfigCmd.PersistentFlags().String("cluster-name", "", "Cluster name")
	kubeconfigCmd.PersistentFlags().String("region", "", "region")
	kubeconfigCmd.PersistentFlags().Bool("merge", false, "Merge into existing kubeconfig")
}
