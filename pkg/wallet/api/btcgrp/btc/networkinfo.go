package btc

import (
	"encoding/json"

	"github.com/pkg/errors"
)

// GetNetworkInfoResult is response type of PRC `getnetworkinfo`
type GetNetworkInfoResult struct {
	Version            BTCVersion     `json:"version"`
	Subversion         string         `json:"subversion"`
	Protocolversion    uint32         `json:"protocolversion"`
	Localservices      string         `json:"localservices"`
	Localservicesnames []string       `json:"localservicesnames"`
	Localrelay         bool           `json:"localrelay"`
	Timeoffset         int64          `json:"timeoffset"`
	Networkactive      bool           `json:"networkactive"`
	Connections        uint32         `json:"connections"`
	Networks           []Network      `json:"networks"`
	Relayfee           float64        `json:"relayfee"`
	Incrementalfee     float64        `json:"incrementalfee"`
	Localaddresses     []LocalAddress `json:"localaddresses"`
	Warnings           string         `json:"warnings"`
}

// BlockchainInfoChain is chain in GetBlockchainInfoResult
type BlockchainInfoChain string

// chain
const (
	BlockchainInfoChainMain    BlockchainInfoChain = "main"
	BlockchainInfoChainTest    BlockchainInfoChain = "test"
	BlockchainInfoChainRegtest BlockchainInfoChain = "regtest"
)

// String converter
func (c BlockchainInfoChain) String() string {
	return string(c)
}

// GetBlockchainInfoResult is response type of PRC `getblockchaininfo`
type GetBlockchainInfoResult struct {
	Chain                BlockchainInfoChain `json:"chain"` // main, test, regtest
	Blocks               uint32              `json:"blocks"`
	Headers              uint32              `json:"headers"`
	Bestblockhash        string              `json:"bestblockhash"`
	Difficulty           float64             `json:"difficulty"`
	Mediantime           uint32              `json:"mediantime"`
	Verificationprogress float64             `json:"verificationprogress"`
	Initialblockdownload bool                `json:"initialblockdownload"`
	Chainwork            string              `json:"chainwork"`
	SizeOnDisk           uint64              `json:"size_on_disk"`
	Pruned               bool                `json:"pruned"`
	SoftForks            SoftForks           `json:"softforks"`
	Warnings             string              `json:"warnings"`
}

// SoftForks is soft fork list
type SoftForks struct {
	Bip34  Fork `json:"bip34"`
	Bip66  Fork `json:"bip66"`
	Bip65  Fork `json:"bip65"`
	Csv    Fork `json:"csv"`
	Segwit Fork `json:"segwit"`
}

// Fork is fork info
type Fork struct {
	Type   string `json:"type"`
	Active bool   `json:"active"`
	Height uint32 `json:"height"`
}

// Network network info
type Network struct {
	Name                      string `json:"name"`
	Limited                   bool   `json:"limited"`
	Reachable                 bool   `json:"reachable"`
	Proxy                     string `json:"proxy"`
	ProxyRandomizeCredentials bool   `json:"proxy_randomize_credentials"`
}

// LocalAddress local address
type LocalAddress struct {
	Address string `json:"address"`
	Port    int    `json:"port"`
	Score   int    `json:"score"`
}

// GetNetworkInfo call RPC `getnetworkinfo`
func (b *Bitcoin) GetNetworkInfo() (*GetNetworkInfoResult, error) {
	rawResult, err := b.Client.RawRequest("getnetworkinfo", []json.RawMessage{})
	if err != nil {
		return nil, errors.Wrap(err, "fail to call json.RawRequest(getnetworkinfo)")
	}

	networkInfoResult := GetNetworkInfoResult{}
	err = json.Unmarshal(rawResult, &networkInfoResult)
	if err != nil {
		return nil, errors.Wrap(err, "fail to call json.Unmarshal()")
	}

	return &networkInfoResult, nil
}

// GetBlockchainInfo call RPC `getblockchaininfo`
func (b *Bitcoin) GetBlockchainInfo() (*GetBlockchainInfoResult, error) {
	rawResult, err := b.Client.RawRequest("getblockchaininfo", []json.RawMessage{})
	if err != nil {
		return nil, errors.Wrap(err, "fail to call json.RawRequest(getblockchaininfo)")
	}

	blockchainInfoResult := GetBlockchainInfoResult{}
	err = json.Unmarshal(rawResult, &blockchainInfoResult)
	if err != nil {
		return nil, errors.Wrap(err, "fail to call json.Unmarshal()")
	}

	return &blockchainInfoResult, nil
}
