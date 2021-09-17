package bch

import (
	"encoding/json"

	"github.com/pkg/errors"

	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/api/btcgrp/btc"
)

// GetAddressInfoResult is response type of RPC `getaddressinfo`
type GetAddressInfoResult struct {
	Address      string `json:"address"`
	ScriptPubKey string `json:"scriptPubKey"`
	Ismine       bool   `json:"ismine"`
	Iswatchonly  bool   `json:"iswatchonly"`
	Isscript     bool   `json:"isscript"`
	Pubkey       string `json:"pubkey,omitempty"`
	Iscompressed bool   `json:"iscompressed,omitempty"`
	Label        string `json:"label,omitempty"`
	Ischange     bool   `json:"ischange"`
	Timestamp    int64  `json:"timestamp,omitempty"`
	Labels       []struct {
		Name    string `json:"name"`
		Purpose string `json:"purpose"`
	} `json:"labels"`
}

// GetAddressInfo can be used as an alternative to `getaccount`, `validateaddress`
func (b *BitcoinCash) GetAddressInfo(addr string) (*btc.GetAddressInfoResult, error) {
	input, err := json.Marshal(addr)
	if err != nil {
		return nil, errors.Wrap(err, "fail to call json.Marchal() in bch")
	}
	rawResult, err := b.Client.RawRequest("getaddressinfo", []json.RawMessage{input})
	if err != nil {
		return nil, errors.Wrapf(err, "fail to call json.RawRequest(getaddressinfo) %s in bch", addr)
	}

	infoResult := GetAddressInfoResult{}
	err = json.Unmarshal(rawResult, &infoResult)
	if err != nil {
		return nil, errors.Wrap(err, "fail to call json.Unmarshal(rawResult) in bch")
	}

	// convert bch result to btc
	return &btc.GetAddressInfoResult{
		Address:      infoResult.Address,
		ScriptPubKey: infoResult.ScriptPubKey,
		Ismine:       infoResult.Ismine,
		Solvable:     false,
		Desc:         "",
		Iswatchonly:  infoResult.Iswatchonly,
		Isscript:     infoResult.Isscript,
		Iswitness:    false,
		Pubkey:       infoResult.Pubkey,
		Iscompressed: infoResult.Iscompressed,
		Ischange:     infoResult.Ischange,
		Timestamp:    infoResult.Timestamp,
		Labels:       []string{infoResult.Label},
	}, nil
}
