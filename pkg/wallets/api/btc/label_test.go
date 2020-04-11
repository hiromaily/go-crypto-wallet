package btc_test

import (
	"fmt"
	"testing"

	"github.com/icrowley/fake"

	"github.com/hiromaily/go-bitcoin/pkg/account"
	"github.com/hiromaily/go-bitcoin/pkg/enum"
	"github.com/hiromaily/go-bitcoin/pkg/wallets/key"
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
	keyData := key.NewKey(enum.BTC, wlt.BTC.GetChainConf())

	priv, _, err := keyData.CreateAccount(bSeed, account.AccountTypeClient)
	if err != nil {
		t.Fatal(err)
	}

	walletKeys, err := keyData.CreateKeysWithIndex(priv, 0, 10)
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
