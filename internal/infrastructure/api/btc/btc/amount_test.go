package btc_test

import (
	"testing"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	btcpkg "github.com/hiromaily/go-crypto-wallet/pkg/chains/btc"
)

// TestAmountString is test for AmountString
func TestAmountString(t *testing.T) {
	tests := []struct {
		amt  btcutil.Amount
		want string
	}{
		{54000000, "0.54"},
		{480000, "0.0048"},
		{195200000, "1.952"},
		{1000, "0.00001"},
	}
	for _, val := range tests {
		res := btcpkg.AmountString(val.amt)
		assert.Equal(t, val.want, res, "AmountString() result mismatch")
	}
}

// TestAmountToDecimal is test for AmountToDecimal
func TestAmountToDecimal(t *testing.T) {
	tests := []struct {
		amt  btcutil.Amount
		want string
	}{
		{54000000, "0.54"},
		{480000, "0.0048"},
		{195200000, "1.952"},
		{1000, "0.00001"},
	}
	for _, val := range tests {
		res, err := btcpkg.AmountToDecimal(val.amt)
		require.NoError(t, err, "AmountToDecimal() should not return error")
		assert.Equal(t, val.want, res.String(), "AmountString() result mismatch")
	}
}

// TestFloatToAmount is test for FloatToAmount
func TestFloatToAmount(t *testing.T) {
	tests := []struct {
		bit  float64
		want btcutil.Amount
	}{
		{0.54, 54000000},
		{0.0048, 480000},
		{1.9520, 195200000},
		{0.000010, 1000},
	}
	for _, val := range tests {
		amt, err := btcpkg.FloatToAmount(val.bit)
		require.NoError(t, err, "FloatToAmount() should not return error")
		assert.Equal(t, val.want, amt, "FloatToAmount() result mismatch")

		t.Logf("satoshi: %d, %v", amt, amt)
		// satoshi: 54000000, 0.54 BTC
		// satoshi: 195200000, 1.952 BTC
	}
}

// TestStrToAmount is test for StrToAmount
func TestStrToAmount(t *testing.T) {
	tests := []struct {
		bit  string
		want btcutil.Amount
	}{
		{"10", 1000000000},
		{"8.04", 804000000},
		{"0.5", 50000000},
		{"0.0025", 250000},
		{"0.0000040000", 400},
		{"0.000010", 1000},
	}
	for _, val := range tests {
		amt, err := btcpkg.StrToAmount(val.bit)
		require.NoError(t, err, "StrToAmount() should not return error")
		assert.Equal(t, val.want, amt, "StrToAmount() result mismatch")

		t.Logf("satoshi: %d, %v", amt, amt)
		// satoshi: 1000000000, 10 BTC
		// satoshi: 50000000, 0.5 BTC
	}
}

// TestStrSatoshiToAmount is test for StrSatoshiToAmount
func TestStrSatoshiToAmount(t *testing.T) {
	tests := []struct {
		satoshi string
		want    float64
	}{
		{"1000000000", 10},
		{"804000000", 8.04},
		{"50000000", 0.5},
		{"250000", 0.0025},
	}
	for _, val := range tests {
		amt, err := btcpkg.StrSatoshiToAmount(val.satoshi)
		require.NoError(t, err, "StrSatoshiToAmount() should not return error")
		assert.Equal(t, val.want, amt.ToBTC(), "StrToAmount() result mismatch")

		t.Logf("satoshi: %d, %v", amt, amt)
		// satoshi: 1000000000, 10 BTC
		// satoshi: 50000000, 0.5 BTC
	}
}

// TestCalculation is test for calculation of amount
func TestCalculation(t *testing.T) {
	tests := []struct {
		val1 float64
		val2 float64
		want float64
	}{
		{0.156, 0.3, 0.456},
		{2.567, 0.111, 2.678},
	}
	for _, v := range tests {
		amt1, _ := btcutil.NewAmount(v.val1)
		amt2, _ := btcutil.NewAmount(v.val2)
		res := (amt1 + amt2).ToBTC()
		assert.Equal(t, v.want, res, "Calculation result mismatch")
		t.Logf("%f + %f = %f", v.val1, v.val2, res)
	}
}
