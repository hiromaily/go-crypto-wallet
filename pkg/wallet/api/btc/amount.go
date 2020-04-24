package btc

import (
	"github.com/ericlagergren/decimal"
	"github.com/volatiletech/sqlboiler/types"
	"strconv"
	"strings"

	"github.com/btcsuite/btcutil"
	"github.com/pkg/errors"
)

//0.00000001 BTC=1 Satoshi
//btcutil.Amount => Satoshi
//float64 => BTC

// AmountString converts amount `1.0 BTC` to `1.0` as string
func (b *Bitcoin) AmountString(amt btcutil.Amount) string {
	return strings.TrimRight(amt.String(), " BTC")
}

// AmountToDecimal converts amount `1.0 BTC` to `1.0` as decimal
func (b *Bitcoin) AmountToDecimal(amt btcutil.Amount) types.Decimal {
	strAmt := strings.TrimRight(amt.String(), " BTC")
	dAmt := types.Decimal{Big: new(decimal.Big)}
	dAmt.Big, _ = dAmt.SetString(strAmt)
	return dAmt
}

// FloatToAmount converts float to amount
// e.g. 0.54 to 54000000
func (b *Bitcoin) FloatToAmount(f float64) (btcutil.Amount, error) {
	amt, err := btcutil.NewAmount(f)
	if err != nil {
		return 0, errors.Wrapf(err, "fail to call btcutil.NewAmount(%f)", f)
	}

	return amt, nil
}

// StrToAmount converts string to amount
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

// StrSatoshiToAmount converts satoshi of string type to amount
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
