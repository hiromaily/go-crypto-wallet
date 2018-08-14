package api

import (
	"github.com/btcsuite/btcutil"
	"github.com/pkg/errors"
)

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
func (b *Bitcoin) ValidateAddress(addr string) error {
	address, err := b.DecodeAddress(addr)
	if err != nil {
		return errors.Errorf("DecodeAddress(%s): error: %v", addr, err)
	}
	_, err = b.client.ValidateAddress(address)
	if err != nil {
		return errors.Errorf("client.ValidateAddress(%s): error: %v", addr, err)
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
	//log.Println(string(b))

	return nil
}

// DecodeAddress string型のアドレスをDecodeしてAddress型に変換する
func (b *Bitcoin) DecodeAddress(addr string) (btcutil.Address, error) {
	address, err := btcutil.DecodeAddress(addr, b.chainConf)
	if err != nil {
		return nil, errors.Errorf("btcutil.DecodeAddress() error: %v", err)
	}
	return address, nil
}
