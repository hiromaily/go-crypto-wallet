package api

import (
	"encoding/json"

	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcutil"
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
	Label     string `json:"label"`
	Timestamp int64  `json:"timestamp"`
	Labels    []struct {
		Name    string `json:"name"`
		Purpose string `json:"purpose"`
	} `json:"labels"`
}

// CreateNewAddress アカウント名から新しいアドレスを生成する
// 常に新しいアドレスが生成される
// これによって作成されたアカウントはbitcoin core側のwalletで管理される
// TODO:おそらく本番では使わない
func (b *Bitcoin) CreateNewAddress(accountName string) (btcutil.Address, error) {
	addr, err := b.client.GetNewAddress(accountName)
	if err != nil {
		return nil, errors.Errorf("client.GetNewAddress(%s): error: %v", accountName, err)
	}

	return addr, nil
}

// GetAccountAddress アカウントに紐づくアドレスを返す
// => アカウントに紐づくアドレスが無い場合は新規アドレスを作成し、そのアドレス値を返す
// => アカウントに紐づくアドレスが既にある場合は新規作成せずに既存アドレス値を返す。
// => 既存アドレスがある場合でもそれが使用済み（一度以上BTCを受け取った）の場合には新規アドレスを作成し、そのアドレス値を返す。
func (b *Bitcoin) GetAccountAddress(accountName string) (btcutil.Address, error) {
	addr, err := b.client.GetAccountAddress(accountName)
	if err != nil {
		return nil, errors.Errorf("client.GetAccountAddress(%s): error: %v", accountName, err)
	}

	return addr, nil
}

// GetAddressesByAccount アカウント名から紐づくすべてのアドレスを取得する
// TODO:本番でaddressに対してアカウント名を紐付けることをするかどうか未定
func (b *Bitcoin) GetAddressesByAccount(accountName string) ([]btcutil.Address, error) {
	//Returns the list of addresses for the given account.
	//アカウント名に紐づくアドレス一覧
	//[
	//	"2MvqroV1E2FjrRaWjQuLcfotZabhkSuYNGi",
	//	"2NGVf4fHRkWCtUBw4VXCtBMJdZBdgL7ffAq"
	//]

	addrs, err := b.client.GetAddressesByAccount(accountName)
	if err != nil {
		return nil, errors.Errorf("client.GetAddressesByAccount(%s): error: %v", accountName, err)
	}

	return addrs, nil
}

// ValidateAddress 渡されたアドレスの整合性をチェックする
// TODO: こちらの機能はCayenne側でも必要だが、Cayenneから直接利用する場合、Bitcoin Coreの機能に依存しているので、煩雑になってしまう
func (b *Bitcoin) ValidateAddress(addr string) (*btcjson.ValidateAddressWalletResult, error) {
	address, err := b.DecodeAddress(addr)
	if err != nil {
		return nil, errors.Errorf("DecodeAddress(%s): error: %v", addr, err)
	}
	res, err := b.client.ValidateAddress(address)
	if err != nil {
		return nil, errors.Errorf("client.ValidateAddress(%s): error: %v", addr, err)
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
		return nil, errors.Errorf("btcutil.DecodeAddress() error: %v", err)
	}
	return address, nil
}

// GetAddressInfo getaddressinfo RPC をcallする
// version0.18より、getaccountは呼び出せなくなるので、こちらをcallすること
// 従来のvalidateaddressより取得していたaddressの詳細情報もこちらから取得可能
func (b *Bitcoin) GetAddressInfo(addr string) (*GetAddressInfoResult, error) {
	input, err := json.Marshal(string(addr))
	if err != nil {
		return nil, errors.Errorf("json.Marchal(): error: %v", err)
	}
	rawResult, err := b.client.RawRequest("getaddressinfo", []json.RawMessage{input})
	if err != nil {
		return nil, errors.Errorf("json.RawRequest(getaddressinfo): error: %v", err)
	}

	infoResult := GetAddressInfoResult{}
	err = json.Unmarshal([]byte(rawResult), &infoResult)
	if err != nil {
		return nil, errors.Errorf("json.Unmarshal(): error: %v", err)
	}

	return &infoResult, nil
}
