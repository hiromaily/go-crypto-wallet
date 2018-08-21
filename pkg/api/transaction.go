package api

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"log"

	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
	"github.com/pkg/errors"
)

//TODO:参考に(中国語のサイト)
//https://www.haowuliaoa.com/article/info/11350.html

//transactionの命名ルールについて
//*wire.MsgTxなどは、txがsuffixとして扱われているため、それに習う。e.g. hexTx, hashTxなど

//入金時における、トランザクション作成順序
//[Online]  CreateRawTransaction
//[Offline] SignRawTransactionByHex

// FundRawTransactionResult fundrawtransactionをcallしたresponseの型
type FundRawTransactionResult struct {
	Hex       string `json:"hex"`
	Fee       int64  `json:"fee"`
	Changepos int64  `json:"changepos"`
}

// ToHex 16進数のstringに変換する
func (b *Bitcoin) ToHex(tx *wire.MsgTx) (string, error) {
	buf := bytes.NewBuffer(make([]byte, 0, tx.SerializeSize()))
	if err := tx.Serialize(buf); err != nil {
		return "", errors.Errorf("tx.Serialize(): error: %v", err)
	}
	return hex.EncodeToString(buf.Bytes()), nil
}

//ToMsgTx 16進数のstringから、wire.MsgTxに変換する
func (b *Bitcoin) ToMsgTx(txHex string) (*wire.MsgTx, error) {
	byteHex, err := hex.DecodeString(txHex)
	if err != nil {
		return nil, errors.Errorf("hex.DecodeString(): error: %v", err)
	}

	var msgTx wire.MsgTx
	if err := msgTx.Deserialize(bytes.NewReader(byteHex)); err != nil {
		return nil, err
	}

	return &msgTx, nil
}

// DecodeRawTransaction Hex stringをデコードして、Rawのトランザクションデータに変換する
func (b *Bitcoin) DecodeRawTransaction(hexTx string) (*btcjson.TxRawResult, error) {
	byteHex, err := hex.DecodeString(hexTx)
	if err != nil {
		return nil, errors.Errorf("hex.DecodeString(): error: %v", err)
	}
	resTx, err := b.client.DecodeRawTransaction(byteHex)
	if err != nil {
		return nil, errors.Errorf("client.DecodeRawTransaction(): error: %v", err)
	}

	return resTx, nil
}

// GetRawTransactionByHex Hexからトランザクションを取得する
func (b *Bitcoin) GetRawTransactionByHex(strHashTx string) (*btcutil.Tx, error) {

	hashTx, err := chainhash.NewHashFromStr(strHashTx)
	if err != nil {
		return nil, errors.Errorf("chainhash.NewHashFromStr(%s): error: %v", strHashTx, err)
	}

	tx, err := b.client.GetRawTransaction(hashTx)
	if err != nil {
		return nil, errors.Errorf("GetRawTransaction(hash): error: %v", err)
	}
	//MsgTx()
	//tx.MsgTx()

	return tx, nil
}

// GetTransactionByTxID txIDからトランザクション詳細を取得する
func (b *Bitcoin) GetTransactionByTxID(txID string) (*btcjson.GetTransactionResult, error) {
	// Transaction詳細を取得(必要な情報があるかどうか不明)
	hashTx, err := chainhash.NewHashFromStr(txID)
	if err != nil {
		return nil, errors.Errorf("chainhash.NewHashFromStr(%s): error: %v", txID, err)
	}
	resTx, err := b.client.GetTransaction(hashTx)
	if err != nil {
		return nil, errors.Errorf("GetTransaction(%s): error: %v", hashTx, err)
	}
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

	return resTx, nil
}

