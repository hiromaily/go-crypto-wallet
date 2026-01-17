package bch

import (
	"encoding/json"
	"fmt"

	"github.com/btcsuite/btcd/btcutil"

	dtobtc "github.com/hiromaily/go-crypto-wallet/internal/application/dto/btc"
	apibtcimpl "github.com/hiromaily/go-crypto-wallet/internal/infrastructure/api/btc/btc"
	bchutil "github.com/hiromaily/go-crypto-wallet/pkg/cryptocurrency/bch"
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
func (b *BitcoinCash) GetAddressInfo(addr string) (*dtobtc.AddressInfo, error) {
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

	// convert bch result to btc infrastructure type, then to application DTO
	btcResult := &apibtcimpl.GetAddressInfoResult{
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
	}

	return apibtcimpl.ToAddressInfo(btcResult), nil
}

// GetAddressesByLabel overrides Bitcoin's GetAddressesByLabel to use BCH address decoding
func (b *BitcoinCash) GetAddressesByLabel(labelName string) ([]btcutil.Address, error) {
	// This override is necessary because the parent Bitcoin.GetAddressesByLabel calls
	// Bitcoin.DecodeAddress internally, which doesn't understand BCH CashAddr format.
	logger.Debug("getting addresses by label", "label", labelName)

	// input for rpc api
	input, err := json.Marshal(labelName)
	if err != nil {
		return nil, fmt.Errorf("fail to call json.Marchal() in bch: %w", err)
	}
	// call getaddressesbylabel
	rawResult, err := b.Client.RawRequest("getaddressesbylabel", []json.RawMessage{input})
	if err != nil {
		logger.Debug("getaddressesbylabel RPC failed", "label", labelName, "error", err)
		return nil, fmt.Errorf("fail to call json.RawRequest(getaddressesbylabel) for label %s in bch: %w", labelName, err)
	}

	// unmarshal response
	var labels map[string]apibtcimpl.Purpose
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
	resAddrs := make([]btcutil.Address, 0, len(labels))
	for key := range labels {
		// key is address string
		var address btcutil.Address
		address, err = b.DecodeAddress(key)
		if err != nil {
			logger.Error(
				"fail to call b.DecodeAddress()",
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
