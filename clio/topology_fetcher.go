package clio

import (
	"context"
	"github.com/blockblu-io/cardano-ops-lib/topology"
	"github.com/blockblu-io/clio-topology-sync/io"
	log "github.com/sirupsen/logrus"
	"math"
	"time"
)


type TopologyFetcher interface {
	// Run
	Run(ctx context.Context, api API)
}

type topologyFetcher struct {
	OutputFilePath string
	FixedFilePath  string
	MaxPeers       uint
	NetworkMagic   uint
}

func NewTopologyFetcher(outputFilePath, fixedFilePath string, maxPeers uint, networkMagic uint) TopologyFetcher {
	return &topologyFetcher{
		OutputFilePath: outputFilePath,
		FixedFilePath:  fixedFilePath,
		MaxPeers:       maxPeers,
		NetworkMagic:   networkMagic,
	}
}

func (fetcher *topologyFetcher) Run(ctx context.Context, api API) {
	gap := 0 * time.Second
	for {
		time.Sleep(gap)
		fctx, fCancel := context.WithCancel(ctx)
		select {
		case <-time.After(gap):
			gap = fetcher.runFetcher(fctx, api)
		case <-ctx.Done():
			fCancel()
			break
		}
	}
}

func (fetcher *topologyFetcher) runFetcher(ctx context.Context, api API) time.Duration {
	start := time.Now()
	log.Infof("fetching topology from CLIO API.")
	var top *topology.Topology
	for n := 1; n <= 3; n++ {
		select {
		case <-ctx.Done():
			return 0 * time.Second
		case <-time.After(time.Duration(math.Pow(10.0, float64(n-1))) * time.Second):
			var err error
			top, err = api.FetchTopology(fetcher.MaxPeers, fetcher.NetworkMagic)
			if err == nil {
				break
			}
			log.Errorf("couldn't fetch the topology (%d. attempt). %s", n, err.Error())
		}
	}
	size := 0
	if top != nil && top.Producers != nil {
		size = len(top.Producers)
	}
	log.Infof("fetched the topology of size %d.", size)
	if fetcher.FixedFilePath != "" {
		log.Info("merging fetched topology with the fixed one.")
		fixedTop, err := io.ReadTopologyFromFile(fetcher.FixedFilePath)
		if err != nil {
			log.Error(err.Error())
		}
		top = topology.Merge(fixedTop, top)
		size = len(top.Producers)
		log.Infof("new merged topology of size %d.", size)
	}
	if top != nil {
		err := io.WriteTopologyToFile(top, fetcher.OutputFilePath)
		if err != nil {
			log.Error(err.Error())
		}
	}
	return start.Add(1 * time.Hour).Sub(time.Now())
}
