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

// SignRawTransactionResult is response type of PRC `signrawtransactionwithwallet`
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

// CreateRawTransaction create raw transaction
//  - for receipt/transfer action
//func (b *Bitcoin) CreateRawTransaction(receiverAddr string, amount btcutil.Amount, txInputs []btcjson.TransactionInput) (*wire.MsgTx, error) {
//	receiverAddrDecoded, err := btcutil.DecodeAddress(receiverAddr, b.GetChainConf())
//	if err != nil {
//		return nil, errors.Wrapf(err, "fail to call btcutil.DecodeAddress(%s)", receiverAddr)
//	}
//
//	txOutputs := make(map[btcutil.Address]btcutil.Amount)
//	txOutputs[receiverAddrDecoded] = amount //satoshi
//
//	// CreateRawTransaction
//	return b.CreateRawTransactionWithOutput(txInputs, txOutputs)
//}

// CreateRawTransaction create raw transaction
//  - for payment action
func (b *Bitcoin) CreateRawTransaction(inputs []btcjson.TransactionInput, outputs map[btcutil.Address]btcutil.Amount) (*wire.MsgTx, error) {
	lockTime := int64(0) //TODO:Raw locktime what value is exactly required??

	// CreateRawTransaction
	msgTx, err := b.client.CreateRawTransaction(inputs, outputs, &lockTime)
	if err != nil {
		return nil, errors.Wrap(err, "fail to call btcutil.CreateRawTransaction()")
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
// - TODO: this code can be shared with SignRawTransactionWithKey() to some extend
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
		grok.Value(signRawTxResult)
		return nil, false, errors.Errorf("result of `signrawtransactionwithwallet` includes error: %s", signRawTxResult.Errors[0].Error)
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

// SignRawTransactionWithKey sign on raw unsigned tx for `multisig address`
// - for multisig
func (b *Bitcoin) SignRawTransactionWithKey(tx *wire.MsgTx, privKeysWIF []string, prevtxs []PrevTx) (*wire.MsgTx, bool, error) {
	//if b.Version() >= ctype.BTCVer17 {

	// hex tx
	hexTx, err := b.ToHex(tx)
	if err != nil {
		return nil, false, errors.Wrap(err, "fail to call btc.ToHex(tx)")
	}

	input1, err := json.Marshal(hexTx)
	if err != nil {
		return nil, false, errors.Wrap(err, "fail to call json.Marchal(txHex)")
	}

	// private keys
	input2, err := json.Marshal(privKeysWIF)
	if err != nil {
		return nil, false, errors.Wrap(err, "fail to call json.Marchal(privKeysWIF)")
	}

	// prevtxs
	input3, err := json.Marshal(prevtxs)
	if err != nil {
		return nil, false, errors.Errorf("fail to call json.Marchal(prevtxs)")
	}

	// call api `signrawtransactionwithkey`
	rawResult, err := b.client.RawRequest("signrawtransactionwithkey", []json.RawMessage{input1, input2, input3})
	if err != nil {
		return nil, false, errors.Wrap(err, "fail to call json.RawRequest(signrawtransactionwithkey)")
	}

	signRawTxResult := SignRawTransactionResult{}
	err = json.Unmarshal([]byte(rawResult), &signRawTxResult)
	if err != nil {
		return nil, false, errors.Wrap(err, "fail to call json.Unmarshal(rawResult)")
	}
	//Note: if signature is not completed yet, error would occur
	// - ignore error if retured msgTx is not blank、and returned hex is changed from given hex as parameter
	//	above error would be like `Signature must be zero for failed CHECK(MULTI)SIG operation`
	if len(signRawTxResult.Errors) != 0 {
		if signRawTxResult.Hex == "" || hexTx == signRawTxResult.Hex {
			grok.Value(signRawTxResult)
			return nil, false, errors.Errorf("result of `signrawtransactionwithwallet` includes error: %s", signRawTxResult.Errors[0].Error)
		}
		b.logger.Warn("result of `signrawtransactionwithwallet` includes error", zap.Any("errors", signRawTxResult.Errors))
	}

	msgTx, err := b.ToMsgTx(signRawTxResult.Hex)
	if err != nil {
		return nil, false, errors.Wrap(err, "fail to call btc.ToMsgTx(hex)")
	}

	return msgTx, signRawTxResult.Complete, nil
}

// SendTransactionByHex send raw transaction by hex string
func (b *Bitcoin) SendTransactionByHex(hex string) (*chainhash.Hash, error) {
	msgTx, err := b.ToMsgTx(hex)
	if err != nil {
		return nil, err
	}

	// send
	hash, err := b.sendRawTransaction(msgTx)
	if err != nil {
		return nil, errors.Wrap(err, "fail to call btc.SendRawTransaction()")
	}

	//txID
	//hash.String()
	return hash, nil
}

// SendTransactionByByte send raw transaction by byte array
func (b *Bitcoin) SendTransactionByByte(rawTx []byte) (*chainhash.Hash, error) {

	//[]byte to wireTx
	wireTx := new(wire.MsgTx)
	r := bytes.NewBuffer(rawTx)

	if err := wireTx.Deserialize(r); err != nil {
		return nil, errors.Wrap(err, "fail to call wireTx.Deserialize()")
	}

	// send
	hash, err := b.sendRawTransaction(wireTx)
	if err != nil {
		return nil, errors.Wrap(err, "fail to call btc.SendRawTransaction()")
	}

	//txID
	//hash.String()
	return hash, nil
}

// sendRawTransaction send raw transaction
func (b *Bitcoin) sendRawTransaction(tx *wire.MsgTx) (*chainhash.Hash, error) {
	// send
	hash, err := b.client.SendRawTransaction(tx, true)
	if err != nil {
		// error occurred when trying to send tx with minimum fee(1Satoshi)
		//  -26: 66: min relay fee not met
		return nil, errors.Wrap(err, "fail to call btc.client.SendRawTransaction()")
	}

	return hash, nil
}

// Sign sign on unsigned tx
// [WIP] implement signiture without bitcoin core. this code is not fixed yet
// if implementation is done, bitcoin core is not required anymore in keygen/sign wallet
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
	//TODO: how to check isSigned
	//TODO: it's difficult if each TxIn has different Key

	// to sign
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
