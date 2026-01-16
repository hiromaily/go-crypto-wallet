package btc

import (
	"testing"

	"github.com/btcsuite/btcd/btcutil/hdkeychain"
	"github.com/btcsuite/btcd/chaincfg"
)

const benchmarkXpub = "xpub6ERApfZwUNrhLCkDtcHTcxd75RbzS1ed54G1LkBUHQVHQKqhMkhgbmJbZRkrgZw4koxb5JaHWkY4ALHY" +
	"2grBGRjaDMzQLcgJvLJuZZvRcEL"

func BenchmarkGenerateDescriptors(b *testing.B) {
	service := NewDescriptorService(&chaincfg.MainNetParams)
	xpub, err := hdkeychain.NewKeyFromString(benchmarkXpub)
	if err != nil {
		b.Fatalf("parse xpub: %v", err)
	}

	for i := 0; i < b.N; i++ {
		if _, genErr := service.GenerateTaprootDescriptor("a1b2c3d4", "/86'/0'/0'", xpub, false); genErr != nil {
			b.Fatalf("generate taproot descriptor: %v", genErr)
		}
		if _, genErr := service.GenerateBech32Descriptor("a1b2c3d4", "/84'/0'/0'", xpub, true); genErr != nil {
			b.Fatalf("generate bech32 descriptor: %v", genErr)
		}
	}
}

func BenchmarkParseDescriptor(b *testing.B) {
	parser := NewDescriptorParser()
	desc := "wpkh([a1b2c3d4/84'/0'/0']" + benchmarkXpub + "/0/*)"

	for i := 0; i < b.N; i++ {
		if _, err := parser.Parse(desc); err != nil {
			b.Fatalf("parse descriptor: %v", err)
		}
	}
}
