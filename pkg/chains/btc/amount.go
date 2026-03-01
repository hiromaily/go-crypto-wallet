package btc

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/quagmt/udecimal"
)

// 0.00000001 BTC=1 Satoshi
// btcutil.Amount => Satoshi
// float64 => BTC

// AmountString converts amount `1.0 BTC` to `1.0` as string.
func AmountString(amt btcutil.Amount) string {
	s := strings.TrimRight(amt.String(), " BTC")
	// Remove trailing zeros after decimal point
	if strings.Contains(s, ".") {
		s = strings.TrimRight(s, "0")
		s = strings.TrimRight(s, ".")
	}
	return s
}

// AmountToDecimal converts amount `1.0 BTC` to `1.0` as decimal.
func AmountToDecimal(amt btcutil.Amount) (udecimal.Decimal, error) {
	strAmt := strings.TrimRight(amt.String(), " BTC")
	// Remove trailing zeros after decimal point
	if strings.Contains(strAmt, ".") {
		strAmt = strings.TrimRight(strAmt, "0")
		strAmt = strings.TrimRight(strAmt, ".")
	}
	return udecimal.Parse(strAmt)
}

// FloatToDecimal converts float to decimal.
func FloatToDecimal(f float64) (udecimal.Decimal, error) {
	return udecimal.NewFromFloat64(f)
}

// FloatToAmount converts float to amount.
// e.g. 0.54 to 54000000
func FloatToAmount(f float64) (btcutil.Amount, error) {
	amt, err := btcutil.NewAmount(f)
	if err != nil {
		return 0, fmt.Errorf("fail to call btcutil.NewAmount(%f): %w", f, err)
	}
	return amt, nil
}

// StrToAmount converts string to amount.
func StrToAmount(s string) (btcutil.Amount, error) {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, fmt.Errorf("fail to call strconv.ParseFloat(%s): %w", s, err)
	}
	amt, err := btcutil.NewAmount(f)
	if err != nil {
		return 0, fmt.Errorf("fail to call btcutil.NewAmount(%f): %w", f, err)
	}
	return amt, nil
}

// StrSatoshiToAmount converts satoshi of string type to amount.
func StrSatoshiToAmount(s string) (btcutil.Amount, error) {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, fmt.Errorf("fail to call strconv.ParseFloat(%s): %w", s, err)
	}
	val, err := btcutil.NewAmount(f / float64(100000000))
	if err != nil {
		return 0, fmt.Errorf("fail to call btcutil.NewAmount(%f): %w", f, err)
	}
	return val, nil
}

// CalculateNewFee adjusts a fee by multiplying it by adjustmentFee.
// Used by both BTC and BCH fee calculation logic.
func CalculateNewFee(fee btcutil.Amount, adjustmentFee float64) (btcutil.Amount, error) {
	return FloatToAmount(fee.ToBTC() * adjustmentFee)
}
