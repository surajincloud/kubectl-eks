package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version of kubectl-eks",
	Long:  "Print the version of kubectl-eks",
	RunE:  version,
}

func version(cmd *cobra.Command, args []string) error {
	fmt.Println("v0.4.1")
	return nil
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
