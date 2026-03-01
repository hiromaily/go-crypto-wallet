package rpc

import (
	"encoding/json"
	"fmt"
)

// LoggingResult is the wire-format response of the logging command.
type LoggingResult struct {
	Net         bool `json:"net"`
	Tor         bool `json:"tor"`
	Mempool     bool `json:"mempool"`
	HTTP        bool `json:"http"`
	Bench       bool `json:"bench"`
	Zmq         bool `json:"zmq"`
	Walletdb    bool `json:"walletdb"`
	RPC         bool `json:"rpc"`
	Estimatefee bool `json:"estimatefee"`
	Addrman     bool `json:"addrman"`
	Selectcoins bool `json:"selectcoins"`
	Reindex     bool `json:"reindex"`
	Cmpctblock  bool `json:"cmpctblock"`
	Rand        bool `json:"rand"`
	Prune       bool `json:"prune"`
	Proxy       bool `json:"proxy"`
	Mempoolrej  bool `json:"mempoolrej"`
	Libevent    bool `json:"libevent"`
	Coindb      bool `json:"coindb"`
	Qt          bool `json:"qt"`
	Leveldb     bool `json:"leveldb"`
	Validation  bool `json:"validation"`
}

// Logging calls the logging RPC and returns the current logging configuration.
func (c *Client) Logging() (*LoggingResult, error) {
	rawResult, err := c.client.RawRequest("logging", []json.RawMessage{})
	if err != nil {
		return nil, fmt.Errorf("fail to call RawRequest(logging): %w", err)
	}
	var result LoggingResult
	if err := json.Unmarshal(rawResult, &result); err != nil {
		return nil, fmt.Errorf("fail to call json.Unmarshal(logging): %w", err)
	}
	return &result, nil
}
