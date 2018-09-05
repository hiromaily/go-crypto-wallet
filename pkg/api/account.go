package api

import (
	"github.com/pkg/errors"
)

// GetAccount 渡されたアドレスから該当するアカウント名を取得する
// version0.18より、getaccountは呼び出せなくなるので、GetAddressInfo()をcallすること
func (b *Bitcoin) GetAccount(addr string) (string, error) {
	address, err := b.DecodeAddress(addr)
	if err != nil {
		return "", errors.Errorf("DecodeAddress(%s): error: %s", addr, err)
	}

	accountName, err := b.client.GetAccount(address)
	if err != nil {
		return "", errors.Errorf("client.GetAccount(): error: %s", err)
	}

	return accountName, nil
}

// SetAccount 既存のimport済のアドレスにアカウント名をセットする
// version0.18より、setaccountは呼び出せなくなるので、SetLabel()をcallすること
func (b *Bitcoin) SetAccount(addr, account string) error {
	address, err := b.DecodeAddress(addr)
	if err != nil {
		return errors.Errorf("DecodeAddress(%s): error: %s", addr, err)
	}

	err = b.client.SetAccount(address, account)
	if err != nil {
		return errors.Errorf("client.SetAccount(): error: %s", err)
	}

	return nil
}
