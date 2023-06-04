/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/eks"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/spf13/cobra"
	"github.com/surajincloud/kubectl-eks/pkg/kube"
)

// suggestionCmd represents the suggestion command
var suggestionCmd = &cobra.Command{
	Use:   "suggest-ami",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: suggestion,
}

func suggestion(cmd *cobra.Command, args []string) error {

	ctx := context.Background()

	// read flag values
	clusterName, _ := cmd.Flags().GetString("cluster-name")
	region, _ := cmd.Flags().GetString("region")

	// get Clustername
	clusterName, err := kube.GetClusterName(clusterName)
	if err != nil {
		log.Fatal(err)
	}
	// get region
	region, err = kube.GetRegion(region)
	if err != nil {
		log.Fatal(err)
	}
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err != nil {
		log.Fatal(err)
	}
	// getcluser info
	// Grab the cluster version
	client2 := eks.NewFromConfig(cfg)
	descluster, _ := client2.DescribeCluster(ctx, &eks.DescribeClusterInput{
		Name: aws.String(clusterName),
	})
	fmt.Printf("Cluster name: %v\n", clusterName)
	fmt.Printf("Cluster version: %v\n", aws.ToString(descluster.Cluster.Version))
	clusterVersion := aws.ToString(descluster.Cluster.Version)
	client3 := ssm.NewFromConfig(cfg)
	out, err := client3.GetParameter(ctx, &ssm.GetParameterInput{
		Name: aws.String(fmt.Sprintf("/aws/service/eks/optimized-ami/%s/amazon-linux-2/recommended/release_version", clusterVersion))})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Recommended Release version for EKS optimised amazon linux AMI: %v\n", aws.ToString(out.Parameter.Value))

	out, err = client3.GetParameter(ctx, &ssm.GetParameterInput{
		Name: aws.String(fmt.Sprintf("/aws/service/eks/optimized-ami/%s/amazon-linux-2/recommended/image_id", clusterVersion))})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Recommended Image ID for EKS optimised amazon linux AMI: %v\n", aws.ToString(out.Parameter.Value))

	return nil

}
func init() {
	rootCmd.AddCommand(suggestionCmd)
	suggestionCmd.PersistentFlags().String("cluster-name", "", "Cluster name")
	suggestionCmd.PersistentFlags().String("region", "", "region")
}
