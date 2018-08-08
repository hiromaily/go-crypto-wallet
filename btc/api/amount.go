package api

import (
	"github.com/btcsuite/btcutil"
	"github.com/pkg/errors"
)

// ConvertToSatoshi bitcoinをSatoshiに変換する
func ConvertToSatoshi(f float64) (btcutil.Amount, error) {
	// Amount
	// Satoshiに変換しないといけない
	// 1Satoshi＝0.00000001BTC
	amt, err := btcutil.NewAmount(f)
	if err != nil {
		return 0, errors.Errorf("btcutil.NewAmount(%f): error: %v", f, err)
	}

	return amt, nil
}

// ListAccounts これは単純にアカウントの資産一覧が表示されるだけ
func (b *Bitcoin) ListAccounts() (map[string]btcutil.Amount, error) {
	listAmts, err := b.Client.ListAccounts()
	if err != nil {
		return nil, errors.Errorf("ListAccounts(): error: %v", err)
	}

	return listAmts, nil
}

// GetBalanceByAccount アカウントに対してのBalanceを取得する
func (b *Bitcoin) GetBalanceByAccount(accountName string) (btcutil.Amount, error) {
	amt, err := b.Client.GetBalance(accountName)
	if err != nil {
		return 0, errors.Errorf("GetBalance(%s): error: %v", accountName, err)
	}

	return amt, nil
}
