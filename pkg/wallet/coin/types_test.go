package coin

import (
	"testing"

	"github.com/btcsuite/btcd/chaincfg"
)

func TestCoinType(t *testing.T) {
	type args struct {
		strCoinTypeCode string
		conf            *chaincfg.Params
	}
	tests := []struct {
		name string
		args args
		want CoinType
	}{
		{
			name: "bitcoin, mainnet",
			args: args{"btc", &chaincfg.MainNetParams},
			want: CoinTypeBitcoin,
		},
		{
			name: "bitcoin, testnet",
			args: args{"btc", &chaincfg.TestNet3Params},
			want: CoinTypeTestnet,
		},
		{
			name: "blank, testnet",
			args: args{"", &chaincfg.TestNet3Params},
			want: CoinTypeTestnet,
		},
		{
			name: "bitcoin cash, mainnet",
			args: args{"bch", &chaincfg.MainNetParams},
			want: CoinTypeBitcoinCash,
		},
		{
			name: "bitcoin cash, testnet",
			args: args{"bch", &chaincfg.TestNet3Params},
			want: CoinTypeTestnet,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			coinTypeCode := CoinTypeCode(tt.args.strCoinTypeCode)
			if got := coinTypeCode.CoinType(tt.args.conf); got != tt.want {
				t.Errorf("CoinType() = %d, want %d", got, tt.want)
			}
		})
	}
}
