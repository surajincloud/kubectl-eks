/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/antchfx/jsonquery"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go-v2/service/eks"
	"github.com/markusmobius/go-dateparser"
	"github.com/spf13/cobra"
	awspkg "github.com/surajincloud/kubectl-eks/pkg/aws"
	"github.com/surajincloud/kubectl-eks/pkg/kube"
)

var Follow bool
var SinceTime string

func init() {
	logsCmd.Flags().BoolVarP(&Follow, "follow", "f", false, "Follow logs (not available for node file queries)")
	logsCmd.Flags().StringVar(&SinceTime, "since", "1 hour ago", "What time logs should start from")
}

// logsCmd represents the logs command
var logsCmd = &cobra.Command{
	Use:               "logs [flags] LOG_SOURCE",
	ValidArgsFunction: validateArgs,
	Example: `    kubectl eks logs kube-apiserver
    kubectl eks logs NODE [kubelet]
  
  Query multiple log sources:
    kubectl eks logs api audit scheduler
    kubectl eks logs NODE kubelet containerd`,
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

	// we only look at the first argument to determine if it is a control plane log source or node
	logTarget := args[0]
	logsChan := make(chan string)
	logsDoneChan := make(chan bool, len(args))

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

	// handle ctl+c interrupt
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		os.Exit(0)
	}()

	// check first argument if it is a control plane log source or node
	if contains(cloudwatchLogStreams, logTarget) || (logTarget == "all") {

		if logTarget == "all" {
			args = args[1:]
			args = append(args, "scheduler", "kube-apiserver-audit", "kube-controller-manager", "kube-apiserver", "authenticator", "cloud-controller-manager")
		}

		var cwlStreamInput cloudwatchlogs.DescribeLogStreamsInput
		var streams *cloudwatchlogs.DescribeLogStreamsOutput
		var limit int32 = 100
		var cwlGroupPrefix string

		cwlGroupName := "/aws/eks/" + clusterName + "/cluster"
		cwlStreamInput.LogGroupName = &cwlGroupName

		// aws config
		ctx := context.Background()
		// read flag values
		region, _ := cmd.Flags().GetString("region")

		cfg, err := awspkg.GetAWSConfig(ctx, region)
		if err != nil {
			log.Fatal(err)
		}

		// TODO check if logging is enabled
		// TODO allow user to enable logging
		svc := eks.NewFromConfig(cfg)
		eksClusterInput := &eks.DescribeClusterInput{
			Name: aws.String(clusterName),
		}
		result, _ := svc.DescribeCluster(ctx, eksClusterInput)
		// check if logging is enabled
		if !*result.Cluster.Logging.ClusterLogging[0].Enabled == true {
			fmt.Println("Logging is not enabled for this cluster. Please enable logging and try again.")
			os.Exit(1)
		}

		cwl := cloudwatchlogs.NewFromConfig(cfg)
		// verify the group exists first
		err = ensureLogGroupExists(cwlGroupName, ctx, cwl)
		//TODO prompt if user wants to enable logs
		if err != nil {
			panic(err)
		}

		var fetchStream string
		matchedStream := 0
		for _, logSource := range args {

			cwlGroupPrefix = getLogStreamPrefix(logSource)
			cwlStreamInput.LogStreamNamePrefix = &cwlGroupPrefix

			streams, err = cwl.DescribeLogStreams(ctx, &cwlStreamInput)
			if err != nil {
				log.Fatal(err)
			}

			if len(streams.LogStreams) == 0 {
				fmt.Fprintln(os.Stderr, "No log streams found for", cwlGroupPrefix)
			} else if cwlGroupPrefix == "kube-apiserver" {
				// we have to make sure kube-apiserver doesn't add kube-apiserver-audit logs
				for _, stream := range streams.LogStreams {
					if !strings.Contains(*stream.LogStreamName, "kube-apiserver-audit") {
						fetchStream = *stream.LogStreamName
						// only match the first stream
						break
					}
				}
				if fetchStream == "" {
					fmt.Fprintln(os.Stderr, "No log streams found for", cwlGroupPrefix)
				} else {
					go getLogEvents(&cwlGroupName, &fetchStream, &limit, logsChan, logsDoneChan, ctx, cwl, len(args))
					matchedStream++

				}
			} else {
				fetchStream = *streams.LogStreams[0].LogStreamName
				go getLogEvents(&cwlGroupName, &fetchStream, &limit, logsChan, logsDoneChan, ctx, cwl, len(args))
				matchedStream++
			}
		}

		// use a count to make sure we have a goroutine fetching logs
		if matchedStream > 0 {
			// Print each line from logsChan
			for log := range logsChan {
				// print the log
				fmt.Println(log)

				// once all event streams are done, close the channel
				if len(logsDoneChan) == len(args) {
					close(logsChan)
				}
			}
		}
	} else {
		// we assume the target is a node instead of control plane
		nodeList, err := kube.GetNodes(KubernetesConfigFlags)
		if err != nil {
			return err
		}

		var nodeMatched bool = false
		var currentNodeSlice []string
		logTargetSlice := strings.Split(logTarget, ".")

		for _, i := range nodeList {
			// match node based on substring eg. ip-192-168-1-1
			currentNodeSlice = strings.Split(i.Name, ".")
			if currentNodeSlice[0] == logTargetSlice[0] {
				nodeMatched = true
				var query []string

				// use all additional arguments as services to query
				if len(args) > 1 {
					query = args[1:]
				} else {
					query = append(query, "kubelet")
				}

				// validate kubelet settings for remote logs
				if validateKubeletConfig(i.Name) {
					// get logs and assume query is journald and can accept sinceTime
					go getNodeLogs(i.Name, query, false, logsChan, logsDoneChan)

					// print each line from logsChan
					for log := range logsChan {

						fmt.Println(log)

						if len(logsDoneChan) == 1 {
							close(logsChan)
						}
					}
				}
			}
		}
		if nodeMatched {
			return nil
		} else {
			fmt.Printf("Node %s not found\nTo query control plane logs please see options in --help output.\n", logTarget)
		}
	}
	return nil
}

