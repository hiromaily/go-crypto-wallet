package btc

import (
	"encoding/json"
	"fmt"

	"github.com/btcsuite/btcd/btcutil"

	dtobtc "github.com/hiromaily/go-crypto-wallet/internal/application/dto/btc"
	"github.com/hiromaily/go-crypto-wallet/pkg/logger"
)

// GetAddressInfoResult is response type of RPC `getaddressinfo`
type GetAddressInfoResult struct {
	Address              string   `json:"address"`
	ScriptPubKey         string   `json:"scriptPubKey"`
	Ismine               bool     `json:"ismine"`
	Solvable             bool     `json:"solvable,omitempty"`
	Desc                 string   `json:"desc,omitempty"`
	Iswatchonly          bool     `json:"iswatchonly"`
	Isscript             bool     `json:"isscript"`
	Iswitness            bool     `json:"iswitness,omitempty"`
	Pubkey               string   `json:"pubkey,omitempty"`
	Iscompressed         bool     `json:"iscompressed,omitempty"`
	Ischange             bool     `json:"ischange"`
	Timestamp            int64    `json:"timestamp,omitempty"`
	HDKeyPath            string   `json:"hdkeypath,omitempty"`
	HDSeedID             string   `json:"hdseedid,omitempty"`
	HDMasterFingerprint  string   `json:"hdmasterfingerprint,omitempty"`
	Labels               []string `json:"labels"`
}

// ValidateAddressResult is response type of RPC `validateaddress`
type ValidateAddressResult struct {
	IsValid           bool   `json:"isvalid"`
	Address           string `json:"address"`
	ScriptPubKey      string `json:"scriptPubKey"`
	IsScript          bool   `json:"isscript"`
	IsWitness         bool   `json:"iswitness"`
	WitnessVersion    uint32 `json:"witness_version,omitempty"`
	WitnessProgramHex string `json:"witness_program,omitempty"`
}

// GetLabelName returns label name
func (a *GetAddressInfoResult) GetLabelName() string {
	if len(a.Labels) != 0 {
		return a.Labels[0]
	}
	return ""
}

// Purpose stores part of response of PRC `getaddressesbylabel`
type Purpose struct {
	Purpose string `json:"purpose"`
}

// GetAddressInfo can be used as an alternative to `getaccount`, `validateaddress`
func (b *Bitcoin) GetAddressInfo(addr string) (*dtobtc.AddressInfo, error) {
	input, err := json.Marshal(addr)
	if err != nil {
		return nil, fmt.Errorf("fail to call json.Marchal(): %w", err)
	}
	rawResult, err := b.Client.RawRequest("getaddressinfo", []json.RawMessage{input})
	if err != nil {
		return nil, fmt.Errorf("fail to call json.RawRequest(getaddressinfo) %s: %w", addr, err)
	}

	infoResult := GetAddressInfoResult{}
	err = json.Unmarshal(rawResult, &infoResult)
	if err != nil {
		return nil, fmt.Errorf("fail to call json.Unmarshal(rawResult): %w", err)
	}

	return ToAddressInfo(&infoResult), nil
}

// GetAddressesByLabel returns addresses of account(label)
// Note: even if client has 5 addresses, it returns 15 addresses
//
//	it seems 3 different address types are returned respectively
//
// For now, it would be better to stop using it
func (b *Bitcoin) GetAddressesByLabel(labelName string) ([]btcutil.Address, error) {
	logger.Debug("getting addresses by label", "label", labelName)

	// input for rpc api
	input, err := json.Marshal(labelName)
	if err != nil {
		return nil, fmt.Errorf("fail to call json.Marchal(): %w", err)
	}
	// call getaddressesbylabel
	rawResult, err := b.Client.RawRequest("getaddressesbylabel", []json.RawMessage{input})
	if err != nil {
		logger.Debug("getaddressesbylabel RPC failed", "label", labelName, "error", err)
		return nil, fmt.Errorf("fail to call json.RawRequest(getaddressesbylabel) for label %s: %w", labelName, err)
	}

	// unmarshal response
	var labels map[string]Purpose
	err = json.Unmarshal(rawResult, &labels)
	if err != nil {
		return nil, fmt.Errorf("fail to call json.Unmarshal(rawResult): %w", err)
	}

	if len(labels) == 0 {
		logger.Debug("no addresses found for label", "label", labelName)
		return nil, nil
	}

	logger.Debug("found addresses for label",
		"label", labelName,
		"raw_count", len(labels))

	// retrieve
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

// ValidateAddress validate address
func (b *Bitcoin) ValidateAddress(addr string) (*dtobtc.ValidateAddressResult, error) {
	input, err := json.Marshal(addr)
	if err != nil {
		return nil, fmt.Errorf("json.Marchal(): error: %s", err)
	}
	rawResult, err := b.Client.RawRequest("validateaddress", []json.RawMessage{input})
	if err != nil {
		return nil, fmt.Errorf("fail to call json.RawRequest(validateaddress): %w", err)
	}

	result := ValidateAddressResult{}
	err = json.Unmarshal(rawResult, &result)
	if err != nil {
		return nil, fmt.Errorf("json.Unmarshal(): error: %s", err)
	}
	if !result.IsValid {
		return nil, fmt.Errorf("this address is invalid: %v", result)
	}

	return ToValidateAddressResult(&result), nil
}

// DecodeAddress decode string address to type Address
func (b *Bitcoin) DecodeAddress(addr string) (btcutil.Address, error) {
	address, err := btcutil.DecodeAddress(addr, b.chainConf)
	if err != nil {
		return nil, fmt.Errorf("fail to call btcutil.DecodeAddress(): %w", err)
	}
	return address, nil
}
