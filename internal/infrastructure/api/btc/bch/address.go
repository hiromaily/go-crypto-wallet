package bch

import (
	"encoding/json"
	"fmt"

	"github.com/btcsuite/btcd/btcutil"

	bchutil "github.com/hiromaily/go-crypto-wallet/pkg/chains/bch"
	btcrpc "github.com/hiromaily/go-crypto-wallet/pkg/chains/btc/rpc"
	"github.com/hiromaily/go-crypto-wallet/pkg/logger"
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
func (b *BitcoinCash) GetAddressInfo(addr string) (*btcrpc.GetAddressInfoResult, error) {
	input, err := json.Marshal(addr)
	if err != nil {
		return nil, fmt.Errorf("fail to call json.Marchal() in bch: %w", err)
	}
	rawResult, err := b.Client.RawRequest("getaddressinfo", []json.RawMessage{input})
	if err != nil {
		return nil, fmt.Errorf("fail to call json.RawRequest(getaddressinfo) %s in bch: %w", addr, err)
	}

	infoResult := GetAddressInfoResult{}
	err = json.Unmarshal(rawResult, &infoResult)
	if err != nil {
		return nil, fmt.Errorf("fail to call json.Unmarshal(rawResult) in bch: %w", err)
	}

	return &btcrpc.GetAddressInfoResult{
		Address:      infoResult.Address,
		ScriptPubKey: infoResult.ScriptPubKey,
		IsMine:       infoResult.Ismine,
		IsWatchOnly:  infoResult.Iswatchonly,
		IsScript:     infoResult.Isscript,
		Pubkey:       infoResult.Pubkey,
		IsCompressed: infoResult.Iscompressed,
		IsChange:     infoResult.Ischange,
		Timestamp:    infoResult.Timestamp,
		Labels:       btcrpc.FlexibleLabels{infoResult.Label},
	}, nil
}

// GetAddressesByLabel overrides Bitcoin's GetAddressesByLabel to use BCH address decoding
func (b *BitcoinCash) GetAddressesByLabel(labelName string) ([]btcutil.Address, error) {
	// This override is necessary because the parent Bitcoin.GetAddressesByLabel calls
	// Bitcoin.DecodeAddress internally, which doesn't understand BCH CashAddr format.
	logger.Debug("BCH GetAddressesByLabel called", "label", labelName)

	// input for rpc api
	input, err := json.Marshal(labelName)
	if err != nil {
		return nil, fmt.Errorf("fail to call json.Marchal() in bch: %w", err)
	}
	// call getaddressesbylabel
	rawResult, err := b.Client.RawRequest("getaddressesbylabel", []json.RawMessage{input})
	if err != nil {
		logger.Debug("getaddressesbylabel RPC failed", "label", labelName, "error", err)
		return nil, fmt.Errorf(
			"fail to call json.RawRequest(getaddressesbylabel) for label %s in bch: %w",
			labelName, err)
	}

	// unmarshal response - only keys (address strings) are needed
	var labels map[string]json.RawMessage
	err = json.Unmarshal(rawResult, &labels)
	if err != nil {
		return nil, fmt.Errorf("fail to call json.Unmarshal(rawResult) in bch: %w", err)
	}

	if len(labels) == 0 {
		logger.Debug("no addresses found for label", "label", labelName)
		return nil, nil
	}

	logger.Debug("found addresses for label",
		"label", labelName,
		"raw_count", len(labels))

	// retrieve - use BCH DecodeAddress for proper CashAddr handling
	// With GetBalanceByAccount override, b is *BitcoinCash so b.DecodeAddress calls BCH version
	resAddrs := make([]btcutil.Address, 0, len(labels))
	for key := range labels {
		// key is address string
		address, err := b.DecodeAddress(key)
		if err != nil {
			logger.Error(
				"failed to decode address",
				"address", key,
				"label", labelName,
				"error", err)
			continue
		}

		resAddrs = append(resAddrs, address)
	}

	logger.Debug("successfully decoded addresses",
		"label", labelName,
		"decoded_count", len(resAddrs))

	return resAddrs, nil
}

// DecodeAddress overrides Bitcoin's DecodeAddress to handle BCH CashAddr format
func (b *BitcoinCash) DecodeAddress(addr string) (btcutil.Address, error) {
	// Use BCH-specific decoder that handles CashAddr format
	logger.Debug("BCH DecodeAddress called", "address", addr)
	address, err := bchutil.DecodeAddress(addr, b.GetChainConf())
	if err != nil {
		return nil, fmt.Errorf("fail to call bchutil.DecodeAddress() in bch: %w", err)
	}
	return address, nil
}
