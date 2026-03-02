package rpc

import (
	"encoding/json"
	"fmt"
)

// ListUnspentResult is the wire-format response of the listunspent command.
type ListUnspentResult struct {
	TxID          string  `json:"txid"`
	Vout          uint32  `json:"vout"`
	Address       string  `json:"address"`
	Label         string  `json:"label"`
	RedeemScript  string  `json:"redeemScript"`
	WitnessScript string  `json:"witnessScript"`
	ScriptPubKey  string  `json:"scriptPubKey"`
	Amount        float64 `json:"amount"`
	Confirmations int64   `json:"confirmations"`
	Spendable     bool    `json:"spendable"`
	Solvable      bool    `json:"solvable"`
	Desc          string  `json:"desc"`
	Safe          bool    `json:"safe"`
}

// ListUnspent calls listunspent with a minimum confirmation count.
func (c *rpcClient) ListUnspent(minConf uint64) ([]ListUnspentResult, error) {
	input, err := json.Marshal(minConf)
	if err != nil {
		return nil, fmt.Errorf("fail to call json.Marshal(minConf): %w", err)
	}
	rawResult, err := c.client.RawRequest("listunspent", []json.RawMessage{input})
	if err != nil {
		return nil, fmt.Errorf("fail to call RawRequest(listunspent): %w", err)
	}
	var result []ListUnspentResult
	if err := json.Unmarshal(rawResult, &result); err != nil {
		return nil, fmt.Errorf("fail to call json.Unmarshal(listunspent): %w", err)
	}
	return result, nil
}

// ListUnspentByAddresses calls listunspent with min/max confirmation counts and an address filter.
func (c *rpcClient) ListUnspentByAddresses(minConf, maxConf uint64, addresses []string) ([]ListUnspentResult, error) {
	input1, err := json.Marshal(minConf)
	if err != nil {
		return nil, fmt.Errorf("fail to call json.Marshal(minConf): %w", err)
	}
	input2, err := json.Marshal(maxConf)
	if err != nil {
		return nil, fmt.Errorf("fail to call json.Marshal(maxConf): %w", err)
	}
	input3, err := json.Marshal(addresses)
	if err != nil {
		return nil, fmt.Errorf("fail to call json.Marshal(addresses): %w", err)
	}
	rawResult, err := c.client.RawRequest("listunspent", []json.RawMessage{input1, input2, input3})
	if err != nil {
		return nil, fmt.Errorf("fail to call RawRequest(listunspent): %w", err)
	}
	var result []ListUnspentResult
	if err := json.Unmarshal(rawResult, &result); err != nil {
		return nil, fmt.Errorf("fail to call json.Unmarshal(listunspent): %w", err)
	}
	return result, nil
}
