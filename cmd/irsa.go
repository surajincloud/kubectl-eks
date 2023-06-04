package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"github.com/surajincloud/kubectl-eks/pkg/kube"
)

// irsaCmd represents the irsa command
var irsaCmd = &cobra.Command{
	Use:   "irsa",
	Short: "List Serviceaccounts with their IRSA information",
	Long:  "List Serviceaccounts with their IRSA information",
	RunE:  irsa,
}

func irsa(cmd *cobra.Command, args []string) error {

	saList, err := kube.GetSA(KubernetesConfigFlags)
	if err != nil {
		return err
	}
	w := tabwriter.NewWriter(os.Stdout, 5, 2, 3, ' ', tabwriter.TabIndent)
	defer w.Flush()
	fmt.Fprintln(w, "NAMESPACE", "\t", "SERVICEACCOUNT", "\t", "IAM-ROLE", "\t", "TOKEN-EXPIRATION")
	for _, i := range saList {
		if i.Annotations["eks.amazonaws.com/role-arn"] != "" {
			fmt.Fprintln(w, i.Namespace, "\t", i.Name, "\t", i.Annotations["eks.amazonaws.com/role-arn"], "\t", i.Annotations["eks.amazonaws.com/token-expiration"])
		}
	}
	return nil
}
func init() {
	rootCmd.AddCommand(irsaCmd)
}
