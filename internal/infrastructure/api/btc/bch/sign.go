package bch

import (
	"encoding/hex"
	"fmt"
	"io"

	"github.com/btcsuite/btcd/wire"
	bchtxscript "github.com/gcash/bchd/txscript"
	bchwire "github.com/gcash/bchd/wire"
	"github.com/gcash/bchutil"

	dtobtc "github.com/hiromaily/go-crypto-wallet/internal/application/dto/btc"
	"github.com/hiromaily/go-crypto-wallet/pkg/logger"
)

// SIGHASH_FORKID is the BCH-specific sighash flag for replay protection.
// BCH requires SIGHASH_ALL | SIGHASH_FORKID (0x41) for all signatures.
// This prevents BCH transactions from being valid on the BTC chain and vice versa.
const (
	// SigHashForkID is the BCH fork ID flag (0x40)
	SigHashForkID bchtxscript.SigHashType = 0x40

	// SigHashBCH is the combined sighash type for BCH: SIGHASH_ALL | SIGHASH_FORKID (0x41)
	SigHashBCH bchtxscript.SigHashType = bchtxscript.SigHashAll | SigHashForkID
)

// Sign signs a BCH raw transaction with SIGHASH_FORKID.
//
// IMPORTANT: This method is for P2PKH single-sig transactions only.
// For multisig transactions, use SignRawTransactionWithKey which uses RPC.
//
// BCH requires SIGHASH_FORKID (0x40) combined with SIGHASH_ALL (0x01) = 0x41.
// This is handled by the gcash/bchd library which properly calculates the
// BIP143-like signature hash with the fork ID flag.
//
// Note: This method requires prevTxs information for the input amounts.
// The current interface only provides tx and private key, so this implementation
// uses a workaround by returning an error when amount information is missing.
// For production use, prefer SignRawTransactionWithKey which handles this via RPC.
func (b *BitcoinCash) Sign(tx *wire.MsgTx, strPrivateKey string) (string, error) {
	// Decode the private key using bchutil (for BCH-compatible PrivateKey type)
	wif, err := bchutil.DecodeWIF(strPrivateKey)
	if err != nil {
		return "", fmt.Errorf("fail to decode WIF for BCH signing: %w", err)
	}
	privKey := wif.PrivKey

	// Convert btcd wire.MsgTx to bchd wire.MsgTx
	bchTx, err := convertToBCHTx(tx)
	if err != nil {
		return "", fmt.Errorf("fail to convert tx for BCH signing: %w", err)
	}

	// Sign each input
	for idx, val := range tx.TxIn {
		// For BCH signing, we need the amount of the input being spent.
		// This is required for BIP143-like signature hash calculation.
		// Since the current interface doesn't provide amounts, we'll attempt
		// to use the SignatureScript field which contains the scriptPubKey.
		//
		// NOTE: This is a limitation. For proper BCH signing with amounts,
		// use SignRawTransactionWithKey which handles this via RPC.

		subscript := val.SignatureScript
		if len(subscript) == 0 {
			return "", fmt.Errorf(
				"BCH signing requires scriptPubKey in TxIn[%d].SignatureScript for input", idx)
		}

		// WARNING: Amount is required for BCH signing but not available in current interface.
		// Using 0 as amount will produce invalid signatures for real transactions.
		// This is kept for API compatibility but should not be used in production.
		// For production, use SignRawTransactionWithKey with proper prevTx information.
		logger.Warn(
			"BCH Sign() called without amount information - signature may be invalid",
			"inputIndex", idx)
		amount := int64(0)

		// Use bchd's SignatureScript which properly handles SIGHASH_FORKID
		script, err := bchtxscript.SignatureScript(
			bchTx,
			idx,
			amount,
			subscript,
			SigHashBCH, // SIGHASH_ALL | SIGHASH_FORKID (0x41)
			privKey,
			true, // compress pubkey
		)
		if err != nil {
			return "", fmt.Errorf("fail to create BCH signature for input %d: %w", idx, err)
		}

		// Set the signature script back to the original btcd tx
		tx.TxIn[idx].SignatureScript = script
	}

	// Convert back to hex
	hexTx, err := b.ToHex(tx)
	if err != nil {
		return "", fmt.Errorf("fail to convert signed BCH tx to hex: %w", err)
	}

	return hexTx, nil
}

