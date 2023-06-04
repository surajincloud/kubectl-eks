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
	"github.com/surajincloud/kubectl-eks/pkg/kube"
)

// addonsCmd represents the addons command
var addonsCmd = &cobra.Command{
	Use:   "addons",
	Short: "List Addons with current and recommended versions",
	Long:  "List Addons with current and recommended versions",
	RunE:  addons,
}

type Addons struct {
	Name           string
	CurrentVersion string
	LatestVersion  string
}

func addons(cmd *cobra.Command, args []string) error {

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
	clusterVersion := aws.ToString(descluster.Cluster.Version)
	// List addons
	resp, err := eksSvc.ListAddons(context.TODO(), &eks.ListAddonsInput{
		ClusterName: &clusterName,
	})
	if err != nil {
		log.Fatal(err)
	}

	w := tabwriter.NewWriter(os.Stdout, 5, 2, 3, ' ', tabwriter.TabIndent)
	defer w.Flush()
	fmt.Fprintln(w, "NAME", "\t", "CURRENT-VERSION", "\t", "LATEST")
	for _, i := range resp.Addons {
		re, _ := eksSvc.DescribeAddon(ctx, &eks.DescribeAddonInput{
			AddonName:   &i,
			ClusterName: &clusterName,
		})
		currentVersion := aws.ToString(re.Addon.AddonVersion)

		resp, err := eksSvc.DescribeAddonVersions(context.TODO(), &eks.DescribeAddonVersionsInput{
			AddonName:         &i,
			KubernetesVersion: &clusterVersion,
		})
		if err != nil {
			log.Fatal(err)
		}

		fmt.Fprintln(w, i, "\t", currentVersion, "\t", aws.ToString(resp.Addons[0].AddonVersions[0].AddonVersion))

	}
	return nil
}

func init() {
	rootCmd.AddCommand(addonsCmd)
	addonsCmd.PersistentFlags().String("cluster-name", "", "Cluster name")
	addonsCmd.PersistentFlags().String("region", "", "region")
}
