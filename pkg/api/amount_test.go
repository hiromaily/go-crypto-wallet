package api_test

import (
	"testing"

	"github.com/btcsuite/btcutil"
	. "github.com/hiromaily/go-bitcoin/pkg/api"
)

var b = Bitcoin{}

// TestFloatBitToAmount
func TestFloatBitToAmount(t *testing.T) {
	var tests = []struct {
		bit  float64
		want btcutil.Amount
	}{
		{0.54, 54000000},
		{1.9520, 195200000},
	}
	for _, val := range tests {
		amt, err := b.FloatBitToAmount(val.bit)
		if err != nil {
			t.Fatal(err)
		}
		if amt != val.want {
			t.Errorf("result is %d however %d is expected.", amt, val.want)
		}
		//amt.ToBTC() //float64

		t.Logf("satoshi: %d, %v", amt, amt)
		//satoshi: 54000000, 0.54 BTC
		//satoshi: 195200000, 1.952 BTC
	}
}

// TestCastStrSatoshiToAmount
func TestCastStrBitToAmount(t *testing.T) {
	var tests = []struct {
		bit  string
		want btcutil.Amount
	}{
		{"10", 1000000000},
		{"0.5", 50000000},
	}
	for _, val := range tests {
		amt, err := b.CastStrBitToAmount(val.bit)
		if err != nil {
			t.Fatal(err)
		}
		if amt != val.want {
			t.Errorf("result is %d however %d is expected.", amt, val.want)
		}

		t.Logf("satoshi: %d, %v", amt, amt)
		//satoshi: 1000000000, 10 BTC
		//satoshi: 50000000, 0.5 BTC
	}
}

// TestCastStrSatoshiToAmount
func TestCastStrSatoshiToAmount(t *testing.T) {
	var tests = []struct {
		satoshi string
		want    float64
	}{
		{"1000000000", 10},
		{"50000000", 0.5},
	}
	for _, val := range tests {
		amt, err := b.CastStrSatoshiToAmount(val.satoshi)
		if err != nil {
			t.Fatal(err)
		}
		if amt.ToBTC() != val.want {
			t.Errorf("result is %d however %f is expected.", amt, val.want)
		}

		t.Logf("satoshi: %d, %v", amt, amt)
		//satoshi: 1000000000, 10 BTC
		//satoshi: 50000000, 0.5 BTC
	}
}
