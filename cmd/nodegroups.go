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

// nodegroupsCmd represents the nodegroups command
var nodegroupsCmd = &cobra.Command{
	Use:   "nodegroups",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: nodegroups,
}

func nodegroups(cmd *cobra.Command, args []string) error {

	ctx := context.Background()

	region, _ := cmd.Flags().GetString("region")

	if region == "" {
		fmt.Println("please pass region name with --region")
		os.Exit(0)
	}

	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err != nil {
		log.Fatal(err)
	}
	client := eks.NewFromConfig(cfg)

	nodegroupsList, err := client.ListNodegroups(ctx, &eks.ListNodegroupsInput{ClusterName: aws.String("staging-01")})
	if err != nil {
		log.Fatal(err)
	}
	w := tabwriter.NewWriter(os.Stdout, 5, 2, 3, ' ', tabwriter.TabIndent)
	defer w.Flush()
	fmt.Fprintln(w, "NAME", "\t", "RELEASE", "\t", "AMI_TYPE", "\t", "LC", "\t", "STATUS")
	for _, i := range nodegroupsList.Nodegroups {
		name := i
		dngp, err := client.DescribeNodegroup(ctx, &eks.DescribeNodegroupInput{ClusterName: aws.String("staging-01"), NodegroupName: aws.String(i)})
		if err != nil {
			log.Fatal(err)
		}

		rv := aws.ToString(dngp.Nodegroup.ReleaseVersion)

		status := dngp.Nodegroup.Status

		amiType := dngp.Nodegroup.AmiType

		ltid := aws.ToString(dngp.Nodegroup.LaunchTemplate.Id)

		fmt.Fprintln(w, name, "\t", rv, "\t", amiType, "\t", ltid, "\t", status)
	}
	return nil
}

func init() {
	rootCmd.AddCommand(nodegroupsCmd)
	addonsCmd.PersistentFlags().String("region", "", "region")
}
