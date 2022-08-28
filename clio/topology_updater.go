package clio

import (
	"context"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// TipUpdater
type TipUpdater interface {
	// Run
	Run(ctx context.Context, api API)
}

type tipUpdater struct {
	Hostname              string
	Port                  uint16
	Valency               uint
	NetworkMagic          uint
	PrometheusEndpointURL string
}

func NewTipUpdater(hostname string, port uint16, valency uint, networkMagic uint, prometheusEndpointURL string) TipUpdater {
	return &tipUpdater{
		Hostname:              hostname,
		Port:                  port,
		Valency:               valency,
		NetworkMagic:          networkMagic,
		PrometheusEndpointURL: prometheusEndpointURL,
	}
}

func (tip *tipUpdater) Run(ctx context.Context, api API) {
	client := http.Client{
		Timeout: 5 * time.Second,
	}
	gap := 0 * time.Second
	for {
		select {
		case <-time.After(gap):
			gap = tip.runTipUpdate(client, api)
		case <-ctx.Done():
			break
		}
	}
}

func (tip *tipUpdater) runTipUpdate(client http.Client, api API) time.Duration {
	log.Info("getting the tip from the node.")
	blockNumber, err := getTip(client, tip.PrometheusEndpointURL)
	if err == nil {
		log.Infof("fetched the tip %d.", blockNumber)
		err = api.UpdateTip(tip.Hostname, tip.Port, tip.Valency, blockNumber, tip.NetworkMagic)
		if err == nil {
			return 1 * time.Hour
		} else {
			log.Errorf("posting the tip failed: %s", err.Error())
			return 10 * time.Minute
		}
	} else {
		log.Errorf("getting the tip failed: %s", err.Error())
		return 1 * time.Minute
	}
}

func getTip(client http.Client, prometheusEndpointURL string) (uint64, error) {
	resp, err := client.Get(prometheusEndpointURL)
	if err == nil {
		defer resp.Body.Close()
		data, err := ioutil.ReadAll(resp.Body)
		if err == nil {
			lines := strings.Split(string(data), "\n")
			for _, line := range lines {
				pair := strings.Split(line, " ")
				if len(pair) == 2 {
					if pair[0] == "cardano_node_metrics_blockNum_int" {
						blockNumber, err := strconv.ParseUint(pair[1], 10, 64)
						if err == nil {
							return blockNumber, nil
						} else {
							return 0, errors.New(fmt.Sprintf("couldn't parse the value: %s", pair[1]))
						}
					}
				}
			}
			return 0, errors.New("couldn't find the 'cardano_node_metrics_blockNum_int' value")
		} else {
			return 0, errors.New(fmt.Sprintf("the body of the call to the prometheus endpoint cannot be read: %s",
				err.Error()))
		}
	} else {
		return 0, errors.New(fmt.Sprintf("the prometheus endpoint at '%s' cannot be reached: %s",
			prometheusEndpointURL, err.Error()))
	}
}
