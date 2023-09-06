package cmd

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/spf13/cobra"
	awspkg "github.com/surajincloud/kubectl-eks/pkg/aws"
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
	region, _ := cmd.Flags().GetString("region")

	// aws config
	cfg, err := awspkg.GetAWSConfig(ctx, region)
	if err != nil {
		log.Fatal(err)
	}
	ec2Client := ec2.NewFromConfig(cfg)

	w := tabwriter.NewWriter(os.Stdout, 5, 2, 3, ' ', tabwriter.TabIndent)
	defer w.Flush()
	fmt.Fprintln(w, "NAME", "\t", "INSTANCE-TYPE", "\t", "OS", "\t", "CAPACITY-TYPE", "\t", "ZONE", "\t", "AMI-ID", "\t", "AMI-NAME", "\t", "AGE")
	for _, i := range nodeList {
		var amiID, amiName, capacityType string
		age := kube.GetAge(i.CreationTimestamp)

		// AMI
		if i.Labels[kube.NodeGroupImage] == "" {
			amiID = i.Labels[kube.KarpenterImage]
			if amiID != "" {
				dis, err := ec2Client.DescribeImages(ctx, &ec2.DescribeImagesInput{ImageIds: []string{amiID}})
				if err != nil {
					log.Fatal(err)
				}
				amiName = aws.ToString(dis.Images[0].Name)
			}

		} else {
			amiID = i.Labels[kube.NodeGroupImage]
			dis, err := ec2Client.DescribeImages(ctx, &ec2.DescribeImagesInput{ImageIds: []string{amiID}})
			if err != nil {
				log.Fatal(err)
			}
			amiName = aws.ToString(dis.Images[0].Name)
		}

		// Capacity Type
		if i.Labels[kube.CapacityTypeLabel] != "" {
			capacityType = i.Labels[kube.CapacityTypeLabel]
		} else {
			capacityType = strings.ToUpper(i.Labels[kube.KarpenterCapacityTypeLabel])

		}
		fmt.Fprintln(w, i.Name, "\t", i.Labels[kube.InstanceTypeLabel], "\t", i.Labels[kube.OsLabel], "\t", capacityType, "\t", i.Labels[kube.ZoneLabel], "\t", amiID, "\t", amiName, "\t", age)
	}
	return nil
}

func init() {
	rootCmd.AddCommand(nodesCmd)
	nodesCmd.PersistentFlags().String("region", "", "region")
}
