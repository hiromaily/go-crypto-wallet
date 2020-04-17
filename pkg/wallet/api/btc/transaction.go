package btc

import (
	"bytes"
	"encoding/hex"
	"encoding/json"

	"github.com/bookerzzz/grok"
	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/hiromaily/go-bitcoin/pkg/account"
)

// refer to https://www.haowuliaoa.com/article/info/11350.html (Chinese site)

// naming regulation for transaction
// as example, *wire.MsgTx is that tx is used as suffix
// so tx should be named like hexTx, hashTx

// SignRawTransactionResult response of api `signrawtransactionwithwallet`
type SignRawTransactionResult struct {
	Hex      string                    `json:"hex"`
	Complete bool                      `json:"complete"`
	Errors   []SignRawTransactionError `json:"errors"`
}

// SignRawTransactionError error object in SignRawTransactionResult
type SignRawTransactionError struct {
	Txid      string `json:"txid"`
	Vout      int64  `json:"vout"`
	ScriptSig string `json:"scriptSig"`
	Sequence  int64  `json:"sequence"`
	Error     string `json:"error"`
}

// AddrsPrevTxs is used when creating tx for multisig address
type AddrsPrevTxs struct {
	Addrs         []string
	PrevTxs       []PrevTx
	SenderAccount account.AccountType
}

// PrevTx is required parameters for api `signrawtransaction` for multisig address
type PrevTx struct {
	Txid         string  `json:"txid"`
	Vout         uint32  `json:"vout"`
	ScriptPubKey string  `json:"scriptPubKey"`
	RedeemScript string  `json:"redeemScript"`
	Amount       float64 `json:"amount"`
}

// FundRawTransactionResult response of api `fundrawtransaction`
type FundRawTransactionResult struct {
	Hex       string `json:"hex"`
	Fee       int64  `json:"fee"`
	Changepos int64  `json:"changepos"`
}

// ToHex convert wire.MsgTx to string of Hexadecimal(16進数)
func (b *Bitcoin) ToHex(tx *wire.MsgTx) (string, error) {
	buf := bytes.NewBuffer(make([]byte, 0, tx.SerializeSize()))
	if err := tx.Serialize(buf); err != nil {
		return "", errors.Wrap(err, "fail to call tx.Serialize()")
	}
	return hex.EncodeToString(buf.Bytes()), nil
}

// ToMsgTx convert string of Hexadecimal(16進数) to wire.MsgTx
func (b *Bitcoin) ToMsgTx(txHex string) (*wire.MsgTx, error) {
	byteHex, err := hex.DecodeString(txHex)
	if err != nil {
		return nil, errors.Wrap(err, "fail to call hex.DecodeString()")
	}

	var msgTx wire.MsgTx
	if err := msgTx.Deserialize(bytes.NewReader(byteHex)); err != nil {
		return nil, errors.Wrap(err, "fail to call hex.Deserialize()")
	}

	return &msgTx, nil
}

// DecodeRawTransaction returns information about a transaction given its serialized byte
func (b *Bitcoin) DecodeRawTransaction(hexTx string) (*btcjson.TxRawResult, error) {
	byteHex, err := hex.DecodeString(hexTx)
	if err != nil {
		return nil, errors.Wrap(err, "fail to call hex.DecodeString()")
	}
	resTx, err := b.client.DecodeRawTransaction(byteHex)
	if err != nil {
		return nil, errors.Wrap(err, "fail to call btc.client.DecodeRawTransaction()")
	}

	return resTx, nil
}

// GetRawTransactionByHex get tx from hex string
// unused for now
func (b *Bitcoin) GetRawTransactionByHex(strHashTx string) (*btcutil.Tx, error) {

	hashTx, err := chainhash.NewHashFromStr(strHashTx)
	if err != nil {
		return nil, errors.Wrapf(err, "fail to call chainhash.NewHashFromStr(%s)", strHashTx)
	}

	tx, err := b.client.GetRawTransaction(hashTx)
	if err != nil {
		return nil, errors.Wrap(err, "fail to call btc.client.GetRawTransaction(hash)")
	}
	//MsgTx()
	//tx.MsgTx()

	return tx, nil
}

