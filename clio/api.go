package clio

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/blockblu-io/cardano-ops-lib/topology"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"
)

// API for topology-management from CLIO for the Cardano blockchain.
type API interface {

	// FetchTopology fetches the topology from the CLIO API for this
	// peer. The topology will include the given number of peers and
	// be for the network with the given magic number. If the
	// fetching fails, then nil and error will be returned. Otherwise
	// the fetched topology is returned.
	FetchTopology(maxPeers uint, networkMagic uint) (*topology.Topology, error)

	// UpdateTip
	UpdateTip(hostname string, port uint16, valency uint, blockNumber uint64, networkMagic uint) error
}

type httpAPI struct {
	apiURL string
	client *http.Client
}

type APIResponse struct {
	ResultCode string `json:"resultcode"`
	Message    string `json:"msg"`
}

// GetAPI gets the API of the topology manager with the given URL.
func GetAPI(apiURL string) API {
	dialer := net.Dialer{
		Timeout: 5 * time.Second,
	}
	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.DialContext = func(ctx context.Context, network, addr string) (net.Conn, error) {
		return dialer.DialContext(ctx, "tcp4", addr)
	}
	client := &http.Client{
		Timeout: 5 * time.Second,
		Transport: transport,
	}
	return &httpAPI{
		apiURL: apiURL,
		client: client,
	}
}

func (clio *httpAPI) FetchTopology(maxPeers uint, networkMagic uint) (*topology.Topology, error) {
	parameterizedURL, err := url.Parse(clio.apiURL)
	if err != nil {
		log.Fatal(err)
	}
	parameterizedURL.Path = path.Join(parameterizedURL.Path, "fetch")
	q := parameterizedURL.Query()
	q.Set("max", fmt.Sprintf("%d", maxPeers))
	q.Set("magic", fmt.Sprintf("%d", networkMagic))
	parameterizedURL.RawQuery = q.Encode()
	resp, err := clio.client.Get(parameterizedURL.String())
	if err != nil {
		return nil, errors.New(fmt.Sprintf("api fetching error: %s", err.Error()))
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("api fetching error: %s", err.Error()))
	}
	var apiResp APIResponse
	err = json.Unmarshal(data, &apiResp)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("api fetching error: %s", err.Error()))
	}
	if !strings.HasPrefix(apiResp.ResultCode, "2") {
		return nil, errors.New(fmt.Sprintf("api fetching error: %s (status code: %s)", apiResp.Message,
			apiResp.ResultCode))
	}
	bodyReader := bytes.NewReader(data)
	top, err := topology.ReadTopology(bodyReader)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("api fetching error: %s", err.Error()))
	}
	return top, nil
}

func (clio *httpAPI) UpdateTip(hostname string, port uint16, valency uint, blockNumber uint64, networkMagic uint) error {
	parameterizedURL, err := url.Parse(clio.apiURL)
	if err != nil {
		log.Fatal(err)
	}
	q := parameterizedURL.Query()
	q.Set("blockNo", fmt.Sprintf("%d", blockNumber))
	q.Set("magic", fmt.Sprintf("%d", networkMagic))
	q.Set("hostname", fmt.Sprintf("%s", hostname))
	q.Set("port", fmt.Sprintf("%d", port))
	q.Set("valency", fmt.Sprintf("%d", valency))
	parameterizedURL.RawQuery = q.Encode()
	resp, err := clio.client.Get(parameterizedURL.String())
	if err != nil {
		return errors.New(fmt.Sprintf("api update error: %s", err.Error()))
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.New(fmt.Sprintf("api update error: %s", err.Error()))
	}
	var apiResp APIResponse
	err = json.Unmarshal(data, &apiResp)
	if err != nil {
		return errors.New(fmt.Sprintf("api update error: %s", err.Error()))
	}
	if !strings.HasPrefix(apiResp.ResultCode, "2") {
		return errors.New(fmt.Sprintf("api update error: %s (status code: %s)", apiResp.Message,
			apiResp.ResultCode))
	}
	return nil
}
