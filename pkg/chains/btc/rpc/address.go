package rpc

import (
	"encoding/json"
	"fmt"
)

// GetAddressInfoResult is the wire-format response of the getaddressinfo command.
type GetAddressInfoResult struct {
	Address             string         `json:"address"`
	ScriptPubKey        string         `json:"scriptPubKey"`
	IsMine              bool           `json:"ismine"`
	Solvable            bool           `json:"solvable,omitempty"`
	Desc                string         `json:"desc,omitempty"`
	IsWatchOnly         bool           `json:"iswatchonly"`
	IsScript            bool           `json:"isscript"`
	IsWitness           bool           `json:"iswitness,omitempty"`
	Pubkey              string         `json:"pubkey,omitempty"`
	IsCompressed        bool           `json:"iscompressed,omitempty"`
	IsChange            bool           `json:"ischange"`
	Timestamp           int64          `json:"timestamp,omitempty"`
	HDKeyPath           string         `json:"hdkeypath,omitempty"`
	HDSeedID            string         `json:"hdseedid,omitempty"`
	HDMasterFingerprint string         `json:"hdmasterfingerprint,omitempty"`
	Labels              FlexibleLabels `json:"labels"`
	// Label is a legacy single-string field used by older Bitcoin Cash nodes.
	// Callers should prefer Labels; this field is only populated by BCH nodes
	// that return "label" instead of "labels".
	Label string `json:"label,omitempty"`
}

// ValidateAddressResult is the wire-format response of the validateaddress command.
type ValidateAddressResult struct {
	IsValid           bool   `json:"isvalid"`
	Address           string `json:"address"`
	ScriptPubKey      string `json:"scriptPubKey"`
	IsScript          bool   `json:"isscript"`
	IsWitness         bool   `json:"iswitness"`
	WitnessVersion    uint32 `json:"witness_version,omitempty"`
	WitnessProgramHex string `json:"witness_program,omitempty"`
}

// Purpose is the intermediate response type element for getaddressesbylabel.
type Purpose struct {
	Purpose string `json:"purpose"`
}

// GetAddressInfo calls the getaddressinfo RPC and returns the raw wire response.
func (c *rpcClient) GetAddressInfo(addr string) (*GetAddressInfoResult, error) {
	input, err := json.Marshal(addr)
	if err != nil {
		return nil, fmt.Errorf("fail to call json.Marshal(addr): %w", err)
	}
	rawResult, err := c.client.RawRequest("getaddressinfo", []json.RawMessage{input})
	if err != nil {
		return nil, fmt.Errorf("fail to call RawRequest(getaddressinfo) %s: %w", addr, err)
	}
	var result GetAddressInfoResult
	if err := json.Unmarshal(rawResult, &result); err != nil {
		return nil, fmt.Errorf("fail to call json.Unmarshal(getaddressinfo): %w", err)
	}
	return &result, nil
}

// ValidateAddress calls the validateaddress RPC and returns the raw wire response.
func (c *rpcClient) ValidateAddress(addr string) (*ValidateAddressResult, error) {
	input, err := json.Marshal(addr)
	if err != nil {
		return nil, fmt.Errorf("fail to call json.Marshal(addr): %w", err)
	}
	rawResult, err := c.client.RawRequest("validateaddress", []json.RawMessage{input})
	if err != nil {
		return nil, fmt.Errorf("fail to call RawRequest(validateaddress): %w", err)
	}
	var result ValidateAddressResult
	if err := json.Unmarshal(rawResult, &result); err != nil {
		return nil, fmt.Errorf("fail to call json.Unmarshal(validateaddress): %w", err)
	}
	return &result, nil
}

// GetAddressesByLabel calls getaddressesbylabel and returns the raw wire map.
func (c *rpcClient) GetAddressesByLabel(labelName string) (map[string]Purpose, error) {
	input, err := json.Marshal(labelName)
	if err != nil {
		return nil, fmt.Errorf("fail to call json.Marshal(labelName): %w", err)
	}
	rawResult, err := c.client.RawRequest("getaddressesbylabel", []json.RawMessage{input})
	if err != nil {
		return nil, fmt.Errorf("fail to call RawRequest(getaddressesbylabel) %s: %w", labelName, err)
	}
	var result map[string]Purpose
	if err := json.Unmarshal(rawResult, &result); err != nil {
		return nil, fmt.Errorf("fail to call json.Unmarshal(getaddressesbylabel): %w", err)
	}
	return result, nil
}
