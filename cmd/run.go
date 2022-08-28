package cmd

import (
	"context"
	"errors"
	"fmt"
	"github.com/blockblu-io/clio-topology-sync/clio"
	"github.com/spf13/cobra"
	"net/url"
	"os"
	"os/signal"
)

var (
	hostname          string
	port              uint16
	valency           uint
	topologyPath      string
	fixedTopologyPath string
	endpointURL       string
	maxPeers          uint
	networkMagic      uint
	runCmd            = &cobra.Command{
		Use:   "run",
		Short: "Run the topology management",
		Long:  `Run the topology management for the given peer`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if hostname == "" {
				return errors.New("hostname must be specified, but was empty")
			}
			if port == 0 {
				return errors.New(fmt.Sprintf("port must be between 1 and 65535, but was %d", port))
			}
			if valency == 0 {
				return errors.New(fmt.Sprintf("valency must be greater than or equal to 1, but was %d", port))
			}
			if endpointURL == "" {
				return errors.New("URL of prometheus endpoint must be specified, but was empty")
			}
			metricUrl, err := url.Parse(endpointURL)
			if err != nil {
				return errors.New(fmt.Sprintf("specified URL for prometheus endpoint isn't valid: %s",
					err.Error()))
			}
			sigChan := make(chan os.Signal, 1)
			signal.Notify(sigChan, os.Interrupt)
			ctx, cancel := context.WithCancel(context.Background())
			api := clio.GetAPI("https://api.clio.one/htopology/v1/")
			fetcher := clio.NewTopologyFetcher(topologyPath, fixedTopologyPath, maxPeers, networkMagic)
			go fetcher.Run(ctx, api)
			tipUpdater := clio.NewTipUpdater(hostname, port, valency, networkMagic, metricUrl.String())
			go tipUpdater.Run(ctx, api)
			select {
			case <-sigChan:
				cancel()
			}
			return nil
		},
	}
)

func init() {
	runCmd.Flags().StringVar(&hostname, "hostname", "", "hostname (IP address or DNS name) of the peer(s)")
	_ = runCmd.MarkFlagRequired("hostname")
	runCmd.Flags().Uint16Var(&port, "port", 0, "port of the peer(s)")
	_ = runCmd.MarkFlagRequired("port")
	runCmd.Flags().UintVar(&valency, "valency", 1, "valency (1 for IP address, and >= 1 for DNS names)")
	runCmd.Flags().StringVar(&topologyPath, "topology-path", "", "path to the output file for the topology")
	_ = runCmd.MarkFlagRequired("topology-path")
	runCmd.Flags().StringVar(&fixedTopologyPath, "fixed-topology-path", "", "path to the output file for the topology")
	runCmd.Flags().StringVar(&endpointURL, "prometheus-endpoint-url", "", "url of the prometheus endpoint of the peer")
	_ = runCmd.MarkFlagRequired("prometheus-endpoint-url")
	runCmd.Flags().UintVar(&maxPeers, "max-peers", 10, "the number of peers that shall be fetched")
	runCmd.Flags().UintVar(&networkMagic, "network-magic", 764824073, "the magic number of the network")
	rootCmd.AddCommand(runCmd)
}
