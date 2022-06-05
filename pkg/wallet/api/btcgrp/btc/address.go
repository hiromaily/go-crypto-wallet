package btc

import (
	"encoding/json"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

// GetAddressInfoResult is response type of RPC `getaddressinfo`
type GetAddressInfoResult struct {
	Address      string   `json:"address"`
	ScriptPubKey string   `json:"scriptPubKey"`
	Ismine       bool     `json:"ismine"`
	Solvable     bool     `json:"solvable,omitempty"`
	Desc         string   `json:"desc,omitempty"`
	Iswatchonly  bool     `json:"iswatchonly"`
	Isscript     bool     `json:"isscript"`
	Iswitness    bool     `json:"iswitness,omitempty"`
	Pubkey       string   `json:"pubkey,omitempty"`
	Iscompressed bool     `json:"iscompressed,omitempty"`
	Ischange     bool     `json:"ischange"`
	Timestamp    int64    `json:"timestamp,omitempty"`
	Labels       []string `json:"labels"`
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
func (b *Bitcoin) GetAddressInfo(addr string) (*GetAddressInfoResult, error) {
	input, err := json.Marshal(addr)
	if err != nil {
		return nil, errors.Wrap(err, "fail to call json.Marchal()")
	}
	rawResult, err := b.Client.RawRequest("getaddressinfo", []json.RawMessage{input})
	if err != nil {
		return nil, errors.Wrapf(err, "fail to call json.RawRequest(getaddressinfo) %s", addr)
	}

	infoResult := GetAddressInfoResult{}
	err = json.Unmarshal(rawResult, &infoResult)
	if err != nil {
		return nil, errors.Wrap(err, "fail to call json.Unmarshal(rawResult)")
	}

	return &infoResult, nil
}

// GetAddressesByLabel returns addresses of account(label)
// Note: even if client has 5 addresses, it returns 15 addresses
//  it seems 3 different address types are returned respectively
// For now, it would be better to stop using it
func (b *Bitcoin) GetAddressesByLabel(labelName string) ([]btcutil.Address, error) {
	// input for rpc api
	input, err := json.Marshal(labelName)
	if err != nil {
		return nil, errors.Wrap(err, "fail to call json.Marchal()")
	}
	// call getaddressesbylabel
	rawResult, err := b.Client.RawRequest("getaddressesbylabel", []json.RawMessage{input})
	if err != nil {
		return nil, errors.Wrap(err, "fail to call json.RawRequest(getaddressesbylabel)")
	}

	// unmarshal response
	var labels map[string]Purpose
	err = json.Unmarshal(rawResult, &labels)
	if err != nil {
		return nil, errors.Wrap(err, "fail to call json.Unmarshal(rawResult)")
	}

	// retrieve
	resAddrs := make([]btcutil.Address, len(labels))
	idx := 0
	for key := range labels {
		// key is address string
		address, err := b.DecodeAddress(key)
		if err != nil {
			b.logger.Error(
				"fail to call b.DecodeAddress()",
				zap.String("address", key),
				zap.Error(err))
			continue
		}

		resAddrs[idx] = address
		idx++
	}

	return resAddrs, nil
}

// ValidateAddress validate address
func (b *Bitcoin) ValidateAddress(addr string) (*ValidateAddressResult, error) {
	input, err := json.Marshal(addr)
	if err != nil {
		return nil, errors.Errorf("json.Marchal(): error: %s", err)
	}
	rawResult, err := b.Client.RawRequest("validateaddress", []json.RawMessage{input})
	if err != nil {
		return nil, errors.Wrap(err, "fail to call json.RawRequest(validateaddress)")
	}

	result := ValidateAddressResult{}
	err = json.Unmarshal(rawResult, &result)
	if err != nil {
		return nil, errors.Errorf("json.Unmarshal(): error: %s", err)
	}
	if !result.IsValid {
		return nil, errors.Errorf("this address is invalid: %v", result)
	}

	return &result, nil
}

// DecodeAddress decode string address to type Address
func (b *Bitcoin) DecodeAddress(addr string) (btcutil.Address, error) {
	address, err := btcutil.DecodeAddress(addr, b.chainConf)
	if err != nil {
		return nil, errors.Wrap(err, "fail to call btcutil.DecodeAddress()")
	}
	return address, nil
}
