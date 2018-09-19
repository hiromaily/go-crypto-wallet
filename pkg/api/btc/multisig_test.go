package btc_test

import (
	"github.com/hiromaily/go-bitcoin/pkg/enum"
	"testing"
)

// TestCreateMultiSig
func TestCreateMultiSig(t *testing.T) {
	//getnewaddress taro 2N7ZwUXpo841GZDpxLGFqrhr1xwMzTba7ZP
	//getnewaddress boss1 2NAm558FWpiaJQLz838vbzBPpqmKxyeyxsu
	//TODO:ここで、AddMultisigAddressを使うのにパラメータとしてaccout名も渡さないといけない。。これをどうすべきか。。。
	//TODO: => おそらくBlankでもいい

	//TODO: Multisigアドレス作成
	resAddr, err := wlt.BTC.AddMultisigAddress(2, []string{"2N7ZwUXpo841GZDpxLGFqrhr1xwMzTba7ZP", "2NAm558FWpiaJQLz838vbzBPpqmKxyeyxsu"}, "multi01", enum.AddressTypeLegacy)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("multisig address: %s, redeemScript: %s", resAddr.Address, resAddr.RedeemScript)
	//multisig address: 2N4Rm1aLPxCcg1H1V96bBzH69vMAipADLCQ, redeemScript: 522103d69e07dbf6da065e6fae1ef5761d029b9ff9143e75d579ffc439d47484044bed2103748797877523b8b36add26c9e0fb6a023f05083dd4056aedc658d2932df1eb6052ae

	resAddr, err = wlt.BTC.AddMultisigAddress(2, []string{"2N7ZwUXpo841GZDpxLGFqrhr1xwMzTba7ZP", "2NAm558FWpiaJQLz838vbzBPpqmKxyeyxsu"}, "multi02", enum.AddressTypeP2shSegwit)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("multisig address: %s, redeemScript: %s", resAddr.Address, resAddr.RedeemScript)

}
