package api

import (
	"github.com/btcsuite/btcutil"
	"github.com/pkg/errors"
)

// CreateNewAddress アカウント名から新しいアドレスを生成する
// これによって作成されたアカウントはbitcoin core側のwalletで管理される
func (b *Bitcoin) CreateNewAddress(accountName string) (btcutil.Address, error) {
	addr, err := b.Client.GetNewAddress(accountName)
	if err != nil {
		return nil, errors.Errorf("GetNewAddress(%s): error: %v", accountName, err)
	}

	return addr, nil
}

// GetAddressesByAccount アカウント名から紐づくすべてのアドレスを取得する
func (b *Bitcoin) GetAddressesByAccount(accountName string) ([]btcutil.Address, error) {
	addrs, err := b.Client.GetAddressesByAccount(accountName)
	if err != nil {
		return nil, errors.Errorf("GetAddressesByAccount(%s): error: %v", accountName, err)
	}

	return addrs, nil
}
