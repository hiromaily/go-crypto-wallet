package api

import (
	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
)

//TODO:参考に
//https://www.haowuliaoa.com/article/info/11350.html

// GetTransactionByTxID txIDからトランザクション詳細を取得する
func (b *Bitcoin) GetTransactionByTxID(txID string) (*btcjson.GetTransactionResult, error) {
	// Transaction詳細を取得(必要な情報があるかどうか不明)
	hash, err := chainhash.NewHashFromStr(txID)
	if err != nil {
		return nil, err
	}
	return b.Client.GetTransaction(hash)
}

// CreateRawTransaction Rawトランザクションを作成する
func (b *Bitcoin) CreateRawTransaction(sendAddr string, amount btcutil.Amount, inputs []btcjson.TransactionInput) (*wire.MsgTx, error) {
	sendAddrDecoded, err := btcutil.DecodeAddress(sendAddr, b.GetChainConf())
	//TODO:sendAddrの厳密なチェックがセキュリティ的に必要な場面もありそう
	if err != nil {
		return nil, err
	}

	outputs := make(map[btcutil.Address]btcutil.Amount)
	outputs[sendAddrDecoded] = amount //satoshi
	lockTime := int64(0)              //TODO:ここは何をいれるべき？
	return b.Client.CreateRawTransaction(inputs, outputs, &lockTime)
}

func (b *Bitcoin) GetTxOutByTxID(txID string, index uint32) (*btcjson.GetTxOutResult, error) {
	hash, err := chainhash.NewHashFromStr(txID)
	if err != nil {
		return nil, err
	}

	// Gettxout
	// txHash *chainhash.Hash, index uint32, mempool bool
	//return b.Client.GetTxOut(hash, 0, false)
	return b.Client.GetTxOut(hash, index, false)
	//log.Printf("TxOut: %v\n", txOut)
	//grok.Value(txOut)
}
