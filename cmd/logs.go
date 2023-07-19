/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"log"
	"os/exec"
	"strings"

	"github.com/antchfx/jsonquery"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/spf13/cobra"
	awspkg "github.com/surajincloud/kubectl-eks/pkg/aws"
	"github.com/surajincloud/kubectl-eks/pkg/kube"
)

// logsCmd represents the logs command
var logsCmd = &cobra.Command{
	Use:               "logs [flags] LOG_SOURCE",
	ValidArgsFunction: validateArgs,
	Example: `  kubectl eks logs kube-apiserver
  kubectl eks logs NODE [query]`,
	ValidArgs: []string{
		"kube-scheduler",
		"kube-apiserver",
		"kube-apiserver-audit",
		"kube-controller-manager",
		"authenticator",
		"cloud-controller-manager"},
	ArgAliases: []string{
		"scheduler",
		"audit",
		"kube-audit",
		"cm",
		"controller-manager",
		"api",
		"apiserver",
		"auth",
		"ccm",
		"cloud-controller"},
	Args:  cobra.MinimumNArgs(1),
	Short: "Get logs from EKS control plane or nodes",
	Long:  "Get logs from EKS control plane or nodes",
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
			// TODO allow for following tokens for more logs
			// TODO allow -f follow for logs
			// gotToken = nextToken
			// nextToken = *resp.NextForwardToken

			// if gotToken == nextToken {
			// 	break
			// }

			fmt.Println("  ", *event.Message)
		}
	} else {
		// we need to assume the target is a node instead of control plane
		nodeList, err := kube.GetNodes(KubernetesConfigFlags)
		if err != nil {
			return err
		}

		var nodeMatched bool = false
		var currentNodeSlice []string
		logTargetSlice := strings.Split(logTarget, ".")

		for _, i := range nodeList {
			// match node based on substring
			currentNodeSlice = strings.Split(i.Name, ".")
			if currentNodeSlice[0] == logTargetSlice[0] {
				nodeMatched = true
				var query string
				// create URL for log fetching
				rawURL := "/api/v1/nodes/" + i.Name + "/proxy/logs/?query="
				if len(args) > 1 {
					query = args[1]
				} else {
					query = "kubelet"
				}
				// validate kubelet settings for remote logs
				if validateKubeletConfig(i.Name) {
					kubeLogsCmdOutput, err := exec.Command("kubectl", "get", "--raw", rawURL+query).Output()
					if err != nil {
						log.Fatal(err)
					}
					fmt.Printf("%s\n", kubeLogsCmdOutput)
				}
			}
		}
		if nodeMatched {
			return nil
		} else {
			fmt.Printf("Node %s not found\n", logTarget)
		}
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

// Convert possible log source aliases to full log stream prefix
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

// for shell completion
func validateArgs(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	var comps []string
	if len(args) == 0 {
		comps = cobra.AppendActiveHelp(comps, "Please select a log source valid options are: kube-scheduler, kube-apiserver-audit, kube-controller-manager, kube-apiserver, authenticator, cloud-controller-manager")
	} else if len(args) == 1 {
		comps = cobra.AppendActiveHelp(comps, "You must specify the URL for the repo you are adding")
	} else {
		comps = cobra.AppendActiveHelp(comps, "This command does not take any more arguments")
	}
	return comps, cobra.ShellCompDirectiveNoFileComp
}

// validates kubelet config for remote logging
func validateKubeletConfig(node string) bool {
	URL := "/api/v1/nodes/" + node + "/proxy/configz"
	kubeletConfigCmdOutput, err := exec.Command("kubectl", "--request-timeout", "20s", "get", "--raw", URL).Output()
	if err != nil {
		log.Fatal(err)
	}
	// b := []byte(kubeletConfigCmdOutput)
	// fmt.Println(b)
	kubeletConfigJson, err := jsonquery.Parse(strings.NewReader(string(kubeletConfigCmdOutput)))
	if err != nil {
		panic(err)
	}

	nodeLogQuery := jsonquery.FindOne(kubeletConfigJson, "kubeletconfig/featureGates/NodeLogQuery")
	// fmt.Printf("%T %v %v\n", nodeLogQuery.Value(), nodeLogQuery.Value(), nodeLogQueryValue)
	systemLogHandler := jsonquery.FindOne(kubeletConfigJson, "kubeletconfig/enableSystemLogHandler")
	systemLogQuery := jsonquery.FindOne(kubeletConfigJson, "kubeletconfig/enableSystemLogQuery")
	if (nodeLogQuery != nil) && (systemLogHandler != nil) && (systemLogQuery != nil) {
		if nodeLogQuery.Value().(bool) {
			return true
		}
	}
	fmt.Println(`	Node is not configured for remote logs.
	Please enable remote logging on the kubelet from the documentation here
	https://kubernetes.io/blog/2023/04/21/node-log-query-alpha/`)
	return false

}
