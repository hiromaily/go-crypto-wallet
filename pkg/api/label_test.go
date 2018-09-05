package api_test

import (
	"fmt"
	"testing"

	"github.com/hiromaily/go-bitcoin/pkg/enum"
	"github.com/icrowley/fake"
)

// TestSetLabel
func TestSetLabel(t *testing.T) {
	if wlt.BTC.Version() < enum.BTCVer17 {
		t.Skipf("BTC version should be above %d", enum.BTCVer17)
	}

	// アドレスを生成
	//_, err := wlt.BTC.ValidateAddress(val.addr)
	bSeed, err := wlt.GenerateSeed()
	if err != nil {
		t.Fatal(err)
	}
	keys, err := wlt.GenerateAccountKey(enum.AccountTypeClient, bSeed, 10)
	if err != nil {
		t.Fatal(err)
	}

	for _, key := range keys {
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

		//
	}

	//var tests = []struct {
	//	addr  string
	//	isErr bool
	//}{
	//	{"2NFXSXxw8Fa6P6CSovkdjXE6UF4hupcTHtr", false},
	//	{"2NDGkbQTwg2v1zP6yHZw3UJhmsBh9igsSos", false},
	//	{"4VHGkbQTGg2vN5P6yHZw3UJhmsBh9igsSos", true},
	//}
	//
	//for _, val := range tests {
	//	//t.Logf("check address: %s", val.addr)
	//	fmt.Printf("check address: %s\n", val.addr)
	//
	//	_, err := wlt.BTC.ValidateAddress(val.addr)
	//	if err != nil && !val.isErr {
	//		t.Errorf("Unexpectedly error occorred. %v", err)
	//	}
	//	if err == nil && val.isErr {
	//		t.Error("Error is expected. However nothing happened.")
	//	}
	//}
}
