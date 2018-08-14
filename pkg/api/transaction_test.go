package api_test

import (
	"testing"
)

// TestGetTransactionByTxID
func TestGetTransactionByTxID(t *testing.T) {
	//実行環境によって、このtxIDが機能しないかも
	txID := "e1fb7fb4c345c81d6852c4cfc51f7a83aaa63587db85317db92bec5839ad3bfb"

	resTx, err := wlt.Btc.GetTransactionByTxID(txID)
	if err != nil {
		t.Errorf("GetTransactionByTxID() error: %v", err)
	}
	t.Log(resTx)
	//type GetTransactionResult struct {
	//	Amount          float64                       `json:"amount"`
	//	Fee             float64                       `json:"fee,omitempty"`
	//	Confirmations   int64                         `json:"confirmations"`
	//	BlockHash       string                        `json:"blockhash"`
	//	BlockIndex      int64                         `json:"blockindex"`
	//	BlockTime       int64                         `json:"blocktime"`
	//	TxID            string                        `json:"txid"`
	//	WalletConflicts []string                      `json:"walletconflicts"`
	//	Time            int64                         `json:"time"`
	//	TimeReceived    int64                         `json:"timereceived"`
	//	Details         []GetTransactionDetailsResult `json:"details"`
	//	Hex             string                        `json:"hex"`
	//}
}
