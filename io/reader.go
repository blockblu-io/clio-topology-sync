package io

import (
	"fmt"
	"github.com/blockblu-io/cardano-ops-lib/topology"
	"os"
)

// TopologyReaderError is an error representing the failure of reading
// the topology file at the given path.
type TopologyReaderError struct {
	FilePath string
	Reason   string
}

func (err *TopologyReaderError) Error() string {
	return fmt.Sprintf("couldn't read topology from file at path '%s': %s", err.FilePath, err.Reason)
}

// ReadTopologyFromFile writes the given topology to the given filepath. If the reading fails,
// then a TopologyReaderError is returned. Otherwise nil will be returned.
func ReadTopologyFromFile(inputFilePath string) (*topology.Topology, error) {
	topFile, err := os.Open(inputFilePath)
	if err != nil {
		return nil, &TopologyReaderError {
			FilePath: inputFilePath,
			Reason: err.Error(),
		}
	}
	defer topFile.Close()
	readTopology, err := topology.ReadTopology(topFile)
	if err != nil {
		return nil, &TopologyReaderError {
			FilePath: inputFilePath,
			Reason: err.Error(),
		}
	}
	return readTopology, nil
}
