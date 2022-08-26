package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"
	"time"

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

	w := tabwriter.NewWriter(os.Stdout, 5, 2, 3, ' ', tabwriter.TabIndent)
	defer w.Flush()
	fmt.Fprintln(w, "NAME", "\t", "ARCH", "\t", "INSTANCE-TYPE", "\t", "OS", "\t", "CAPACITY-TYPE", "\t", "REGION", "\t", "AGE")
	for _, i := range nodeList {
		age := getAge(i.CreationTimestamp)
		fmt.Fprintln(w, i.Name, "\t", i.Labels[kube.ArchLabel], "\t", i.Labels[kube.InstanceTypeLabel], "\t", i.Labels[kube.OsLabel], "\t", i.Labels[kube.CapacityTypeLabel], "\t", i.Labels[kube.ZoneLabel], "\t", age)
	}

	return nil
}

func getAge(creationStamp metav1.Time) string {

	currentTime := time.Now()
	diff := currentTime.Sub(creationStamp.Time)
	return amt.HumanDuration(diff)
}
