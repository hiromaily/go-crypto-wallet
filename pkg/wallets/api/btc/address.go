package btc

import (
	"encoding/json"

	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcutil"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

// GetAddressInfoResult stores response of PRC `getaddressinfo`
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

func (a *GetAddressInfoResult) GetLabelName() string {
	if len(a.Labels) != 0 {
		return a.Labels[0]
	}
	return ""
}

// GetAddressInfo call RPC `getaddressinfo`
//{
//  "address": "mvTRCKpKVUUv3QgMEn838xXDDZS5SSEhnj",
//  "scriptPubKey": "76a914a3deadcdee77d544692dfc64eb321cccc9e036f188ac",
//  "ismine": true,
//  "solvable": true,
//  "desc": "pkh([a3deadcd]02f4d649c24780191d31d4fa23bff91f3fb2646b47d7ef32714e5322059586765e)#q4e9d52m",
//  "iswatchonly": false,
//  "isscript": false,
//  "iswitness": false,
//  "pubkey": "02f4d649c24780191d31d4fa23bff91f3fb2646b47d7ef32714e5322059586765e",
//  "iscompressed": true,
//  "ischange": false,
//  "timestamp": 1,
//  "labels": [
//    "client"
//  ]
//}

// Purpose stores part of response of PRC `getaddressesbylabel`
type Purpose struct {
	Purpose string `json:"purpose"`
}

//{
//  "mvTRCKpKVUUv3QgMEn838xXDDZS5SSEhnj": {
//    "purpose": "receive"
//  },
//  "2MwYRJrBZ4fqAbdME7uCfRisyp3Mp8ooP6Y": {
//    "purpose": "receive"
//  },
//  "tb1q5002mn0wwl25g6fdl3jwkvsueny7qdh3a7670e": {
//    "purpose": "receive"
//  }
//}

// it can be used as an alternative to `getaccount`, `validateaddress`
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
func (b *Bitcoin) GetAddressesByLabel(labelName string) ([]btcutil.Address, error) {
	//if b.Version() >= ctype.BTCVer17 {
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
func (b *Bitcoin) ValidateAddress(addr string) (*btcjson.ValidateAddressWalletResult, error) {
	address, err := b.DecodeAddress(addr)
	if err != nil {
		return nil, errors.Wrapf(err, "fail to call btc.DecodeAddress(%s)", addr)
	}
	res, err := b.client.ValidateAddress(address)
	if err != nil {
		return nil, errors.Errorf("client.ValidateAddress(%s): error: %s", addr, err)
	}

	return res, nil
}

// DecodeAddress decode string address to type Address
func (b *Bitcoin) DecodeAddress(addr string) (btcutil.Address, error) {
	address, err := btcutil.DecodeAddress(addr, b.chainConf)
	if err != nil {
		return nil, errors.Wrap(err, "fail to call btcutil.DecodeAddress()")
	}
	return address, nil
}
