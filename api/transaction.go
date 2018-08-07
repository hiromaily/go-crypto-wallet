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

// GetTxOutByTxID TxOutを指定したトランザクションIDから取得する
func (b *Bitcoin) GetTxOutByTxID(txID string, index uint32) (*btcjson.GetTxOutResult, error) {
	hash, err := chainhash.NewHashFromStr(txID)
	if err != nil {
		return nil, err
	}

	// Gettxout
	// txHash *chainhash.Hash, index uint32, mempool bool
	//return b.Client.GetTxOut(hash, 0, false)
	return b.Client.GetTxOut(hash, index, false)
	//log.Printf("TxOut: %v\n", txOut): Output
	//grok.Value(txOut)
	//value *GetTxOutResult = {
	//	BestBlock string = "00000000000000a080461b99935872934fa35bc705453f9f2ad7712b2228e849" 64
	//	Confirmations int64 = 473
	//	Value float64 = 0.65
	//	ScriptPubKey ScriptPubKeyResult = {
	//		Asm string = "OP_HASH160 b24f4d8c8237c73a7299b6e790b309d477bb509c OP_EQUAL" 60
	//		Hex string = "a914b24f4d8c8237c73a7299b6e790b309d477bb509c87" 46
	//		ReqSigs int32 = 1
	//		Type string = "scripthash" 10
	//		Addresses []string = [
	//			0 string = "2N9W3GVam33jQc5FbkLKwMqH7RkvkYK7xvz" 35
	//		]
	//	}
	//	Coinbase bool = false
	//}
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

//Signrawtransaction
//TODO: It should be implemented on Cold Strage
//この処理がHotwallet内で動くということは、重要な情報がwallet内に含まれてしまっているということでは？
//signed, isSigned, err := bit.Client.SignRawTransaction(msgTx)
//if err != nil {
//	log.Fatal(err)
//}
//log.Printf("Signrawtransaction: %v\n", signed)
//log.Printf("Signrawtransaction isSigned: %v\n", isSigned)

//Sendrawtransaction

//TODO:トランザクションのkbに応じて、手数料を算出