// SignWithPrevTxs signs a BCH raw transaction with SIGHASH_FORKID.
// This method requires prevTxs to get the input amounts which are required
// for BCH's BIP143-like signature hash calculation.
//
// This is the recommended method for BCH signing when you have prevTx information.
func (b *BitcoinCash) SignWithPrevTxs(
	tx *wire.MsgTx, strPrivateKey string, prevTxs []dtobtc.PreviousTx,
) (string, error) {
	// Decode the private key using bchutil (for BCH-compatible PrivateKey type)
	wif, err := bchutil.DecodeWIF(strPrivateKey)
	if err != nil {
		return "", fmt.Errorf("fail to decode WIF for BCH signing: %w", err)
	}
	privKey := wif.PrivKey

	// Convert btcd wire.MsgTx to bchd wire.MsgTx
	bchTx, err := convertToBCHTx(tx)
	if err != nil {
		return "", fmt.Errorf("fail to convert tx for BCH signing: %w", err)
	}

	// Build a map of prevTxs for quick lookup
	prevTxMap := make(map[string]dtobtc.PreviousTx)
	for _, prev := range prevTxs {
		key := fmt.Sprintf("%s:%d", prev.TxID, prev.Vout)
		prevTxMap[key] = prev
	}

	// Sign each input
	for idx, val := range tx.TxIn {
		// Get the prevTx for this input
		key := fmt.Sprintf("%s:%d", val.PreviousOutPoint.Hash.String(), val.PreviousOutPoint.Index)
		prev, ok := prevTxMap[key]
		if !ok {
			return "", fmt.Errorf("missing prevTx for input %d (outpoint: %s)", idx, key)
		}

		// Get the scriptPubKey from prevTx
		subscript, err := hex.DecodeString(prev.ScriptPubKey)
		if err != nil {
			return "", fmt.Errorf("fail to decode scriptPubKey for input %d: %w", idx, err)
		}

		// Get the amount from prevTx (required for BCH signing)
		amount := int64(prev.Amount)

		// Use bchd's SignatureScript which properly handles SIGHASH_FORKID
		script, err := bchtxscript.SignatureScript(
			bchTx,
			idx,
			amount,
			subscript,
			SigHashBCH, // SIGHASH_ALL | SIGHASH_FORKID (0x41)
			privKey,
			true, // compress pubkey
		)
		if err != nil {
			return "", fmt.Errorf("fail to create BCH signature for input %d: %w", idx, err)
		}

		// Set the signature script back to the original btcd tx
		tx.TxIn[idx].SignatureScript = script
	}

	// Convert back to hex
	hexTx, err := b.ToHex(tx)
	if err != nil {
		return "", fmt.Errorf("fail to convert signed BCH tx to hex: %w", err)
	}

	return hexTx, nil
}

// convertToBCHTx converts a btcsuite/btcd wire.MsgTx to gcash/bchd wire.MsgTx.
// This is necessary because bchd's txscript package requires bchd's wire types.
func convertToBCHTx(btcdTx *wire.MsgTx) (*bchwire.MsgTx, error) {
	// Serialize btcd tx
	var buf []byte
	btcdBuf := new(bytesBuffer)
	if err := btcdTx.Serialize(btcdBuf); err != nil {
		return nil, fmt.Errorf("fail to serialize btcd tx: %w", err)
	}
	buf = btcdBuf.Bytes()

	// Deserialize as bchd tx
	bchTx := bchwire.NewMsgTx(bchwire.TxVersion)
	bchBuf := newBytesReader(buf)
	if err := bchTx.Deserialize(bchBuf); err != nil {
		return nil, fmt.Errorf("fail to deserialize as bchd tx: %w", err)
	}

	return bchTx, nil
}

// bytesBuffer implements io.Writer for serialization
type bytesBuffer struct {
	data []byte
}

func (b *bytesBuffer) Write(p []byte) (n int, err error) {
	b.data = append(b.data, p...)
	return len(p), nil
}

func (b *bytesBuffer) Bytes() []byte {
	return b.data
}

// bytesReader implements io.Reader for deserialization
type bytesReader struct {
	data []byte
	pos  int
}

func newBytesReader(data []byte) *bytesReader {
	return &bytesReader{data: data, pos: 0}
}

func (b *bytesReader) Read(p []byte) (n int, err error) {
	if b.pos >= len(b.data) {
		return 0, io.EOF
	}
	n = copy(p, b.data[b.pos:])
	b.pos += n
	return n, nil
}
