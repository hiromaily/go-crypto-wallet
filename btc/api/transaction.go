package api

import (
	"bytes"
	"encoding/hex"

	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
	"github.com/pkg/errors"
)

//TODO:参考に
//https://www.haowuliaoa.com/article/info/11350.html

// GetTransactionByTxID txIDからトランザクション詳細を取得する
func (b *Bitcoin) GetTransactionByTxID(txID string) (*btcjson.GetTransactionResult, error) {
	// Transaction詳細を取得(必要な情報があるかどうか不明)
	hash, err := chainhash.NewHashFromStr(txID)
	if err != nil {
		return nil, errors.Errorf("chainhash.NewHashFromStr(%s): error: %v", txID, err)
	}
	txResult, err := b.Client.GetTransaction(hash)
	if err != nil {
		return nil, errors.Errorf("GetTransaction(%s): error: %v", hash, err)
	}

	return txResult, nil
}

// GetTxOutByTxID TxOutを指定したトランザクションIDから取得する
func (b *Bitcoin) GetTxOutByTxID(txID string, index uint32) (*btcjson.GetTxOutResult, error) {
	hash, err := chainhash.NewHashFromStr(txID)
	if err != nil {
		return nil, errors.Errorf("chainhash.NewHashFromStr(%s): error: %v", txID, err)
	}

	// Gettxout
	// txHash *chainhash.Hash, index uint32, mempool bool
	//return b.Client.GetTxOut(hash, 0, false)
	txOutResult, err := b.Client.GetTxOut(hash, index, false)
	if err != nil {
		return nil, errors.Errorf("GetTxOut(%s, %d, false): error: %v", hash, index, err)
	}

	return txOutResult, nil
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

func (b *Bitcoin) ToHex(tx *wire.MsgTx) (string, error) {
	buf := bytes.NewBuffer(make([]byte, 0, tx.SerializeSize()))
	if err := tx.Serialize(buf); err != nil {
		return "", err
	}
	return hex.EncodeToString(buf.Bytes()), nil
}

// CreateRawTransaction Rawトランザクションを作成する
func (b *Bitcoin) CreateRawTransaction(sendAddr string, amount btcutil.Amount, inputs []btcjson.TransactionInput) (*wire.MsgTx, error) {
	sendAddrDecoded, err := btcutil.DecodeAddress(sendAddr, b.GetChainConf())
	//TODO:sendAddrの厳密なチェックがセキュリティ的に必要な場面もありそう
	//TODO:このタイミングでfeeの計算も必要っぽい
	//TODO:トランザクションのkbに応じて、手数料を算出
	//TODO:でもfeeのパラメータを入れるのは、sendrawtransaction
	if err != nil {
		return nil, errors.Errorf("btcutil.DecodeAddress(%s): error: %v", sendAddr, err)
	}

	outputs := make(map[btcutil.Address]btcutil.Amount)
	outputs[sendAddrDecoded] = amount //satoshi
	lockTime := int64(0)              //TODO:ここは何をいれるべき？
	msgTx, err := b.Client.CreateRawTransaction(inputs, outputs, &lockTime)
	if err != nil {
		return nil, errors.Errorf("btcutil.CreateRawTransaction(): error: %v", err)
	}
	return msgTx, nil
}

// SignRawTransaction
func (b *Bitcoin) SignRawTransaction(tx *wire.MsgTx) (*wire.MsgTx, error) {
	//TODO: It should be implemented on Cold Strage
	//この処理がHotwallet内で動くということは、重要な情報がwallet内に含まれてしまっているということでは？
	msgTx, isSigned, err := b.Client.SignRawTransaction(tx)
	if err != nil {
		return nil, errors.Errorf("SignRawTransaction(): error: %v", err)
	}
	if !isSigned {
		return nil, errors.New("SignRawTransaction() can not sign on given transaction")
	}

	return msgTx, nil
}

func (b *Bitcoin) SendRawTransaction(tx *wire.MsgTx) (*chainhash.Hash, error) {
	hash, err := b.Client.SendRawTransaction(tx, true)
	if err != nil {
		return nil, errors.Errorf("SendRawTransaction(): error: %v", err)
	}

	return hash, nil
}

//func (c *Connector) SendTransaction(rawTx []byte) error {
//	m := crypto.NewMetric(c.cfg.DaemonCfg.Name, string(c.cfg.Asset),
//		MethodSendTransaction, c.cfg.Metrics)
//	defer m.Finish()
//
//	wireTx := new(wire.MsgTx)
//	r := bytes.NewBuffer(rawTx)
//
//	if err := wireTx.Deserialize(r); err != nil {
//		m.AddError(errToSeverity(ErrDeserialiseTx))
//		return errors.Errorf("unable to deserialize raw tx: %v", err)
//	}
//
//	_, err := c.client.SendRawTransaction(wireTx, true)
//	if err != nil {
//		m.AddError(errToSeverity(ErrSendTx))
//		return errors.Errorf("unable to send transaction: %v", err)
//	}
//
//	return nil
//}
