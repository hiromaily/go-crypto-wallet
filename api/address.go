package api

import (
	"github.com/btcsuite/btcutil"
)

// CreateNewAddress アカウント名から新しいアドレスを生成する
// これによって作成されたアカウントはbitcoin core側のwalletで管理される
func (b *Bitcoin) CreateNewAddress(accountName string) (btcutil.Address, error) {
	return b.Client.GetNewAddress(accountName)
}

// GetAddressesByAccount アカウント名から紐づくすべてのアドレスを取得する
func (b *Bitcoin) GetAddressesByAccount(name string) ([]btcutil.Address, error) {
	return b.Client.GetAddressesByAccount(name)
}
