package btc

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/ericlagergren/decimal"
	"github.com/pkg/errors"
)

// 0.00000001 BTC=1 Satoshi
// btcutil.Amount => Satoshi
// float64 => BTC

// AmountString converts amount `1.0 BTC` to `1.0` as string
func (*Bitcoin) AmountString(amt btcutil.Amount) string {
	s := strings.TrimRight(amt.String(), " BTC")
	// Remove trailing zeros after decimal point
	if strings.Contains(s, ".") {
		s = strings.TrimRight(s, "0")
		s = strings.TrimRight(s, ".")
	}
	return s
}

// AmountToDecimal converts amount `1.0 BTC` to `1.0` as decimal
func (*Bitcoin) AmountToDecimal(amt btcutil.Amount) *decimal.Big {
	strAmt := strings.TrimRight(amt.String(), " BTC")
	// Remove trailing zeros after decimal point
	if strings.Contains(strAmt, ".") {
		strAmt = strings.TrimRight(strAmt, "0")
		strAmt = strings.TrimRight(strAmt, ".")
	}
	dAmt := new(decimal.Big)
	dAmt, _ = dAmt.SetString(strAmt)
	return dAmt
}

// FloatToDecimal converts float to decimal
func (*Bitcoin) FloatToDecimal(f float64) *decimal.Big {
	strAmt := fmt.Sprintf("%f", f)
	dAmt := new(decimal.Big)
	dAmt, _ = dAmt.SetString(strAmt)
	return dAmt
}

// FloatToAmount converts float to amount
// e.g. 0.54 to 54000000
func (*Bitcoin) FloatToAmount(f float64) (btcutil.Amount, error) {
	amt, err := btcutil.NewAmount(f)
	if err != nil {
		return 0, errors.Wrapf(err, "fail to call btcutil.NewAmount(%f)", f)
	}

	return amt, nil
}

// StrToAmount converts string to amount
func (*Bitcoin) StrToAmount(s string) (btcutil.Amount, error) {
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
func (*Bitcoin) StrSatoshiToAmount(s string) (btcutil.Amount, error) {
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
