package cmd

import (
	"context"
	"fmt"
	"log"
	"os"
	"text/tabwriter"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/spf13/cobra"
	"github.com/surajincloud/kubectl-eks/pkg/kube"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	amt "k8s.io/apimachinery/pkg/util/duration"
)

// nodesCmd represents the nodes command
var nodesCmd = &cobra.Command{
	Use:   "nodes",
	Short: "List all EKS Nodes",
	Long:  `A better way to list EKS nodes`,

	RunE: nodes,
}

func init() {
	rootCmd.AddCommand(nodesCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// nodesCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// nodesCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func nodes(cmd *cobra.Command, args []string) error {
	nodeList, err := kube.GetNodes(KubernetesConfigFlags)
	if err != nil {
		return err
	}

	ctx := context.Background()
	region := "eu-west-1"
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err != nil {
		log.Fatal(err)
	}
	// client := eks.NewFromConfig(cfg)
	client2 := ec2.NewFromConfig(cfg)

	w := tabwriter.NewWriter(os.Stdout, 5, 2, 3, ' ', tabwriter.TabIndent)
	defer w.Flush()
	fmt.Fprintln(w, "NAME", "\t", "INSTANCE-TYPE", "\t", "OS", "\t", "CAPACITY-TYPE", "\t", "REGION", "\t", "AMI-ID", "\t", "AMI-NAME", "\t", "AGE")
	for _, i := range nodeList {
		age := getAge(i.CreationTimestamp)
		img := i.Labels[kube.NodeGroupImage]
		dis, _ := client2.DescribeImages(ctx, &ec2.DescribeImagesInput{ImageIds: []string{img}})

		amiName := aws.ToString(dis.Images[0].Name)
		fmt.Fprintln(w, i.Name, "\t", i.Labels[kube.InstanceTypeLabel], "\t", i.Labels[kube.OsLabel], "\t", i.Labels[kube.CapacityTypeLabel], "\t", i.Labels[kube.ZoneLabel], "\t", i.Labels[kube.NodeGroupImage], "\t", amiName, "\t", age)
	}
	// SSM parameter
	//aws ssm get-parameter --name /aws/service/eks/optimized-ami/1.24/amazon-linux-2/recommended/image_id --region region-code --query "Parameter.Value" --output text
	client3 := ssm.NewFromConfig(cfg)
	out, err := client3.GetParameter(ctx, &ssm.GetParameterInput{
		Name: aws.String("/aws/service/eks/optimized-ami/1.23/amazon-linux-2/recommended/image_id"),
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(aws.ToString(out.Parameter.Value))
	return nil
}

func getAge(creationStamp metav1.Time) string {

	currentTime := time.Now()
	diff := currentTime.Sub(creationStamp.Time)
	return amt.HumanDuration(diff)
}
