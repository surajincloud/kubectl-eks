/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/spf13/cobra"
	awspkg "github.com/surajincloud/kubectl-eks/pkg/aws"
	"github.com/surajincloud/kubectl-eks/pkg/kube"
)

// logsCmd represents the logs command
var logsCmd = &cobra.Command{
	Use:   "logs",
	Args:  cobra.ExactArgs(1),
	Short: "Get logs from an EKS cluster control plane or nodes",
	Long:  "Allows you to see logs from different components of the EKS control plane or from nodes",
	RunE:  logs,
}

func logs(cmd *cobra.Command, args []string) error {

	// pass empty string to let fuction get name
	clusterName, err := kube.GetClusterName(*KubernetesConfigFlags.ClusterName)
	if err != nil {
		return err
	}

	logTarget := args[0]

	cloudwatchLogStreams := []string{
		"scheduler",
		"kube-scheduler",
		"audit",
		"kube-audit",
		"kube-apiserver-audit",
		"cm",
		"controller-manager",
		"kube-controller-manager",
		"api",
		"apiserver",
		"kube-apiserver",
		"auth",
		"authenticator",
		"ccm",
		"cloud-controller",
		"cloud-controller-manager",
	}

	if contains(cloudwatchLogStreams, logTarget) {
		cwlGroupPrefix := getLogStreamPrefix(logTarget)

		cwlGroupName := "/aws/eks/" + clusterName + "/cluster"
		var cwlStreamInput cloudwatchlogs.DescribeLogStreamsInput
		cwlStreamInput.LogGroupName = &cwlGroupName
		cwlStreamInput.LogStreamNamePrefix = &cwlGroupPrefix

		// aws config
		ctx := context.Background()
		// read flag values
		region, _ := cmd.Flags().GetString("region")

		cfg, err := awspkg.GetAWSConfig(ctx, region)
		if err != nil {
			log.Fatal(err)
		}

		cwl := cloudwatchlogs.NewFromConfig(cfg)
		// verify the group exists first
		err = ensureLogGroupExists(cwlGroupName, ctx, cwl)
		//TODO prompt if user wants to enable logs
		if err != nil {
			panic(err)
		}

		var newestStream string
		var limit int32 = 100
		streams, err := cwl.DescribeLogStreams(ctx, &cwlStreamInput)
		if err != nil {
			log.Fatal(err)
		}
		// for _, stream := range streams.LogStreams {
		// 	fmt.Printf("%s\n", *stream.LogStreamName)
		// }
		newestStream = *streams.LogStreams[0].LogStreamName

		resp, err := getLogEvents(&cwlGroupName, &newestStream, &limit, ctx, cwl)
		if err != nil {
			fmt.Println("Got error getting log events:")
			return err
		}

		// gotToken := ""
		// nextToken := ""

		for _, event := range resp.Events {
			// gotToken = nextToken
			// nextToken = *resp.NextForwardToken

			// if gotToken == nextToken {
			// 	break
			// }

			fmt.Println("  ", *event.Message)
		}
	} else {
		fmt.Println("Logs not found")
		return nil
	}
	return nil
}

// ensureLogGroupExists first checks if the log group exists
func ensureLogGroupExists(name string, ctx context.Context, cwl *cloudwatchlogs.Client) error {
	resp, err := cwl.DescribeLogGroups(ctx, &cloudwatchlogs.DescribeLogGroupsInput{})
	if err != nil {
		return err
	}

	for _, logGroup := range resp.LogGroups {
		if *logGroup.LogGroupName == name {
			return nil
		}
	}

	return err
}

func getLogEvents(logGroupName *string, logStreamName *string, limit *int32, ctx context.Context, cwl *cloudwatchlogs.Client) (*cloudwatchlogs.GetLogEventsOutput, error) {

	resp, err := cwl.GetLogEvents(ctx, &cloudwatchlogs.GetLogEventsInput{
		Limit:         limit,
		LogGroupName:  logGroupName,
		LogStreamName: logStreamName,
	})
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func init() {
	rootCmd.AddCommand(logsCmd)
	logsCmd.PersistentFlags().String("region", "", "region")
}

// contains checks if a string is present in a slice
func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

func getLogStreamPrefix(logTarget string) string {
	schedulerSlice := []string{"scheduler", "kube-scheduler"}
	auditSlice := []string{"audit", "kube-audit", "kube-apiserver-audit"}
	cmSlice := []string{"cm", "controller-manager", "kube-controller-manager"}
	apiSlice := []string{"api", "apiserver", "kube-apiserver"}
	authSlice := []string{"auth", "authenticator"}
	ccmSlice := []string{"ccm", "cloud-controller", "cloud-controller-manager"}

	if contains(schedulerSlice, logTarget) {
		return "kube-scheduler"
	} else if contains(auditSlice, logTarget) {
		return "kube-apiserver-audit"
	} else if contains(cmSlice, logTarget) {
		return "kube-controller-manager"
	} else if contains(apiSlice, logTarget) {
		return "kube-apiserver"
	} else if contains(authSlice, logTarget) {
		return "authenticator"
	} else if contains(ccmSlice, logTarget) {
		return "cloud-controller-manager"
	} else {
		return ""
	}

}
