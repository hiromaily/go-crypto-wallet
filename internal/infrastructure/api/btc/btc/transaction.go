package btc

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/bookerzzz/grok"
	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"

	dtobtc "github.com/hiromaily/go-crypto-wallet/internal/application/dto/btc"
	btcrpc "github.com/hiromaily/go-crypto-wallet/pkg/chains/btc/rpc"
	"github.com/hiromaily/go-crypto-wallet/pkg/logger"
)

// refer to https://www.haowuliaoa.com/article/info/11350.html (Chinese site)

// naming regulation for transaction
//  - as example, *wire.MsgTx is that tx is used as suffix
//  - so tx should be named like hexTx, hashTx

// ToHex convert wire.MsgTx to string of Hexadecimal(16進数)
func (*Bitcoin) ToHex(tx *wire.MsgTx) (string, error) {
	buf := bytes.NewBuffer(make([]byte, 0, tx.SerializeSize()))
	if err := tx.Serialize(buf); err != nil {
		return "", fmt.Errorf("fail to call tx.Serialize(): %w", err)
	}
	return hex.EncodeToString(buf.Bytes()), nil
}

// ToMsgTx convert string of Hexadecimal(16進数) to wire.MsgTx
func (*Bitcoin) ToMsgTx(txHex string) (*wire.MsgTx, error) {
	byteHex, err := hex.DecodeString(txHex)
	if err != nil {
		return nil, fmt.Errorf("fail to call hex.DecodeString(): %w", err)
	}

	var msgTx wire.MsgTx
	if err = msgTx.Deserialize(bytes.NewReader(byteHex)); err != nil {
		return nil, fmt.Errorf("fail to call hex.Deserialize(): %w", err)
	}

	return &msgTx, nil
}

// DecodeRawTransaction decodes a serialized transaction hex string back into a JSON object
func (b *Bitcoin) DecodeRawTransaction(hexTx string) (*btcrpc.TxRawResult, error) {
	return b.pkgrpc.DecodeRawTransaction(hexTx)
}

// GetTransactionByTxID get transaction result by txID
func (b *Bitcoin) GetTransactionByTxID(txID string) (*btcrpc.GetTransactionResult, error) {
	result, err := b.pkgrpc.GetTransaction(txID)
	if err != nil {
		return nil, fmt.Errorf("fail to call btcrpc.GetTransaction(%s): %w", txID, err)
	}

	return result, nil
}

