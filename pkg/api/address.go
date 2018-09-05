package api

import (
	"encoding/json"

	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcutil"
	"github.com/hiromaily/go-bitcoin/pkg/enum"
	"github.com/hiromaily/go-bitcoin/pkg/logger"
	"github.com/pkg/errors"
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

// CreateNewAddress アカウント名から新しいアドレスを生成する
// 常に新しいアドレスが生成される
// これによって作成されたアカウントはbitcoin core側のwalletで管理される
// TODO:おそらく本番では使わない
func (b *Bitcoin) CreateNewAddress(accountName string) (btcutil.Address, error) {
	addr, err := b.client.GetNewAddress(accountName)
	if err != nil {
		return nil, errors.Errorf("client.GetNewAddress(%s): error: %s", accountName, err)
	}

	return addr, nil
}

//GetAccountAddress アカウントに紐づくアドレスを返す
// Deprecated, will be removed in V0.18.
func (b *Bitcoin) GetAccountAddress(accountName string) (btcutil.Address, error) {
	if b.Version() >= enum.BTCVer17 {
		//TODO:複数件取得に変わるため、ゆくゆくは、この呼び出すはやめて、GetAddressesByLabel()を呼び出すようにする
		addrs, err := b.GetAddressesByLabel(accountName)
		if err != nil {
			return nil, errors.Errorf("BTC.GetAddressesByLabel() error: %s", err)
		}
		if len(addrs) == 0 {
			return nil, nil
		}
		return addrs[0], nil
	}
	return b.getAccountAddress(accountName)
}

// GetAccountAddress アカウントに紐づくアドレスを返す
// => アカウントに紐づくアドレスが無い場合は新規アドレスを作成し、そのアドレス値を返す
// => アカウントに紐づくアドレスが既にある場合は新規作成せずに既存アドレス値を返す。
// => 既存アドレスがある場合でもそれが使用済み（一度以上BTCを受け取った）の場合には新規アドレスを作成し、そのアドレス値を返す。
// Deprecated, will be removed in V0.18.
func (b *Bitcoin) getAccountAddress(accountName string) (btcutil.Address, error) {
	addr, err := b.client.GetAccountAddress(accountName)
	if err != nil {
		return nil, errors.Errorf("client.GetAccountAddress(%s): error: %s", accountName, err)
	}

	return addr, nil
}

// GetAddressesByLabel 指定したラベルに紐づくaddressをすべて返す
func (b *Bitcoin) GetAddressesByLabel(labelName string) ([]btcutil.Address, error) {
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
			logger.Errorf("b.DecodeAddress(%s) error: %s", key, err)
			continue
		}

		resAddrs[idx] = address
		idx++
	}

	return resAddrs, nil
}

// GetAddressesByAccount アカウント名から紐づくすべてのアドレスを取得する
func (b *Bitcoin) GetAddressesByAccount(accountName string) ([]btcutil.Address, error) {
	if b.Version() >= enum.BTCVer17 {
		//TODO:複数件取得に変わるため、ゆくゆくは、この呼び出すはやめて、GetAddressesByLabel()を呼び出すようにする
		addrs, err := b.GetAddressesByLabel(accountName)
		if err != nil {
			return nil, errors.Errorf("BTC.GetAddressesByLabel() error: %s", err)
		}
		return addrs, nil
	}
	return b.getAddressesByAccount(accountName)
}

// getAddressesByAccount アカウント名から紐づくすべてのアドレスを取得する
// Deprecated, will be removed in V0.18
func (b *Bitcoin) getAddressesByAccount(accountName string) ([]btcutil.Address, error) {
	//Returns the list of addresses for the given account.
	//アカウント名に紐づくアドレス一覧
	//[
	//	"2MvqroV1E2FjrRaWjQuLcfotZabhkSuYNGi",
	//	"2NGVf4fHRkWCtUBw4VXCtBMJdZBdgL7ffAq"
	//]

	addrs, err := b.client.GetAddressesByAccount(accountName)
	if err != nil {
		return nil, errors.Errorf("client.GetAddressesByAccount(%s): error: %s", accountName, err)
	}

	return addrs, nil
}

// ValidateAddress 渡されたアドレスの整合性をチェックする
// TODO: こちらの機能はCayenne側でも必要だが、Cayenneから直接利用する場合、Bitcoin Coreの機能に依存しているので、煩雑になってしまう
func (b *Bitcoin) ValidateAddress(addr string) (*btcjson.ValidateAddressWalletResult, error) {
	address, err := b.DecodeAddress(addr)
	if err != nil {
		return nil, errors.Errorf("DecodeAddress(%s): error: %s", addr, err)
	}
	res, err := b.client.ValidateAddress(address)
	if err != nil {
		return nil, errors.Errorf("client.ValidateAddress(%s): error: %s", addr, err)
	}
	//debug
	//type ValidateAddressWalletResult struct {
	//	IsValid      bool     `json:"isvalid"`
	//	Address      string   `json:"address,omitempty"`
	//	IsMine       bool     `json:"ismine,omitempty"`
	//	IsWatchOnly  bool     `json:"iswatchonly,omitempty"`
	//	IsScript     bool     `json:"isscript,omitempty"`
	//	PubKey       string   `json:"pubkey,omitempty"`
	//	IsCompressed bool     `json:"iscompressed,omitempty"`
	//	Account      string   `json:"account,omitempty"`
	//	Addresses    []string `json:"addresses,omitempty"`
	//	Hex          string   `json:"hex,omitempty"`
	//	Script       string   `json:"script,omitempty"`
	//	SigsRequired int32    `json:"sigsrequired,omitempty"`
	//}

	//b, _ := json.MarshalIndent(acc, "", " ")
	//logger.Info(string(b))

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
