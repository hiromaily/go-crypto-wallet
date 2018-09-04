package api

import (
	"github.com/pkg/errors"
)

//TODO:test用にtestアカウントを作成するのはいいかもしれない。

// GetAccount 渡されたアドレスから該当するアカウント名を取得する
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