// GetTxOutByTxID get txOut by txID and index
//   - it's not used anywhere
func (b *Bitcoin) GetTxOutByTxID(txID string, index uint32) (*btcjson.GetTxOutResult, error) {
	hash, err := chainhash.NewHashFromStr(txID)
	if err != nil {
		return nil, fmt.Errorf("fail to call chainhash.NewHashFromStr(%s): %w", txID, err)
	}

	// Gettxout / txHash *chainhash.Hash, index uint32, mempool bool
	//  client.GetTxOut is not outdated yet (at bitcoin core 0.19)
	txOutResult, err := b.btcdClient.GetTxOut(hash, index, false)
	if err != nil {
		return nil, fmt.Errorf("fail to call btc.client.GetTxOut(%s, %d, false): %w", hash, index, err)
	}

	return txOutResult, nil
	// value *GetTxOutResult = {
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

// GetRawTransactionByHex get tx from hex string
// unused for now
func (b *Bitcoin) GetRawTransactionByHex(strHashTx string) (*btcutil.Tx, error) {
	hashTx, err := chainhash.NewHashFromStr(strHashTx)
	if err != nil {
		return nil, fmt.Errorf("fail to call chainhash.NewHashFromStr(%s): %w", strHashTx, err)
	}

	tx, err := b.btcdClient.GetRawTransaction(hashTx)
	if err != nil {
		return nil, fmt.Errorf("fail to call btc.client.GetRawTransaction(hash): %w", err)
	}
	// MsgTx()
	// tx.MsgTx()

	return tx, nil
}

// CreateRawTransaction create raw transaction
//   - for payment action
func (b *Bitcoin) CreateRawTransaction(
	inputs []btcjson.TransactionInput, outputs map[btcutil.Address]btcutil.Amount,
) (*wire.MsgTx, error) {
	lockTime := int64(0) // TODO:Raw locktime what value is exactly required??

	// CreateRawTransaction
	msgTx, err := b.btcdClient.CreateRawTransaction(inputs, outputs, &lockTime)
	if err != nil {
		return nil, fmt.Errorf("fail to call btcutil.CreateRawTransaction(): %w", err)
	}

	return msgTx, nil
}

// FundRawTransaction Add inputs to a transaction until it has enough in value to meet its out value.
// TODO: unused for now, but it looks useful
func (b *Bitcoin) FundRawTransaction(hexTx string) (*btcrpc.FundRawTransactionResult, error) {
	return b.FundRawTransactionWithOptions(hexTx, nil)
}

// FundRawTransactionWithOptions adds inputs to a transaction with custom options.
// Bitcoin Core v28.0+ supports max_tx_weight parameter to limit transaction size.
func (b *Bitcoin) FundRawTransactionWithOptions(
	hexTx string,
	opts *btcrpc.FundRawTransactionOptions,
) (*btcrpc.FundRawTransactionResult, error) {
	if opts == nil {
		opts = &btcrpc.FundRawTransactionOptions{}
	}

	// Set fee rate if not specified
	if opts.FeeRate == 0 {
		feePerKb, err := b.pkgrpc.EstimateSmartFee(int(b.confirmationBlock))
		if err != nil {
			return nil, fmt.Errorf("fail to call btc.EstimateSmartFee(): %w", err)
		}
		opts.FeeRate = feePerKb
	}

	result, err := b.pkgrpc.FundRawTransaction(hexTx, opts)
	if err != nil {
		// error: -4: Insufficient funds
		return nil, fmt.Errorf("fail to call btcrpc.FundRawTransaction(): %w", err)
	}

	return result, nil
}

// prepareSignRawTransactionInputs prepares common inputs for sign raw transaction operations
func (b *Bitcoin) prepareSignRawTransactionInputs(
	tx *wire.MsgTx, prevtxs []dtobtc.PreviousTx,
) (hexTx string, pkgPrevTxs []btcrpc.PrevTx, err error) {
	pkgPrevTxs = fromDTOPreviousTxToLocal(prevtxs)

	hexTx, err = b.ToHex(tx)
	if err != nil {
		return "", nil, fmt.Errorf("fail to call btc.ToHex(tx): %w", err)
	}

	return hexTx, pkgPrevTxs, nil
}

// processSignRawTransactionResult processes the result from sign raw transaction operations
func (b *Bitcoin) processSignRawTransactionResult(
	tx *wire.MsgTx, hexTx string, result *btcrpc.SignRawTransactionResult, strictErrorCheck bool,
) (*wire.MsgTx, bool, error) {
	if len(result.Errors) != 0 {
		if strictErrorCheck {
			grok.Value(result)
			return nil, false, fmt.Errorf(
				"result of sign raw transaction includes error: %s",
				result.Errors[0].Error)
		}
		// For SignRawTransactionWithKey: allow partial signatures
		if result.Hex == "" || hexTx == result.Hex {
			grok.Value(result)
			return nil, false, fmt.Errorf(
				"result of sign raw transaction includes error: %s",
				result.Errors[0].Error)
		}
		logger.Warn(
			"result of sign raw transaction includes error",
			"errors", result.Errors)
	}

	msgTx, err := b.ToMsgTx(result.Hex)
	if err != nil {
		return nil, false, fmt.Errorf("fail to call btc.ToMsgTx(hex): %w", err)
	}

	// Debug to compare tx before and after
	if !result.Complete {
		logger.Debug("sign is not completed yet")
		b.debugCompareTx(tx, msgTx)
	}

	return msgTx, result.Complete, nil
}

// SignRawTransaction signs a raw unsigned transaction for non-multisig addresses (e.g., client account).
// This function supports all Bitcoin address types including Taproot with automatic signature type selection:
//   - P2PKH (Legacy): ECDSA signature
//   - P2SH-SegWit: ECDSA signature with witness data
//   - P2WPKH (Bech32): ECDSA signature with witness data
//   - P2TR (Taproot): Schnorr signature (BIP340) with witness data
//
// The signature type is automatically determined by Bitcoin Core based on the scriptPubKey type
// of the input being spent. Bitcoin Core v22.0+ is required for Taproot/Schnorr support.
//
// This would be used for deposit action (client -> deposit).
// For multisig addresses, refer to SignRawTransactionWithKey().
func (b *Bitcoin) SignRawTransaction(tx *wire.MsgTx, prevtxs []dtobtc.PreviousTx) (*wire.MsgTx, bool, error) {
	hexTx, pkgPrevTxs, err := b.prepareSignRawTransactionInputs(tx, prevtxs)
	if err != nil {
		return nil, false, err
	}

	pkgResult, err := b.pkgrpc.SignRawTransactionWithWallet(hexTx, pkgPrevTxs)
	if err != nil {
		return nil, false, fmt.Errorf("fail to call btcrpc.SignRawTransactionWithWallet(): %w", err)
	}

	return b.processSignRawTransactionResult(tx, hexTx, pkgResult, true)
}

// SignRawTransactionWithKey signs a raw unsigned transaction for multisig addresses.
// This function supports all Bitcoin address types including Taproot with automatic signature type selection:
//   - P2SH (Legacy multisig): ECDSA signatures
//   - P2SH-SegWit: ECDSA signatures with witness data
//   - P2WSH (Bech32 multisig): ECDSA signatures with witness data
//   - P2TR (Taproot): Schnorr signatures (BIP340) - Note: Taproot multisig uses MuSig2 or script path
//
// The signature type is automatically determined by Bitcoin Core based on the scriptPubKey type
// of the input being spent. Bitcoin Core v22.0+ is required for Taproot/Schnorr support.
//
// For Taproot, note that native Taproot multisig would use MuSig2 (key path) or script path spending,
// which is different from traditional multisig and may require additional implementation.
//
// This would be used for payment and transfer actions from multisig accounts.
func (b *Bitcoin) SignRawTransactionWithKey(
	tx *wire.MsgTx, privKeysWIF []string, prevtxs []dtobtc.PreviousTx,
) (*wire.MsgTx, bool, error) {
	hexTx, pkgPrevTxs, err := b.prepareSignRawTransactionInputs(tx, prevtxs)
	if err != nil {
		return nil, false, err
	}

	pkgResult, err := b.pkgrpc.SignRawTransactionWithKey(hexTx, privKeysWIF, pkgPrevTxs)
	if err != nil {
		return nil, false, fmt.Errorf("fail to call btcrpc.SignRawTransactionWithKey(): %w", err)
	}

	// Note: For multisig, partial signatures are allowed
	return b.processSignRawTransactionResult(tx, hexTx, pkgResult, false)
}

// SendTransactionByHex send raw transaction by hex string
func (b *Bitcoin) SendTransactionByHex(hexTx string) (*chainhash.Hash, error) {
	msgTx, err := b.ToMsgTx(hexTx)
	if err != nil {
		return nil, err
	}

	// send
	hash, err := b.sendRawTransaction(msgTx)
	if err != nil {
		return nil, fmt.Errorf("fail to call btc.SendRawTransaction(): %w", err)
	}

	// txID
	// hash.String()
	return hash, nil
}

// SendTransactionByByte send raw transaction by byte array
func (b *Bitcoin) SendTransactionByByte(rawTx []byte) (*chainhash.Hash, error) {
	// []byte to wireTx
	wireTx := new(wire.MsgTx)
	r := bytes.NewBuffer(rawTx)

	if err := wireTx.Deserialize(r); err != nil {
		return nil, fmt.Errorf("fail to call wireTx.Deserialize(): %w", err)
	}

	// send
	hash, err := b.sendRawTransaction(wireTx)
	if err != nil {
		return nil, fmt.Errorf("fail to call btc.SendRawTransaction(): %w", err)
	}

	// txID
	// hash.String()
	return hash, nil
}

// sendRawTransaction send raw transaction
func (b *Bitcoin) sendRawTransaction(tx *wire.MsgTx) (*chainhash.Hash, error) {
	// send
	hash, err := b.btcdClient.SendRawTransaction(tx, true)
	if err != nil {
		// error occurred when trying to send tx with minimum fee(1Satoshi)
		//  -26: 66: min relay fee not met
		if strings.Contains(err.Error(), "-27: Transaction already in block chain") {
			logger.Info("Transaction already in block chain")
			return nil, nil
		}
		return nil, fmt.Errorf("fail to call btc.client.SendRawTransaction(): %w", err)
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
		var script []byte
		script, err = txscript.SignatureScript(tx, idx, val.SignatureScript, txscript.SigHashAll, privKey, false)
		if err != nil {
			return "", err
		}
		tx.TxIn[idx].SignatureScript = script
	}
	// TODO: how to check isSigned
	// TODO: it's difficult if each TxIn has different Key

	// to sign
	hexTx, err := b.ToHex(tx)
	if err != nil {
		return "", err
	}

	return hexTx, nil
}

func (b *Bitcoin) debugCompareTx(tx1, tx2 *wire.MsgTx) {
	logger.Debug("compare tx before and after")
	hexTx1, err := b.ToHex(tx1)
	if err != nil {
		logger.Debug("fail to call btc.ToHex(tx1)", "error", err)
	}
	hexTx2, err := b.ToHex(tx2)
	if err != nil {
		logger.Debug("fail to call btc.ToHex(tx2)", "error", err)
	}
	if hexTx1 == hexTx2 {
		logger.Debug("hexTx before is same to hexTx after. Something program is wrong")
	}
}
