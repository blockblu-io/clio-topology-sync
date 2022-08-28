package main

import (
	"fmt"
	"github.com/blockblu-io/clio-topology-sync/cmd"
	log "github.com/sirupsen/logrus"
	"os"
)

func main() {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})
	if err := cmd.Execute(); err != nil {
		_,_ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
