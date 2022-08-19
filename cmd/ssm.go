/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/spf13/cobra"
	"github.com/surajincloud/kubectl-eks/pkg/kube"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/mmmorris1975/ssm-session-client/ssmclient"
)

// ssmCmd represents the ssm command
var ssmCmd = &cobra.Command{
	Use:   "ssm",
	Short: "Access given EKS node via SSM",
	Long: `SSM Access to given EKS Node
	IAM Roles needs to be attached to given EKS Node`,
	RunE: performSSM,
}

func init() {
	rootCmd.AddCommand(ssmCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// ssmCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// ssmCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func performSSM(cmd *cobra.Command, args []string) error {
	nodeList, err := kube.GetNodes(KubernetesConfigFlags)
	if err != nil {
		return err
	}
	var givenNode string
	for _, i := range nodeList {
		if i.Name == args[0] {
			// https://github.com/aws/containers-roadmap/issues/1395
			str := strings.Split(i.Spec.ProviderID, "/")
			givenNode = str[len(str)-1]
		}
	}
	fmt.Printf("SSM into node %v\n", givenNode)

	target := givenNode

	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	// A 3rd argument can be passed to specify a command to run before turning the shell over to the user
	log.Fatal(ssmclient.ShellSession(cfg, target))

	return nil

}
