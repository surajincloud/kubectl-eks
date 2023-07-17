package cmd

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/spf13/cobra"
	"github.com/surajincloud/kubectl-eks/pkg/kube"

	"github.com/mmmorris1975/ssm-session-client/ssmclient"
	awspkg "github.com/surajincloud/kubectl-eks/pkg/aws"
)

// ssmCmd represents the ssm command
var ssmCmd = &cobra.Command{
	Use:   "ssm [flags] node",
	Args:  cobra.ExactArgs(1),
	Short: "Access given EKS node via SSM",
	Long: `
	SSM Access to given EKS Node
	IAM Roles needs to be attached to given EKS Node
	Check docs: https://surajincloud.github.io/kubectl-eks/usage/#access-to-eks-node-via-ssm`,
	RunE: performSSM,
}

func performSSM(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	nodeList, err := kube.GetNodes(KubernetesConfigFlags)
	if err != nil {
		return err
	}
	var givenNode, region string
	for _, i := range nodeList {
		if i.Name == args[0] {
			// https://github.com/aws/containers-roadmap/issues/1395
			str := strings.Split(i.Spec.ProviderID, "/")
			givenNode = str[len(str)-1]
			region = i.Labels["topology.kubernetes.io/region"]
		}
	}
	fmt.Printf("SSM into node %v in region %v\n", givenNode, region)

	target := givenNode

	// aws config
	cfg, err := awspkg.GetAWSConfig(ctx, region)
	if err != nil {
		log.Fatal(err)
	}

	// A 3rd argument can be passed to specify a command to run before turning the shell over to the user
	log.Fatal(ssmclient.ShellPluginSession(cfg, target))
	return nil

}

func init() {
	rootCmd.AddCommand(ssmCmd)
}
