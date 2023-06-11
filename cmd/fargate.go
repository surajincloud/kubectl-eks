/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"log"
	"os"
	"text/tabwriter"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/eks"
	"github.com/spf13/cobra"
	awspkg "github.com/surajincloud/kubectl-eks/pkg/aws"
	"github.com/surajincloud/kubectl-eks/pkg/kube"
)

// fargateCmd represents the fargate command
var fargateCmd = &cobra.Command{
	Use:   "fargate",
	Short: "List fargate profiles",
	Long:  "List fargate profiles",
	RunE:  fargate,
}

func fargate(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	// read flag values
	region, _ := cmd.Flags().GetString("region")

	// get Clustername
	clusterName, err := kube.GetClusterName(*KubernetesConfigFlags.ClusterName)
	if err != nil {
		log.Fatal(err)
	}

	// aws config
	cfg, err := awspkg.GetAWSConfig(ctx, region)
	if err != nil {
		log.Fatal(err)
	}

	// Create an EKS client using the loaded configuration
	client := eks.NewFromConfig(cfg)

	// Retrieve a list of Fargate profiles
	input := &eks.ListFargateProfilesInput{
		ClusterName: &clusterName,
	}
	output, err := client.ListFargateProfiles(ctx, input)
	if err != nil {
		fmt.Println("Failed to list EKS Fargate profiles:", err)
		os.Exit(1)
	}

	w := tabwriter.NewWriter(os.Stdout, 5, 2, 3, ' ', tabwriter.TabIndent)
	defer w.Flush()
	fmt.Fprintln(w, "NAME", "\t", "PROFILE_ARN", "\t", "STATUS")
	for _, profile := range output.FargateProfileNames {
		out, err := client.DescribeFargateProfile(ctx, &eks.DescribeFargateProfileInput{
			ClusterName:        aws.String(clusterName),
			FargateProfileName: aws.String(profile),
		})
		if err != nil {
			log.Fatal(err)
		}
		fmt.Fprintln(w, profile, "\t", aws.ToString(out.FargateProfile.FargateProfileArn), "\t", aws.ToString(out.FargateProfile.FargateProfileArn))

	}
	return nil
}

func init() {
	rootCmd.AddCommand(fargateCmd)
	fargateCmd.PersistentFlags().String("region", "", "region")
}
