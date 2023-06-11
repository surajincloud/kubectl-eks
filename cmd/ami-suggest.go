/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/eks"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/spf13/cobra"
	awspkg "github.com/surajincloud/kubectl-eks/pkg/aws"
	"github.com/surajincloud/kubectl-eks/pkg/kube"
)

// suggestionCmd represents the suggestion command
var suggestionCmd = &cobra.Command{
	Use:   "suggest-ami",
	Short: "Suggest recommended version of AMI",
	Long:  "Suggest recommended version of AMI",
	RunE:  suggestion,
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
	// aws config
	cfg, err := awspkg.GetAWSConfig(ctx, region)
	if err != nil {
		log.Fatal(err)
	}
	// getcluser info
	// Grab the cluster version
	client2 := eks.NewFromConfig(cfg)
	descluster, _ := client2.DescribeCluster(ctx, &eks.DescribeClusterInput{
		Name: aws.String(clusterName),
	})
	clusterVersion := aws.ToString(descluster.Cluster.Version)
	client3 := ssm.NewFromConfig(cfg)
	fmt.Println("Recommended versions for:")
	fmt.Println("Amazon Linux 2:")
	out, err := client3.GetParameter(ctx, &ssm.GetParameterInput{
		Name: aws.String(fmt.Sprintf("/aws/service/eks/optimized-ami/%s/amazon-linux-2/recommended/release_version", clusterVersion))})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Release version: %v\n", aws.ToString(out.Parameter.Value))

	out, err = client3.GetParameter(ctx, &ssm.GetParameterInput{
		Name: aws.String(fmt.Sprintf("/aws/service/eks/optimized-ami/%s/amazon-linux-2/recommended/image_id", clusterVersion))})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Image ID: %v\n", aws.ToString(out.Parameter.Value))
	fmt.Println("Bottlerocket:")
	out, err = client3.GetParameter(ctx, &ssm.GetParameterInput{
		Name: aws.String(fmt.Sprintf("/aws/service/bottlerocket/aws-k8s-%s/x86_64/latest/image_id", clusterVersion))})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Image ID: %v\n", aws.ToString(out.Parameter.Value))

	return nil

}
func init() {
	rootCmd.AddCommand(suggestionCmd)
	suggestionCmd.PersistentFlags().String("cluster-name", "", "Cluster name")
	suggestionCmd.PersistentFlags().String("region", "", "region")
}