// GetTransactionByTxID get transaction result by txID
func (b *Bitcoin) GetTransactionByTxID(txID string) (*btcjson.GetTransactionResult, error) {
	hashTx, err := chainhash.NewHashFromStr(txID)
	if err != nil {
		return nil, errors.Wrapf(err, "fail to call chainhash.NewHashFromStr(%s)", txID)
	}
	resTx, err := b.client.GetTransaction(hashTx)
	if err != nil {
		return nil, errors.Wrapf(err, "fail to call btc.client.GetTransaction(%s)", hashTx)
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

// GetTxOutByTxID get txOut by txID and index
func (b *Bitcoin) GetTxOutByTxID(txID string, index uint32) (*btcjson.GetTxOutResult, error) {
	hash, err := chainhash.NewHashFromStr(txID)
	if err != nil {
		return nil, errors.Wrapf(err, "fail to call chainhash.NewHashFromStr(%s)", txID)
	}

	// Gettxout / txHash *chainhash.Hash, index uint32, mempool bool
	txOutResult, err := b.client.GetTxOut(hash, index, false)
	if err != nil {
		return nil, errors.Wrapf(err, "fail to call btc.client.GetTxOut(%s, %d, false)", hash, index)
	}

	return txOutResult, nil
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
// [Noted] 手数料を考慮せず、全額送金しようとすると、SendRawTransaction()で、`min relay fee not met`エラーが発生する
func (b *Bitcoin) CreateRawTransaction(receiverAddr string, amount btcutil.Amount, inputs []btcjson.TransactionInput) (*wire.MsgTx, error) {
	//TODO:sendAddrの厳密なチェックがセキュリティ的に必要な場面もありそう
	sendAddrDecoded, err := btcutil.DecodeAddress(receiverAddr, b.GetChainConf())
	if err != nil {
		return nil, errors.Wrapf(err, "fail to call btcutil.DecodeAddress(%s)", receiverAddr)
	}

	// パラメータを作成する
	outputs := make(map[btcutil.Address]btcutil.Amount)
	outputs[sendAddrDecoded] = amount //satoshi

	// CreateRawTransaction
	return b.CreateRawTransactionWithOutput(inputs, outputs)
}

// CreateRawTransactionWithOutput 出金時に1対多で送信する場合に利用するトランザクションを作成する
func (b *Bitcoin) CreateRawTransactionWithOutput(inputs []btcjson.TransactionInput, outputs map[btcutil.Address]btcutil.Amount) (*wire.MsgTx, error) {
	lockTime := int64(0) //TODO:Raw locktime ここは何をいれるべき？

	// CreateRawTransaction
	msgTx, err := b.client.CreateRawTransaction(inputs, outputs, &lockTime)
	if err != nil {
		return nil, errors.Errorf("btcutil.CreateRawTransaction(): error: %s", err)
	}

	return msgTx, nil
}

// FundRawTransaction Add inputs to a transaction until it has enough in value to meet its out value.
// TODO: unused for now, but it looks useful
func (b *Bitcoin) FundRawTransaction(hex string) (*FundRawTransactionResult, error) {
	//fundrawtransaction
	//https://bitcoincore.org/en/doc/0.19.0/rpc/rawtransactions/fundrawtransaction/

	//hex
	bHex, err := json.Marshal(hex)
	if err != nil {
		return nil, errors.Wrap(err, "fail to call json.Marchal(hex)")
	}

	//fee rate
	feePerKb, err := b.EstimateSmartFee()
	if err != nil {
		return nil, errors.Wrap(err, "fail to call btc.EstimateSmartFee()")
	}

	bFeeRate, err := json.Marshal(struct {
		FeeRate float64 `json:"feeRate"`
	}{
		FeeRate: feePerKb,
	})
	if err != nil {
		return nil, errors.Wrap(err, "fail to call json.Marchal(feeRate)")
	}

	rawResult, err := b.client.RawRequest("fundrawtransaction", []json.RawMessage{bHex, bFeeRate})
	if err != nil {
		//error: -4: Insufficient funds
		return nil, errors.Wrap(err, "fail to call json.RawRequest(fundrawtransaction)")
	}

	fundRawTransactionResult := FundRawTransactionResult{}
	err = json.Unmarshal([]byte(rawResult), &fundRawTransactionResult)
	if err != nil {
		return nil, errors.Wrap(err, "fail to call json.Unmarshal(rawResult)")
	}

	return &fundRawTransactionResult, nil
}

// SignRawTransaction sign on raw unsigned tx for `not multisig address` like client account
// - this would be used for receipt action
// - for multisig, refer to `SignRawTransactionWithKey()`
// - TODO: this code can be shared with SignRawTransactionWithKey()
func (b *Bitcoin) SignRawTransaction(tx *wire.MsgTx, prevtxs []PrevTx) (*wire.MsgTx, bool, error) {
	//hex tx
	hexTx, err := b.ToHex(tx)
	if err != nil {
		return nil, false, errors.Wrap(err, "fail to call btc.ToHex(tx)")
	}

	input1, err := json.Marshal(hexTx)
	if err != nil {
		return nil, false, errors.Wrap(err, "fail to call json.Marchal(hexTx)")
	}

	//prevtxs
	input2, err := json.Marshal(prevtxs)
	if err != nil {
		return nil, false, errors.Wrap(err, "fail to call json.Marchal(prevtxs)")
	}

	// call api `signrawtransactionwithwallet`
	rawResult, err := b.client.RawRequest("signrawtransactionwithwallet", []json.RawMessage{input1, input2})
	if err != nil {
		return nil, false, errors.Wrap(err, "fail to call json.RawRequest(signrawtransactionwithwallet)")
	}
	signRawTxResult := SignRawTransactionResult{}
	err = json.Unmarshal([]byte(rawResult), &signRawTxResult)
	if err != nil {
		return nil, false, errors.Wrap(err, "fail to call json.Unmarshal(rawResult)")
	}
	if len(signRawTxResult.Errors) != 0 {
		return nil, false, errors.Errorf("fail to call json.RawRequest(signrawtransactionwithwallet): error: %s", signRawTxResult.Errors[0].Error)
	}

	msgTx, err := b.ToMsgTx(signRawTxResult.Hex)
	if err != nil {
		return nil, false, errors.Wrap(err, "fail to call btc.ToMsgTx(hex)")
	}

	//Debug to compare tx before and after
	if !signRawTxResult.Complete {
		b.logger.Debug("sign is not completed yet")
		b.debugCompareTx(tx, msgTx)
	}

	return msgTx, signRawTxResult.Complete, nil
}

// SignRawTransaction sign on raw unsigned tx for `not multisig address` like client account
// - this would be used for receipt action
// - for multisig, refer to `SignRawTransactionWithKey()`

// SignRawTransactionWithKey sign on raw unsigned tx for `multisig address`
// - for multisig
func (b *Bitcoin) SignRawTransactionWithKey(tx *wire.MsgTx, privKeysWIF []string, prevtxs []PrevTx) (*wire.MsgTx, bool, error) {
	//if b.Version() >= ctype.BTCVer17 {
	//hex tx
	hexTx, err := b.ToHex(tx)
	if err != nil {
		return nil, false, errors.Errorf("BTC.ToHex(tx): error: %s", err)
	}

	input1, err := json.Marshal(hexTx)
	if err != nil {
		return nil, false, errors.Errorf("json.Marchal(txHex): error: %s", err)
	}

	//private keys
	input2, err := json.Marshal(privKeysWIF)
	if err != nil {
		return nil, false, errors.Errorf("json.Marchal(privKeysWIF): error: %s", err)
	}

	//prevtxs
	input3, err := json.Marshal(prevtxs)
	if err != nil {
		return nil, false, errors.Errorf("json.Marchal(prevtxs): error: %s", err)
	}

	rawResult, err := b.client.RawRequest("signrawtransactionwithkey", []json.RawMessage{input1, input2, input3})
	if err != nil {
		return nil, false, errors.Errorf("json.RawRequest(signrawtransactionwithkey): error: %s", err)
	}

	//SignRawTransactionResult
	signRawTxResult := SignRawTransactionResult{}
	err = json.Unmarshal([]byte(rawResult), &signRawTxResult)
	if err != nil {
		return nil, false, errors.Errorf("json.Unmarshal(): error: %s", err)
	}
	//TODO:戻り値のmsgTxがブランクではない、かつ値が初期値と変化がある場合は、このエラーはskip可能
	//	Signature must be zero for failed CHECK(MULTI)SIG operation
	//  =>こちらのエラーはOK
	if len(signRawTxResult.Errors) != 0 {
		if signRawTxResult.Hex == "" || hexTx == signRawTxResult.Hex {
			grok.Value(signRawTxResult)
			return nil, false, errors.Errorf("json.RawRequest(signrawtransactionwithkey): error: %s", signRawTxResult.Errors[0].Error)
		}
		b.logger.Debug("Errors in signRawTxResult", zap.Any("errors", signRawTxResult.Errors[0].Error))
	}

	msgTx, err := b.ToMsgTx(signRawTxResult.Hex)
	if err != nil {
		return nil, false, errors.Errorf("BTC.ToMsgTx(hex): error: %s", err)
	}

	return msgTx, signRawTxResult.Complete, nil
}

// SendTransactionByHex 外部から渡されたバイト列からRawトランザクションを送信する
// オンラインで実行される必要があるため、watchOnlyWallet専用
func (b *Bitcoin) SendTransactionByHex(hex string) (*chainhash.Hash, error) {
	// Hexからトランザクションを取得
	msgTx, err := b.ToMsgTx(hex)
	if err != nil {
		return nil, err
	}

	//送信
	hash, err := b.sendRawTransaction(msgTx)
	//hash, err := b.client.SendRawTransaction(msgTx, true)
	if err != nil {
		return nil, errors.Errorf("BTC.SendRawTransaction(): error: %s", err)
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
		return nil, errors.Errorf("wireTx.Deserialize(): error: %s", err)
	}

	//送信
	hash, err := b.sendRawTransaction(wireTx)
	if err != nil {
		return nil, errors.Errorf("BTC.SendRawTransaction(): error: %v", err)
	}

	return hash, nil
}

// sendRawTransaction Rawトランザクションを送信する
func (b *Bitcoin) sendRawTransaction(tx *wire.MsgTx) (*chainhash.Hash, error) {
	//送信
	hash, err := b.client.SendRawTransaction(tx, true)
	if err != nil {
		//feeを1Satoshiで試してみたら、
		//-26: 66: min relay fee not metが出た
		return nil, errors.Errorf("client.SendRawTransaction(): error: %s", err)
	}

	return hash, nil
}

// Sign 署名を行う without Bitcoin Core [WIP]
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

func (b *Bitcoin) debugCompareTx(tx1, tx2 *wire.MsgTx) {
	b.logger.Debug("compare tx before and after")
	hexTx1, err := b.ToHex(tx1)
	if err != nil {
		b.logger.Debug("fail to call btc.ToHex(tx1)", zap.Error(err))
	}
	hexTx2, err := b.ToHex(tx2)
	if err != nil {
		b.logger.Debug("fail to call btc.ToHex(tx2)", zap.Error(err))
	}
	if hexTx1 == hexTx2 {
		b.logger.Debug("hexTx before is same to hexTx after. Something program is wrong")
	}
}
