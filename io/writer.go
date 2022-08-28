package io

import (
	"fmt"
	"github.com/blockblu-io/cardano-ops-lib/topology"
	"os"
	"strconv"
)

// TopologyWriteError is an error representing the failure of writing
// a given topology to a given filepath. A reason should be assigned
// to each error.
type TopologyWriteError struct {
	Topology *topology.Topology
	FilePath string
	Reason   string
}

func (err *TopologyWriteError) Error() string {
	size := "nil"
	if err.Topology != nil && err.Topology.Producers != nil {
		size = strconv.Itoa(len(err.Topology.Producers))
	}
	return fmt.Sprintf("couldn't write topology (size: %s) to file at path '%s': %s", size, err.FilePath,
		err.Reason)
}

// WriteTopologyToFile writes the given topology to the given filepath. If the writing fails,
// then a TopologyWriteError is returned. Otherwise nil will be returned.
func WriteTopologyToFile(top *topology.Topology, outputFilePath string) error {
	outputFile, err := os.OpenFile(outputFilePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return &TopologyWriteError{
			Topology: top,
			FilePath: outputFilePath,
			Reason:   err.Error(),
		}
	}
	defer outputFile.Close()
	err = outputFile.Truncate(0)
	if err != nil {
		return &TopologyWriteError{
			Topology: top,
			FilePath: outputFilePath,
			Reason:   err.Error(),
		}
	}
	err = topology.WriteTopology(top, outputFile)
	if err != nil {
		return &TopologyWriteError{
			Topology: top,
			FilePath: outputFilePath,
			Reason:   err.Error(),
		}
	}
	return nil
}
