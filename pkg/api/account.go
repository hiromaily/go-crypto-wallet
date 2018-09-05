package api

import (
	"github.com/hiromaily/go-bitcoin/pkg/enum"
	"github.com/pkg/errors"
)

// GetAccount 渡されたアドレスから該当するアカウント名を取得する
func (b *Bitcoin) GetAccount(addr string) (string, error) {
	if b.Version() >= enum.BTCVer17 {
		res, err := b.GetAddressInfo(addr)
		if err != nil {
			return "", errors.Errorf("BTC.GetAddressInfo() error: %s", err)
		}
		return res.Label, nil
	}
	return b.getAccount(addr)
}

// GetAccount 渡されたアドレスから該当するアカウント名を取得する
// version0.18より、getaccountは呼び出せなくなるので、GetAddressInfo()をcallすること
func (b *Bitcoin) getAccount(addr string) (string, error) {
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
func (b *Bitcoin) SetAccount(addr, account string) error {
	if b.Version() >= enum.BTCVer17 {
		err := b.SetLabel(addr, account)
		if err != nil {
			return errors.Errorf("BTC.SetLabel() error: %s", err)
		}
		return nil
	}
	return b.setAccount(addr, account)
}

// SetAccount 既存のimport済のアドレスにアカウント名をセットする
// version0.18より、setaccountは呼び出せなくなるので、SetLabel()をcallすること
func (b *Bitcoin) setAccount(addr, account string) error {
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
