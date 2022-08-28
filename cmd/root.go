package cmd

import "github.com/spf13/cobra"

const VERSION string = "1.0.0"

var (
	rootCmd = &cobra.Command{
		Use:   "clio-topology-sync",
		Short: "Topology Manager",
		Long: `clio-topology-sync is a CLI tool that handles the topology management through the CLIO "topology-updater" server.`,
	}
)

func Execute() error {
	return rootCmd.Execute()
}