package eth_test

import (
	"testing"

	"github.com/hiromaily/go-crypto-wallet/pkg/testutil"
)

// TestGasPrice is test for GasPrice
func TestGasPrice(t *testing.T) {
	// t.SkipNow()
	et := testutil.GetETH()

	price, err := et.GasPrice()
	if err != nil {
		t.Fatal(err)
	}
	t.Log("gasPrice:", price)
}

// TestEstimateGas is test for EstimateGas
// TODO: implement
//func TestEstimateGas(t *testing.T) {
//}