// GetTxOutByTxID TxOutを指定したトランザクションIDから取得する
func (b *Bitcoin) GetTxOutByTxID(txID string, index uint32) (*btcjson.GetTxOutResult, error) {
	hash, err := chainhash.NewHashFromStr(txID)
	if err != nil {
		return nil, errors.Errorf("chainhash.NewHashFromStr(%s): error: %v", txID, err)
	}

	// Gettxout / txHash *chainhash.Hash, index uint32, mempool bool
	txOutResult, err := b.client.GetTxOut(hash, index, false)
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

// CreateRawTransaction Rawトランザクションを作成する
//  => watch only wallet(online)で利用されることを想定
// こちらは多対1用の送金、つまり入金時、集約用アドレスに一括して送信するケースで利用することを想定
// [Noted] 手数料を考慮せず、全額送金しようとすると、SendRawTransaction()で、`min relay fee not met`
func (b *Bitcoin) CreateRawTransaction(sendAddr string, amount btcutil.Amount, inputs []btcjson.TransactionInput) (*wire.MsgTx, error) {
	//TODO:sendAddrの厳密なチェックがセキュリティ的に必要な場面もありそう
	sendAddrDecoded, err := btcutil.DecodeAddress(sendAddr, b.GetChainConf())
	if err != nil {
		return nil, errors.Errorf("btcutil.DecodeAddress(%s): error: %v", sendAddr, err)
	}

	log.Printf("[Debug] amount:%d, %v", amount, amount) // 1.95 BTC => %v表示だと、単位まで表示される！

	// パラメータを作成する
	outputs := make(map[btcutil.Address]btcutil.Amount)
	outputs[sendAddrDecoded] = amount //satoshi
	//lockTime := int64(0) //TODO:Raw locktime ここは何をいれるべき？

	// CreateRawTransaction
	//TODO:ここから下はCreateRawTransactionWithOutput()を呼び出すことでコードの重複を防いだほうがいい。
	return b.CreateRawTransactionWithOutput(inputs, outputs)
}

// CreateRawTransactionWithOutput 出金時に1対多で送信する場合に利用するトランザクションを作成する
func (b *Bitcoin) CreateRawTransactionWithOutput(inputs []btcjson.TransactionInput, outputs map[btcutil.Address]btcutil.Amount) (*wire.MsgTx, error) {
	lockTime := int64(0) //TODO:Raw locktime ここは何をいれるべき？

	// CreateRawTransaction
	msgTx, err := b.client.CreateRawTransaction(inputs, outputs, &lockTime)
	if err != nil {
		return nil, errors.Errorf("btcutil.CreateRawTransaction(): error: %v", err)
	}

	return msgTx, nil
}

// FundRawTransaction 送信したい金額に応じて、自動的にutxoを算出してくれる
//  現時点で使う予定無し
func (b *Bitcoin) FundRawTransaction(hex string) (*FundRawTransactionResult, error) {
	//fundrawtransaction
	//https://bitcoincore.org/en/doc/0.16.2/rpc/rawtransactions/fundrawtransaction/
	//"{\"changePosition\":2}"

	//hex
	bHex, err := json.Marshal(hex)
	if err != nil {
		return nil, errors.Errorf("json.Marchal(hex): error: %v", err)
	}

	//fee rate
	feePerKb, err := b.EstimateSmartFee()
	if err != nil {
		return nil, errors.Errorf("EstimateSmartFee(): error: %v", err)
	}

	bFeeRate, err := json.Marshal(struct {
		FeeRate float64 `json:"feeRate"`
	}{
		FeeRate: feePerKb,
	})
	if err != nil {
		return nil, errors.Errorf("json.Marchal(feeRate): error: %v", err)
	}

	rawResult, err := b.client.RawRequest("fundrawtransaction", []json.RawMessage{bHex, bFeeRate})
	//rawResult, err := b.client.RawRequest("fundrawtransaction", []json.RawMessage{bHex})
	if err != nil {
		//error: -4: Insufficient funds
		return nil, errors.Errorf("json.RawRequest(fundrawtransaction): error: %v", err)
	}

	fundRawTransactionResult := FundRawTransactionResult{}
	err = json.Unmarshal([]byte(rawResult), &fundRawTransactionResult)
	if err != nil {
		return nil, errors.Errorf("json.Unmarshal(): error: %v", err)
	}

	log.Printf("[Debug]fundRawTransactionResult: %v: %s\n", fundRawTransactionResult, fundRawTransactionResult.Hex)

	return &fundRawTransactionResult, nil
}

// SignRawTransaction *wire.MsgTxからRawのトランザクションに署名する
// こちらは一連の流れを調査するために、CreateRawTransaction()の戻りに合わせたI/F (実際の運用で利用されることはないはず)
// 秘密鍵を保持している側のwallet(つまりcold wallet)で実行することを想定
func (b *Bitcoin) SignRawTransaction(tx *wire.MsgTx) (*wire.MsgTx, bool, error) {
	//署名
	msgTx, isSigned, err := b.client.SignRawTransaction(tx)
	if err != nil {
		return nil, false, errors.Errorf("SignRawTransaction(): error: %v", err)
	}
	//Multisigの場合、これによって署名が終了したか判断するはず
	//if !isSigned {
	//	return nil, errors.New("SignRawTransaction() can not sign on given transaction")
	//}

	return msgTx, isSigned, nil
}

// SignRawTransactionByHex HexからRawトランザクションを生成し、署名する
// 秘密鍵を保持している側のwallet(つまりcold wallet)で実行することを想定
func (b *Bitcoin) SignRawTransactionByHex(hex string) (string, bool, error) {
	// Hexからトランザクションを取得
	msgTx, err := b.ToMsgTx(hex)
	if err != nil {
		return "", false, err
	}

	//署名
	signedTx, isSigned, err := b.SignRawTransaction(msgTx)
	if err != nil {
		return "", false, err
	}

	//Hexに変換
	hexTx, err := b.ToHex(signedTx)
	if err != nil {
		return "", false, errors.Errorf("w.BTC.ToHex(msgTx): error: %v", err)
	}

	//return signedTx, nil
	return hexTx, isSigned, nil
}

// SendRawTransaction Rawトランザクションを送信する
// こちらは一連の流れを調査するために、CreateRawTransaction()の戻りに合わせたI/F (実際の運用で利用されることはないはず)
// オンラインで実行される必要がある
func (b *Bitcoin) SendRawTransaction(tx *wire.MsgTx) (*chainhash.Hash, error) {
	//送信
	hash, err := b.client.SendRawTransaction(tx, true)
	if err != nil {
		return nil, errors.Errorf("client.SendRawTransaction(): error: %v", err)
	}

	return hash, nil
}

// SendTransactionByHex 外部から渡されたバイト列からRawトランザクションを送信する
// オンラインで実行される必要がある
func (b *Bitcoin) SendTransactionByHex(hex string) (*chainhash.Hash, error) {
	// Hexからトランザクションを取得
	msgTx, err := b.ToMsgTx(hex)
	if err != nil {
		return nil, err
	}

	//送信
	hash, err := b.SendRawTransaction(msgTx)
	//hash, err := b.client.SendRawTransaction(msgTx, true)
	if err != nil {
		return nil, errors.Errorf("SendRawTransaction(): error: %v", err)
	}

	//txID
	//hash.String()

	return hash, nil
}

// SendTransactionByByte 外部から渡されたバイト列からRawトランザクションを送信する
// オンラインで実行される必要がある
func (b *Bitcoin) SendTransactionByByte(rawTx []byte) (*chainhash.Hash, error) {

	//[]byte to wireTx
	wireTx := new(wire.MsgTx)
	r := bytes.NewBuffer(rawTx)

	if err := wireTx.Deserialize(r); err != nil {
		return nil, errors.Errorf("wireTx.Deserialize(): error: %v", err)
	}

	//送信
	hash, err := b.SendRawTransaction(wireTx)
	//hash, err := b.client.SendRawTransaction(wireTx, true)
	if err != nil {
		return nil, errors.Errorf("SendRawTransaction(): error: %v", err)
	}

	return hash, nil
}

// SequentialTransaction [Debug用]: 一連の未署名トランザクション作成から送信までの流れ
// TODO:Testに移動するか
func (b *Bitcoin) SequentialTransaction(hex string) (*chainhash.Hash, *btcutil.Tx, error) {
	// Hexからトランザクションを取得
	msgTx, err := b.ToMsgTx(hex)
	if err != nil {
		return nil, nil, err
	}

	//署名(オフライン)
	signedTx, isSigned, err := b.SignRawTransaction(msgTx)
	if err != nil {
		return nil, nil, err
	}
	if !isSigned {
		return nil, nil, errors.New("SignRawTransaction() can not sign on given transaction or multisig may be required")
	}

	//送金(オンライン)
	hash, err := b.SendRawTransaction(signedTx)
	if err != nil {
		return nil, nil, err
	}
	log.Printf("[Debug] txID hash: %s", hash.String())

	//txを取得
	resTx, err := b.GetRawTransactionByHex(hash.String())
	if err != nil {
		return nil, nil, err
	}

	return hash, resTx, nil
}

// Sign 署名を行う
//FIXME: これはColdWallet内で必要となるが、BitcoinCoreの機能が必要ないので、あれば実装しておきたい
func (b *Bitcoin) Sign(tx *wire.MsgTx, strPrivateKey string) (string, error) {
	// Key
	wif, err := btcutil.DecodeWIF(strPrivateKey)
	if err != nil {
		return "", err
	}
	privKey := wif.PrivKey

	// SignatureScript
	for idx, val := range tx.TxIn {
		//
		script, err := txscript.SignatureScript(tx, idx, val.SignatureScript, txscript.SigHashAll, privKey, false)
		if err != nil {
			return "", err
		}
		tx.TxIn[idx].SignatureScript = script
	}
	//TODO: isSignedかどうかをどうチェックするか
	//TODO: TxInごとに異なるKeyの場合は難しい

	//Hexに変換
	hexTx, err := b.ToHex(tx)
	if err != nil {
		return "", err
	}

	return hexTx, nil
}