// get node logs
func getNodeLogs(node string, query []string, fileQuery bool, logsChan chan<- string, logsDoneChan chan<- bool) {

	dt, err := dateparser.Parse(nil, SinceTime)
	if err != nil {
		panic(err)
	}
	// create URL for log fetching
	rawURL := "/api/v1/nodes/" + node + "/proxy/logs/"

	for {

		clientSet, _ := kube.ClientSet(KubernetesConfigFlags)
		req := clientSet.RESTClient().Get().
			AbsPath(rawURL)

		for _, q := range query {
			req.Param("query", q)
		}

		if !fileQuery {
			// sinceTime is ignored for file queries
			req.Param("sinceTime", dt.Time.Format(time.RFC3339))
		}

		resp, err := req.DoRaw(context.Background())
		if err != nil {
			log.Panicln(err, req.URL().String())
		}

		// api returns a byte string
		// convert to string and split by newline to send each line to channel
		for _, logLine := range strings.Split(string(resp), "\n") {
			// fmt.Printf("%v %s", lineNumber, logLine)
			if strings.Contains(logLine, "options present and query resolved to log files") {
				// file queries cannot use sinceTime
				// we catch this error output and run again as file query
				go getNodeLogs(node, query, true, logsChan, logsDoneChan)
				break
			} else if logLine == "" ||
				strings.Contains(logLine, "-- No entries --") ||
				strings.Contains(logLine, "-- Logs begin at ") {
				// don't send log decorations
			} else {
				logsChan <- logLine
			}
		}

		if Follow {
			if fileQuery {
				fmt.Fprintln(os.Stderr, "Cannot follow file queries")
				close(logsChan)
				break
			}
			dt, err = dateparser.Parse(nil, "now")
			if err != nil {
				panic(err)
			}
			time.Sleep(1 * time.Second)
		} else {
			if fileQuery {
				close(logsChan)
			} else {
				logsDoneChan <- true
			}
			break
		}
	}
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

func getLogEvents(logGroupName *string, logStreamName *string, limit *int32, logsChan chan<- string, logsDoneChan chan<- bool, ctx context.Context, cwl *cloudwatchlogs.Client, totalStreams int) {

	dt, err := dateparser.Parse(nil, SinceTime)
	if err != nil {
		panic(err)
	}

	// loop forever if Follow == true
	for {
		resp, err := cwl.GetLogEvents(ctx, &cloudwatchlogs.GetLogEventsInput{
			Limit:         limit,
			LogGroupName:  logGroupName,
			LogStreamName: logStreamName,
			StartTime:     aws.Int64(dt.Time.UnixMilli()),
		})
		if err != nil {
			panic(err)
		}

		for i, event := range resp.Events {
			// TODO allow for following tokens for more logs from different streams

			if i == len(resp.Events)-1 {
				if Follow {
					dt, err = dateparser.Parse(nil, "now")
					if err != nil {
						panic(err)
					}
					// wait 1 sec before querying again
					time.Sleep(1 * time.Second)
				} else {
					logsDoneChan <- true
				}
				logsChan <- *event.Message
			} else {
				logsChan <- *event.Message
			}

		}

	}
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
	rawURL := "/api/v1/nodes/" + node + "/proxy/configz"
	clientSet, _ := kube.ClientSet(KubernetesConfigFlags)
	req := clientSet.RESTClient().Get().
		AbsPath(rawURL).Timeout(20 * time.Second)

	resp, err := req.DoRaw(context.Background())
	if err != nil {
		log.Panicln(err, req.URL().String())
	}

	// read the kueblet config json
	kubeletConfigJson, err := jsonquery.Parse(strings.NewReader(string(resp)))
	if err != nil {
		panic(err)
	}

	// check if the node has the appropriate config for remote logging
	nodeLogQuery := jsonquery.FindOne(kubeletConfigJson, "kubeletconfig/featureGates/NodeLogQuery")
	// fmt.Printf("%T %v %v\n", nodeLogQuery.Value(), nodeLogQuery.Value(), nodeLogQueryValue)
	systemLogHandler := jsonquery.FindOne(kubeletConfigJson, "kubeletconfig/enableSystemLogHandler")
	systemLogQuery := jsonquery.FindOne(kubeletConfigJson, "kubeletconfig/enableSystemLogQuery")
	if (nodeLogQuery != nil) && (systemLogHandler != nil) && (systemLogQuery != nil) {
		if nodeLogQuery.Value().(bool) {
			return true
		}
	}
	fmt.Printf(`	Node %s is not configured for remote logs.
	Please enable remote logging on the kubelet from the documentation here
	Requires Kubernetes 1.27 https://kubernetes.io/blog/2023/04/21/node-log-query-alpha/`, node)

	return false

}
