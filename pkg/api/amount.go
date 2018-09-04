package api

import (
	"github.com/btcsuite/btcutil"
	"github.com/pkg/errors"
	"strconv"
	"strings"
)

//0.00000001 BTC=1 Satoshi
//btcutil.Amount => Satoshi
//float64 => BTC

// AmountString 1.0 BTCといったstringを1.0という単位を除いたstringとして返す
func (b *Bitcoin) AmountString(amt btcutil.Amount) string {
	return strings.TrimRight(amt.String(), " BTC")
}

// FloatBitToAmount BTC(float64)をSatoshi(Amount)に変換する
// e.g. 0.54 to 54000000
func (b *Bitcoin) FloatBitToAmount(f float64) (btcutil.Amount, error) {
	amt, err := btcutil.NewAmount(f)
	if err != nil {
		return 0, errors.Errorf("btcutil.NewAmount(%f): error: %s", f, err)
	}

	return amt, nil
}

// CastStrBitToAmount string型のBitcoinをAmountに変換する
func (b *Bitcoin) CastStrBitToAmount(s string) (btcutil.Amount, error) {

	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, errors.Errorf("strconv.ParseFloat(%s): error: %s", s, err)
	}

	amt, err := btcutil.NewAmount(f)
	if err != nil {
		return 0, errors.Errorf("btcutil.NewAmount(%f): error: %s", f, err)
	}

	//return val.ToBTC(), nil
	return amt, nil
}

// CastStrSatoshiToAmount string型のSatoshiをAmountに変換する
func (b *Bitcoin) CastStrSatoshiToAmount(s string) (btcutil.Amount, error) {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, errors.Errorf("strconv.ParseFloat(%s): error: %s", s, err)
	}

	val, err := btcutil.NewAmount(f / float64(100000000))
	if err != nil {
		return 0, errors.Errorf("btcutil.NewAmount(%f): error: %s", f, err)
	}

	return val, nil
}

// ListAccounts これは単純にアカウントの資産一覧が表示されるだけ
func (b *Bitcoin) ListAccounts() (map[string]btcutil.Amount, error) {
	listAmts, err := b.client.ListAccounts()
	if err != nil {
		return nil, errors.Errorf("client.ListAccounts(): error: %s", err)
	}

	return listAmts, nil
}

// GetBalanceByAccount アカウントに対してのBalanceを取得する
func (b *Bitcoin) GetBalanceByAccount(accountName string) (btcutil.Amount, error) {
	amt, err := b.client.GetBalance(accountName)
	if err != nil {
		return 0, errors.Errorf("client.GetBalance(%s): error: %s", accountName, err)
	}

	return amt, nil
}

// GetBalanceByAccountAndMinConf アカウントに対してのBalanceを取得する
func (b *Bitcoin) GetBalanceByAccountAndMinConf(accountName string, minConf int) (btcutil.Amount, error) {
	amt, err := b.client.GetBalanceMinConf(accountName, minConf)
	if err != nil {
		return 0, errors.Errorf("client.GetBalanceMinConf(%s): error: %s", accountName, err)
	}

	return amt, nil
}
