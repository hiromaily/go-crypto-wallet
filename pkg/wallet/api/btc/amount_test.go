package btc_test

import (
	"testing"

	"github.com/btcsuite/btcutil"

	. "github.com/hiromaily/go-bitcoin/pkg/wallet/api/btc"
)

// TestFloatToAmount
func TestFloatToAmount(t *testing.T) {
	var btc = Bitcoin{}

	var tests = []struct {
		bit  float64
		want btcutil.Amount
	}{
		{0.54, 54000000},
		{0.0048, 480000},
		{1.9520, 195200000},
		{0.000010, 1000},
	}
	for _, val := range tests {
		amt, err := btc.FloatToAmount(val.bit)
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

// TestStrSatoshiToAmount
func TestStrToAmount(t *testing.T) {
	var btc = Bitcoin{}

	var tests = []struct {
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
		amt, err := btc.StrToAmount(val.bit)
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

// TestStrSatoshiToAmount
func TestStrSatoshiToAmount(t *testing.T) {
	var btc = Bitcoin{}

	var tests = []struct {
		satoshi string
		want    float64
	}{
		{"1000000000", 10},
		{"804000000", 8.04},
		{"50000000", 0.5},
		{"250000", 0.0025},
	}
	for _, val := range tests {
		amt, err := btc.StrSatoshiToAmount(val.satoshi)
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

func TestCaluculation(t *testing.T) {
	var tests = []struct {
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
		if res != v.want {
			t.Errorf("result is %f however %f is expected.", res, v.want)
		}
		t.Logf("%f + %f = %f", v.val1, v.val2, res)
	}

}
