package btc

import (
	"encoding/json"

	"github.com/btcsuite/btcutil"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

// GetAddressInfoResult is response type of PRC `getaddressinfo`
type GetAddressInfoResult struct {
	Address      string   `json:"address"`
	ScriptPubKey string   `json:"scriptPubKey"`
	Ismine       bool     `json:"ismine"`
	Solvable     bool     `json:"solvable"`
	Desc         string   `json:"desc"`
	Iswatchonly  bool     `json:"iswatchonly"`
	Isscript     bool     `json:"isscript"`
	Iswitness    bool     `json:"iswitness"`
	Pubkey       string   `json:"pubkey"`
	Iscompressed bool     `json:"iscompressed"`
	Ischange     bool     `json:"ischange"`
	Timestamp    int      `json:"timestamp"`
	Labels       []string `json:"labels"`
}

type ValidateAddressResult struct {
	IsValid           bool   `json:"isvalid"`
	Address           string `json:"address"`
	ScriptPubKey      string `json:"scriptPubKey"`
	IsScript          bool   `json:"isscript"`
	IsWitness         bool   `json:"iswitness"`
	WitnessVersion    int    `json:"witness_version,omitempty"`
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
	input, err := json.Marshal(string(addr))
	if err != nil {
		return nil, errors.Errorf("json.Marchal(): error: %s", err)
	}
	rawResult, err := b.client.RawRequest("getaddressinfo", []json.RawMessage{input})
	if err != nil {
		return nil, errors.Wrap(err, "fail to call json.RawRequest(getaddressinfo)")
	}

	infoResult := GetAddressInfoResult{}
	err = json.Unmarshal([]byte(rawResult), &infoResult)
	if err != nil {
		return nil, errors.Errorf("json.Unmarshal(): error: %s", err)
	}

	return &infoResult, nil
}

// GetAddressesByLabel returns addresses of account(label)
// Note: even if client has 5 addresses, it returns 15 addresses
//  it seems 3 different address types are returned respectively
// For now, it would be better to stop using it
func (b *Bitcoin) GetAddressesByLabel(labelName string) ([]btcutil.Address, error) {
	// input for rpc api
	input, err := json.Marshal(string(labelName))
	if err != nil {
		return nil, errors.Wrap(err, "fail to call json.Marchal()")
	}
	// call getaddressesbylabel
	rawResult, err := b.client.RawRequest("getaddressesbylabel", []json.RawMessage{input})
	if err != nil {
		return nil, errors.Wrap(err, "fail to call json.RawRequest(getaddressesbylabel)")
	}

	// unmarshal response
	var labels map[string]Purpose
	err = json.Unmarshal([]byte(rawResult), &labels)
	if err != nil {
		return nil, errors.Errorf("json.Unmarshal(): error: %s", err)
	}

	// retrieve
	resAddrs := make([]btcutil.Address, len(labels))
	idx := 0
	for key := range labels {
		//key is address string
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
	input, err := json.Marshal(string(addr))
	if err != nil {
		return nil, errors.Errorf("json.Marchal(): error: %s", err)
	}
	rawResult, err := b.client.RawRequest("validateaddress", []json.RawMessage{input})
	if err != nil {
		return nil, errors.Wrap(err, "fail to call json.RawRequest(validateaddress)")
	}

	result := ValidateAddressResult{}
	err = json.Unmarshal([]byte(rawResult), &result)
	if err != nil {
		return nil, errors.Errorf("json.Unmarshal(): error: %s", err)
	}
	if !result.IsValid{
		return nil, errors.Errorf("this address is invalid: %v", result)
	}

	return &result, nil
}

//func (b *Bitcoin) ValidateAddress(addr string) (*btcjson.ValidateAddressWalletResult, error) {
//	address, err := b.DecodeAddress(addr)
//	if err != nil {
//		return nil, errors.Wrapf(err, "fail to call btc.DecodeAddress(%s)", addr)
//	}
//	res, err := b.client.ValidateAddress(address)
//	if err != nil {
//		return nil, errors.Errorf("client.ValidateAddress(%s): error: %s", addr, err)
//	}
//
//	return res, nil
//}

// DecodeAddress decode string address to type Address
func (b *Bitcoin) DecodeAddress(addr string) (btcutil.Address, error) {
	address, err := btcutil.DecodeAddress(addr, b.chainConf)
	if err != nil {
		return nil, errors.Wrap(err, "fail to call btcutil.DecodeAddress()")
	}
	return address, nil
}
