/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/eks"
	"github.com/spf13/cobra"
	"github.com/surajincloud/kubectl-eks/pkg/kube"
)

// nodegroupsCmd represents the nodegroups command
var nodegroupsCmd = &cobra.Command{
	Use:   "nodegroups",
	Short: "List EKS Nodegroups",
	Long:  "List EKS Nodegroups",
	RunE:  nodegroups,
}

func nodegroups(cmd *cobra.Command, args []string) error {

	// AmiTypesMap := map[string]string{
	// 	"AL2_x86_64": "Amazon Linux",
	// }

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
	client := eks.NewFromConfig(cfg)
	nodegroupsList, err := client.ListNodegroups(ctx, &eks.ListNodegroupsInput{ClusterName: aws.String(clusterName)})
	if err != nil {
		log.Fatal(err)
	}
	w := tabwriter.NewWriter(os.Stdout, 5, 2, 3, ' ', tabwriter.TabIndent)
	defer w.Flush()
	fmt.Fprintln(w, "NAME", "\t", "RELEASE", "\t", "AMI_TYPE", "\t", "INSTANCE_TYPES", "\t", "STATUS")
	for _, i := range nodegroupsList.Nodegroups {
		name := i
		dngp, err := client.DescribeNodegroup(ctx, &eks.DescribeNodegroupInput{ClusterName: aws.String(clusterName), NodegroupName: aws.String(i)})
		if err != nil {
			log.Fatal(err)
		}

		rv := aws.ToString(dngp.Nodegroup.ReleaseVersion)

		status := dngp.Nodegroup.Status

		amiType := dngp.Nodegroup.AmiType

		instanceTypes := strings.Join(dngp.Nodegroup.InstanceTypes, ",")
		fmt.Fprintln(w, name, "\t", rv, "\t", amiType, "\t", instanceTypes, "\t", status)
	}
	return nil
}

func init() {
	rootCmd.AddCommand(nodegroupsCmd)
	nodegroupsCmd.PersistentFlags().String("cluster-name", "", "Cluster name")
	nodegroupsCmd.PersistentFlags().String("region", "", "region")
}
