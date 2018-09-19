package btc_test

import (
	"testing"

	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcutil"
)

// TestGetTransactionByTxID
func TestGetTransactionByTxID(t *testing.T) {
	//実行環境によって、このtxIDが機能しないかも
	txID := "e1fb7fb4c345c81d6852c4cfc51f7a83aaa63587db85317db92bec5839ad3bfb"

	resTx, err := wlt.BTC.GetTransactionByTxID(txID)
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

// SequentialTransaction 一連の未署名トランザクション作成から送信までの流れを確認する
// FIXME: パラメータが変更されたため、対応する必要がある
func sequentialTransaction(t *testing.T, hex string) (*chainhash.Hash, *btcutil.Tx, error) {
	// Hexからトランザクションを取得
	//msgTx, err := wlt.BTC.ToMsgTx(hex)
	//if err != nil {
	//	return nil, nil, err
	//}
	//暫定でreturn
	return nil, nil, nil

	// 署名(オフライン)
	//signedTx, isSigned, err := wlt.BTC.SignRawTransaction(msgTx)
	//if err != nil {
	//	return nil, nil, err
	//}
	//if !isSigned {
	//	return nil, nil, errors.New("BTC.SignRawTransaction() can not sign on given transaction or multisig may be required")
	//}
	//
	////送金(オンライン)
	//hash, err := wlt.BTC.SendRawTransaction(signedTx)
	//if err != nil {
	//	return nil, nil, err
	//}
	//
	////txを取得
	//resTx, err := wlt.BTC.GetRawTransactionByHex(hash.String())
	//if err != nil {
	//	return nil, nil, err
	//}
	//
	//return hash, resTx, nil
}
