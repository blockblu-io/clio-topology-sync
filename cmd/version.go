package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of clio-topology-updater",
	Long:  `All software has versions. This is clio-topology-updater's`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("clio-topology-updater %s\n", VERSION)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

