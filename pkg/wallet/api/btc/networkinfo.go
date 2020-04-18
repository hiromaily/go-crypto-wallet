package btc

import (
	"encoding/json"

	"github.com/pkg/errors"

	"github.com/hiromaily/go-bitcoin/pkg/wallet/coin"
)

// GetNetworkInfoResult is response type of PRC `getnetworkinfo`
type GetNetworkInfoResult struct {
	Version            coin.BTCVersion `json:"version"`
	Subversion         string          `json:"subversion"`
	Protocolversion    int             `json:"protocolversion"`
	Localservices      string          `json:"localservices"`
	Localservicesnames []string        `json:"localservicesnames"`
	Localrelay         bool            `json:"localrelay"`
	Timeoffset         int             `json:"timeoffset"`
	Networkactive      bool            `json:"networkactive"`
	Connections        int             `json:"connections"`
	Networks           []Network       `json:"networks"`
	Relayfee           float64         `json:"relayfee"`
	Incrementalfee     float64         `json:"incrementalfee"`
	Localaddresses     []LocalAddress  `json:"localaddresses"`
	Warnings           string          `json:"warnings"`
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
	rawResult, err := b.client.RawRequest("getnetworkinfo", []json.RawMessage{})
	if err != nil {
		return nil, errors.Wrap(err, "fail to call json.RawRequest(getnetworkinfo)")
	}

	networkInfoResult := GetNetworkInfoResult{}
	err = json.Unmarshal([]byte(rawResult), &networkInfoResult)
	if err != nil {
		return nil, errors.Wrap(err, "fail to call json.Unmarshal()")
	}

	return &networkInfoResult, nil
}
