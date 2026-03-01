package rpc

import (
	"encoding/json"
	"fmt"
)

// GetNetworkInfoResult is the wire-format response of the getnetworkinfo command.
type GetNetworkInfoResult struct {
	Version            int32          `json:"version"`
	Subversion         string         `json:"subversion"`
	ProtocolVersion    uint32         `json:"protocolversion"`
	LocalServices      string         `json:"localservices"`
	LocalServicesNames []string       `json:"localservicesnames"`
	LocalRelay         bool           `json:"localrelay"`
	TimeOffset         int64          `json:"timeoffset"`
	NetworkActive      bool           `json:"networkactive"`
	Connections        uint32         `json:"connections"`
	Networks           []Network      `json:"networks"`
	RelayFee           float64        `json:"relayfee"`
	IncrementalFee     float64        `json:"incrementalfee"`
	LocalAddresses     []LocalAddress `json:"localaddresses"`
	Warnings           WarningsField  `json:"warnings"`
}

// Network is a network entry in GetNetworkInfoResult.
type Network struct {
	Name                      string `json:"name"`
	Limited                   bool   `json:"limited"`
	Reachable                 bool   `json:"reachable"`
	Proxy                     string `json:"proxy"`
	ProxyRandomizeCredentials bool   `json:"proxy_randomize_credentials"`
}

// LocalAddress is a local address entry in GetNetworkInfoResult.
type LocalAddress struct {
	Address string `json:"address"`
	Port    int    `json:"port"`
	Score   int    `json:"score"`
}

// GetBlockchainInfoResult is the wire-format response of the getblockchaininfo command.
type GetBlockchainInfoResult struct {
	Chain                string        `json:"chain"`
	Blocks               uint32        `json:"blocks"`
	Headers              uint32        `json:"headers"`
	BestBlockHash        string        `json:"bestblockhash"`
	Difficulty           float64       `json:"difficulty"`
	MedianTime           uint32        `json:"mediantime"`
	VerificationProgress float64       `json:"verificationprogress"`
	InitialBlockDownload bool          `json:"initialblockdownload"`
	ChainWork            string        `json:"chainwork"`
	SizeOnDisk           uint64        `json:"size_on_disk"`
	Pruned               bool          `json:"pruned"`
	SoftForks            SoftForks     `json:"softforks"`
	Warnings             WarningsField `json:"warnings"`
}

// SoftForks contains soft fork statuses.
type SoftForks struct {
	Bip34  Fork `json:"bip34"`
	Bip66  Fork `json:"bip66"`
	Bip65  Fork `json:"bip65"`
	Csv    Fork `json:"csv"`
	Segwit Fork `json:"segwit"`
}

// Fork is a single soft fork entry.
type Fork struct {
	Type   string `json:"type"`
	Active bool   `json:"active"`
	Height uint32 `json:"height"`
}

// GetNetworkInfo calls getnetworkinfo and returns the raw wire response.
func (c *Client) GetNetworkInfo() (*GetNetworkInfoResult, error) {
	rawResult, err := c.client.RawRequest("getnetworkinfo", []json.RawMessage{})
	if err != nil {
		return nil, fmt.Errorf("fail to call RawRequest(getnetworkinfo): %w", err)
	}
	var result GetNetworkInfoResult
	if err := json.Unmarshal(rawResult, &result); err != nil {
		return nil, fmt.Errorf("fail to call json.Unmarshal(getnetworkinfo): %w", err)
	}
	return &result, nil
}

// GetBlockchainInfo calls getblockchaininfo and returns the raw wire response.
func (c *Client) GetBlockchainInfo() (*GetBlockchainInfoResult, error) {
	rawResult, err := c.client.RawRequest("getblockchaininfo", []json.RawMessage{})
	if err != nil {
		return nil, fmt.Errorf("fail to call RawRequest(getblockchaininfo): %w", err)
	}
	var result GetBlockchainInfoResult
	if err := json.Unmarshal(rawResult, &result); err != nil {
		return nil, fmt.Errorf("fail to call json.Unmarshal(getblockchaininfo): %w", err)
	}
	return &result, nil
}

// GetBlockCount calls getblockcount and returns the current block height.
func (c *Client) GetBlockCount() (int64, error) {
	rawResult, err := c.client.RawRequest("getblockcount", []json.RawMessage{})
	if err != nil {
		return 0, fmt.Errorf("fail to call RawRequest(getblockcount): %w", err)
	}
	var count int64
	if err := json.Unmarshal(rawResult, &count); err != nil {
		return 0, fmt.Errorf("fail to call json.Unmarshal(getblockcount): %w", err)
	}
	return count, nil
}
