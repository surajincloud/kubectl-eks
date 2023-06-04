package cmd

import (
	"context"
	"fmt"
	"log"
	"os"
	"text/tabwriter"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/spf13/cobra"
	"github.com/surajincloud/kubectl-eks/pkg/kube"
)

// nodesCmd represents the nodes command
var nodesCmd = &cobra.Command{
	Use:   "nodes",
	Short: "List EKS Nodes",
	Long:  `A better way to list EKS nodes`,

	RunE: nodes,
}

func nodes(cmd *cobra.Command, args []string) error {
	nodeList, err := kube.GetNodes(KubernetesConfigFlags)
	if err != nil {
		return err
	}

	ctx := context.Background()
	// read flag values
	clusterName, _ := cmd.Flags().GetString("cluster-name")
	region, _ := cmd.Flags().GetString("region")

	// get Clustername
	clusterName, err = kube.GetClusterName(clusterName)
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
	ec2Client := ec2.NewFromConfig(cfg)

	w := tabwriter.NewWriter(os.Stdout, 5, 2, 3, ' ', tabwriter.TabIndent)
	defer w.Flush()
	fmt.Fprintln(w, "NAME", "\t", "INSTANCE-TYPE", "\t", "OS", "\t", "CAPACITY-TYPE", "\t", "REGION", "\t", "AMI-ID", "\t", "AMI-NAME", "\t", "AGE")
	for _, i := range nodeList {
		age := kube.GetAge(i.CreationTimestamp)
		img := i.Labels[kube.NodeGroupImage]
		dis, _ := ec2Client.DescribeImages(ctx, &ec2.DescribeImagesInput{ImageIds: []string{img}})

		amiName := aws.ToString(dis.Images[0].Name)
		fmt.Fprintln(w, i.Name, "\t", i.Labels[kube.InstanceTypeLabel], "\t", i.Labels[kube.OsLabel], "\t", i.Labels[kube.CapacityTypeLabel], "\t", i.Labels[kube.ZoneLabel], "\t", i.Labels[kube.NodeGroupImage], "\t", amiName, "\t", age)
	}
	return nil
}

func init() {
	rootCmd.AddCommand(nodesCmd)
}
