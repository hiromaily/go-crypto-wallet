package api

import (
	"encoding/json"

	"github.com/pkg/errors"
)

// GetNetworkInfoResult getnetworkinfo RPC のレスポンス
type GetNetworkInfoResult struct {
	Version         int       `json:"version"`
	Subversion      string    `json:"subversion"`
	Protocolversion int       `json:"protocolversion"`
	Localservices   string    `json:"localservices"`
	Localrelay      bool      `json:"localrelay"`
	Timeoffset      int       `json:"timeoffset"`
	Networkactive   bool      `json:"networkactive"`
	Connections     int       `json:"connections"`
	Networks        []Network `json:"networks"`
	Relayfee        float64   `json:"relayfee"`
	Incrementalfee  float64   `json:"incrementalfee"`
	Localaddresses  []string  `json:"localaddresses"`
	Warnings        string    `json:"warnings"`
}

// Network ネットワーク情報
type Network struct {
	Name                      string `json:"name"`
	Limited                   bool   `json:"limited"`
	Reachable                 bool   `json:"reachable"`
	Proxy                     string `json:"proxy"`
	ProxyRandomizeCredentials bool   `json:"proxy_randomize_credentials"`
}

// GetNetworkInfo getnetworkinfo RPC をcallする
func (b *Bitcoin) GetNetworkInfo() (*GetNetworkInfoResult, error) {
	rawResult, err := b.client.RawRequest("getnetworkinfo", []json.RawMessage{})
	if err != nil {
		return nil, errors.Errorf("json.RawRequest(logging): error: %v", err)
	}

	networkInfoResult := GetNetworkInfoResult{}
	err = json.Unmarshal([]byte(rawResult), &networkInfoResult)
	if err != nil {
		return nil, errors.Errorf("json.Unmarshal(): error: %v", err)
	}

	return &networkInfoResult, nil
}
