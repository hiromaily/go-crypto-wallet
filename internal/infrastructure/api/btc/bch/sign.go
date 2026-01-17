package bch

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/btcsuite/btcd/wire"
	bchtxscript "github.com/gcash/bchd/txscript"
	bchwire "github.com/gcash/bchd/wire"
	"github.com/gcash/bchutil"

	dtobtc "github.com/hiromaily/go-crypto-wallet/internal/application/dto/btc"
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
// IMPORTANT: This method always returns an error because BCH signing requires
// the input amounts for BIP143-like signature hash calculation with SIGHASH_FORKID.
// The standard Sign() interface does not provide this information.
//
// Use one of these alternatives instead:
//   - SignWithPrevTxs(): When you have prevTx information with amounts
//   - SignRawTransactionWithKey(): RPC-based signing (recommended for multisig)
func (*BitcoinCash) Sign(_ *wire.MsgTx, _ string) (string, error) {
	// BCH signing requires input amounts for BIP143-like signature hash calculation.
	// Since this interface doesn't provide amounts, we must return an error.
	// Using a zero amount would produce an invalid signature.
	return "", errors.New("BCH Sign() is not supported: BCH signing requires input amounts for SIGHASH_FORKID. " +
		"Use SignWithPrevTxs() or SignRawTransactionWithKey() instead")
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
	// Serialize btcd tx to a buffer
	var buf bytes.Buffer
	if err := btcdTx.Serialize(&buf); err != nil {
		return nil, fmt.Errorf("fail to serialize btcd tx: %w", err)
	}

	// Deserialize from the buffer as a bchd tx
	bchTx := bchwire.NewMsgTx(bchwire.TxVersion)
	if err := bchTx.Deserialize(&buf); err != nil {
		return nil, fmt.Errorf("fail to deserialize as bchd tx: %w", err)
	}

	return bchTx, nil
}
