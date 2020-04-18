package btc

import (
	"strconv"
	"strings"

	"github.com/btcsuite/btcutil"
	"github.com/pkg/errors"
)

//0.00000001 BTC=1 Satoshi
//btcutil.Amount => Satoshi
//float64 => BTC

// AmountString convert amount `1.0 BTC` to `1.0` as string
func (b *Bitcoin) AmountString(amt btcutil.Amount) string {
	return strings.TrimRight(amt.String(), " BTC")
}

// FloatToAmount convert float to amount
// e.g. 0.54 to 54000000
func (b *Bitcoin) FloatToAmount(f float64) (btcutil.Amount, error) {
	amt, err := btcutil.NewAmount(f)
	if err != nil {
		return 0, errors.Wrapf(err, "fail to call btcutil.NewAmount(%f)", f)
	}

	return amt, nil
}

// StrToAmount convert string to amount
func (b *Bitcoin) StrToAmount(s string) (btcutil.Amount, error) {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, errors.Wrapf(err, "fail to call strconv.ParseFloat(%s)", s)
	}

	amt, err := btcutil.NewAmount(f)
	if err != nil {
		return 0, errors.Wrapf(err, "fail to call btcutil.NewAmount(%f)", f)
	}

	return amt, nil
}

// StrSatoshiToAmount convert satoshi of string type to amount
func (b *Bitcoin) StrSatoshiToAmount(s string) (btcutil.Amount, error) {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, errors.Wrapf(err, "fail to call strconv.ParseFloat(%s)", s)
	}

	val, err := btcutil.NewAmount(f / float64(100000000))
	if err != nil {
		return 0, errors.Wrapf(err, "fail to call btcutil.NewAmount(%f)", f)
	}

	return val, nil
}
