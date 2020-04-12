package btc

import (
	"encoding/json"

	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcutil"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

// GetAddressInfoResult getaddressinfoをcallしたresponseの型
type GetAddressInfoResult struct {
	Address      string `json:"address"`
	ScriptPubKey string `json:"scriptPubKey"`
	Ismine       bool   `json:"ismine"`
	Iswatchonly  bool   `json:"iswatchonly"`
	Isscript     bool   `json:"isscript"`
	Iswitness    bool   `json:"iswitness"`
	Script       string `json:"script"`
	Hex          string `json:"hex"`
	Pubkey       string `json:"pubkey"`
	Embedded     struct {
		Isscript       bool   `json:"isscript"`
		Iswitness      bool   `json:"iswitness"`
		WitnessVersion int    `json:"witness_version"`
		WitnessProgram string `json:"witness_program"`
		Pubkey         string `json:"pubkey"`
		Address        string `json:"address"`
		ScriptPubKey   string `json:"scriptPubKey"`
	} `json:"embedded"`
	Label     string  `json:"label"`
	Timestamp int64   `json:"timestamp"`
	Labels    []Label `json:"labels"`
}

// Label ラベル
type Label struct {
	Name    string `json:"name"`
	Purpose string `json:"purpose"`
}

// Purpose 目的
type Purpose struct {
	Purpose string `json:"purpose"`
}

// GetAddressInfo getaddressinfo RPC をcallする
// version0.18より、getaccountは呼び出せなくなるので、こちらをcallすること
// 従来のvalidateaddressより取得していたaddressの詳細情報もこちらから取得可能
func (b *Bitcoin) GetAddressInfo(addr string) (*GetAddressInfoResult, error) {
	input, err := json.Marshal(string(addr))
	if err != nil {
		return nil, errors.Errorf("json.Marchal(): error: %s", err)
	}
	rawResult, err := b.client.RawRequest("getaddressinfo", []json.RawMessage{input})
	if err != nil {
		return nil, errors.Errorf("json.RawRequest(getaddressinfo): error: %s", err)
	}

	infoResult := GetAddressInfoResult{}
	err = json.Unmarshal([]byte(rawResult), &infoResult)
	if err != nil {
		return nil, errors.Errorf("json.Unmarshal(): error: %s", err)
	}

	return &infoResult, nil
}

// GetAddressesByLabel 指定したラベルに紐づくaddressをすべて返す
func (b *Bitcoin) GetAddressesByLabel(labelName string) ([]btcutil.Address, error) {
	//if b.Version() >= ctype.BTCVer17 {
	input, err := json.Marshal(string(labelName))
	if err != nil {
		return nil, errors.Errorf("json.Marchal(): error: %s", err)
	}
	rawResult, err := b.client.RawRequest("getaddressesbylabel", []json.RawMessage{input})
	if err != nil {
		return nil, errors.Errorf("json.RawRequest(getaddressesbylabel): error: %s", err)
	}

	var labels map[string]Purpose
	err = json.Unmarshal([]byte(rawResult), &labels)
	if err != nil {
		return nil, errors.Errorf("json.Unmarshal(): error: %s", err)
	}

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

// ValidateAddress 渡されたアドレスの整合性をチェックする
func (b *Bitcoin) ValidateAddress(addr string) (*btcjson.ValidateAddressWalletResult, error) {
	address, err := b.DecodeAddress(addr)
	if err != nil {
		return nil, errors.Errorf("DecodeAddress(%s): error: %s", addr, err)
	}
	res, err := b.client.ValidateAddress(address)
	if err != nil {
		return nil, errors.Errorf("client.ValidateAddress(%s): error: %s", addr, err)
	}

	return res, nil
}

// DecodeAddress string型のアドレスをDecodeしてAddress型に変換する
func (b *Bitcoin) DecodeAddress(addr string) (btcutil.Address, error) {
	address, err := btcutil.DecodeAddress(addr, b.chainConf)
	if err != nil {
		return nil, errors.Errorf("btcutil.DecodeAddress() error: %s", err)
	}
	return address, nil
}
