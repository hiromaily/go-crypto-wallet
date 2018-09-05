package api_test

import (
	"fmt"
	"testing"

	"github.com/hiromaily/go-bitcoin/pkg/enum"
	"github.com/hiromaily/go-bitcoin/pkg/key"
	"github.com/icrowley/fake"
)

// TestSetLabel
func TestSetLabel(t *testing.T) {
	if wlt.BTC.Version() < enum.BTCVer17 {
		t.Skipf("BTC version should be above %d", enum.BTCVer17)
	}

	// アドレスを生成
	bSeed, err := key.GenerateSeed()
	if err != nil {
		t.Fatal(err)
	}
	//keys, err := wlt.GenerateAccountKey(enum.AccountTypeClient, bSeed, 10)
	// key生成
	priv, _, err := key.CreateAccount(wlt.BTC.GetChainConf(), bSeed, enum.AccountTypeClient)
	if err != nil {
		t.Fatal(err)
	}

	walletKeys, err := key.CreateKeysWithIndex(wlt.BTC.GetChainConf(), priv, 0, 10)
	if err != nil {
		t.Fatal(err)
	}

	for _, key := range walletKeys {
		// アドレスをimport
		err = wlt.BTC.ImportAddressWithoutReScan(key.P2shSegwit)
		if err != nil {
			t.Fatal(err)
		}

		// ラベルをセット
		labelName := fake.FirstName()
		fmt.Printf("[label name] %s\n", labelName)
		err = wlt.BTC.SetLabel(key.P2shSegwit, labelName)
		if err != nil {
			t.Fatal(err)
		}

		//セットされているか確認
		info, err := wlt.BTC.GetAddressInfo(key.P2shSegwit)
		if err != nil {
			t.Fatal(err)
		}
		if info.Label != labelName {
			t.Errorf("Label:%s is expected but %s is returned", labelName, info.Label)
		}
	}
}
