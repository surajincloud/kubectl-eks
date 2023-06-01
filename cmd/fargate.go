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
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/eks"
	"github.com/spf13/cobra"
)

// fargateCmd represents the fargate command
var fargateCmd = &cobra.Command{
	Use:   "fargate",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: fargate,
}

func fargate(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	// read flag values
	clusterName, _ := cmd.Flags().GetString("cluster-name")
	region, _ := cmd.Flags().GetString("region")

	if clusterName == "" {
		fmt.Println("please pass cluster name with --cluster-name")
		os.Exit(0)
	}
	if region == "" {
		fmt.Println("please pass region name with --region")
		os.Exit(0)
	}

	// Load AWS configuration from environment variables or default configuration files
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err != nil {
		fmt.Println("Failed to load AWS configuration:", err)
		os.Exit(1)
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
		fmt.Fprintln(w, profile, "\t", out.FargateProfile.FargateProfileArn, "\t", out.FargateProfile.FargateProfileArn)

	}
	return nil
}

func init() {
	rootCmd.AddCommand(fargateCmd)
	fargateCmd.PersistentFlags().String("cluster-name", "", "Cluster name")
	fargateCmd.PersistentFlags().String("region", "", "region")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// fargateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// fargateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
