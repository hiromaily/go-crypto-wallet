package coin_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/hiromaily/go-crypto-wallet/internal/domain/coin"
)

func TestIsCoinTypeCode(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		code string
		want bool
	}{
		{"HYC", "hyc", true},
		{"HYT", "hyt", true},
		{"ETH", "eth", true},
		{"BTC", "btc", true},
		{"unknown", "unknown", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.want, coin.IsCoinTypeCode(tt.code))
		})
	}
}

func TestIsERC20Token(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		code string
		want bool
	}{
		{"HYC", "hyc", true},
		{"HYT", "hyt", true},
		{"ETH_not_erc20", "eth", false},
		{"BTC_not_erc20", "btc", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.want, coin.IsERC20Token(tt.code))
		})
	}
}

func TestIsETHGroup(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		code coin.CoinTypeCode
		want bool
	}{
		{"HYC", coin.HYC, true},
		{"HYT", coin.HYT, true},
		{"ETH", coin.ETH, true},
		{"BTC", coin.BTC, false},
		{"XRP", coin.XRP, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.want, coin.IsETHGroup(tt.code))
		})
	}
}

func TestBIP44AccountPath_HYC(t *testing.T) {
	t.Parallel()
	assert.Equal(t, "m/44'/9002'/0'", coin.BIP44AccountPath(coin.HYC, 0))
}
