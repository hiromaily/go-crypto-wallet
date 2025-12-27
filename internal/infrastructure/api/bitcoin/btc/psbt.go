package btc

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"

	"github.com/btcsuite/btcd/btcec/v2/ecdsa"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/btcutil/psbt"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"

	"github.com/hiromaily/go-crypto-wallet/pkg/logger"
)

// PSBTInput represents input metadata for PSBT
type PSBTInput struct {
	WitnessUTXO  *wire.TxOut
	RedeemScript []byte
	SighashType  txscript.SigHashType
}

// PSBTOutput represents output metadata for PSBT
type PSBTOutput struct {
	RedeemScript []byte
}

// ParsedPSBT represents a parsed PSBT with metadata
type ParsedPSBT struct {
	Packet       *psbt.Packet
	InputCount   int
	OutputCount  int
	IsComplete   bool
	HasSignature bool
}

// CreatePSBT creates a PSBT from an unsigned transaction with metadata.
// This function adds all necessary metadata (witness UTXO, redeem scripts) for offline signing.
// Used by Watch wallet to create unsigned PSBTs.
func (b *Bitcoin) CreatePSBT(msgTx *wire.MsgTx, prevTxs []PrevTx) (string, error) {
	// Create PSBT from unsigned transaction
	packet, err := psbt.NewFromUnsignedTx(msgTx)
	if err != nil {
		return "", fmt.Errorf("failed to create PSBT from transaction: %w", err)
	}

	// Create updater to add metadata
	updater, err := psbt.NewUpdater(packet)
	if err != nil {
		return "", fmt.Errorf("failed to create PSBT updater: %w", err)
	}

	// Add metadata for each input from prevTxs
	for i, prevTx := range prevTxs {
		if i >= len(packet.UnsignedTx.TxIn) {
			return "", fmt.Errorf("prevTxs index %d exceeds number of inputs %d", i, len(packet.UnsignedTx.TxIn))
		}

		// Add witness UTXO (required for SegWit/Taproot signing)
		amount, err := btcutil.NewAmount(prevTx.Amount)
		if err != nil {
			return "", fmt.Errorf("failed to parse amount for input %d: %w", i, err)
		}

		scriptPubKey, err := b.decodeHexScript(prevTx.ScriptPubKey)
		if err != nil {
			return "", fmt.Errorf("failed to decode scriptPubKey for input %d: %w", i, err)
		}

		witnessUTXO := &wire.TxOut{
			Value:    int64(amount),
			PkScript: scriptPubKey,
		}

		if err := updater.AddInWitnessUtxo(witnessUTXO, i); err != nil {
			return "", fmt.Errorf("failed to add witness UTXO for input %d: %w", i, err)
		}

		// Add redeem script for P2SH/P2WSH if provided
		if prevTx.RedeemScript != "" {
			redeemScript, err := b.decodeHexScript(prevTx.RedeemScript)
			if err != nil {
				return "", fmt.Errorf("failed to decode redeemScript for input %d: %w", i, err)
			}
			if err := updater.AddInRedeemScript(redeemScript, i); err != nil {
				return "", fmt.Errorf("failed to add redeem script for input %d: %w", i, err)
			}
		}

		// Add sighash type (default to SIGHASH_ALL)
		if err := updater.AddInSighashType(txscript.SigHashAll, i); err != nil {
			return "", fmt.Errorf("failed to add sighash type for input %d: %w", i, err)
		}
	}

	// Serialize PSBT to base64
	psbtBase64, err := b.serializePSBT(packet)
	if err != nil {
		return "", fmt.Errorf("failed to serialize PSBT: %w", err)
	}

	logger.Debug("Created PSBT from transaction",
		"inputs", len(msgTx.TxIn),
		"outputs", len(msgTx.TxOut),
		"txid", msgTx.TxHash().String())

	return psbtBase64, nil
}

// ParsePSBT parses a base64-encoded PSBT and returns metadata.
// Used by all wallets to read PSBT files.
func (b *Bitcoin) ParsePSBT(psbtBase64 string) (*ParsedPSBT, error) {
	// Decode base64 to bytes
	psbtBytes, err := base64.StdEncoding.DecodeString(psbtBase64)
	if err != nil {
		return nil, fmt.Errorf("failed to decode base64 PSBT: %w", err)
	}

	// Parse PSBT using btcd package
	packet, err := psbt.NewFromRawBytes(bytes.NewReader(psbtBytes), false)
	if err != nil {
		return nil, fmt.Errorf("failed to parse PSBT: %w", err)
	}

	// Check if PSBT has any signatures
	hasSignature := b.hasPartialSignatures(packet)

	parsed := &ParsedPSBT{
		Packet:       packet,
		InputCount:   len(packet.Inputs),
		OutputCount:  len(packet.Outputs),
		IsComplete:   packet.IsComplete(),
		HasSignature: hasSignature,
	}

	logger.Debug("Parsed PSBT",
		"inputs", parsed.InputCount,
		"outputs", parsed.OutputCount,
		"complete", parsed.IsComplete,
		"hasSignature", parsed.HasSignature)

	return parsed, nil
}

// ValidatePSBT validates a PSBT structure and checks BIP174 compliance.
// Used by all wallets to verify PSBT before processing.
func (b *Bitcoin) ValidatePSBT(psbtBase64 string) error {
	parsed, err := b.ParsePSBT(psbtBase64)
	if err != nil {
		return fmt.Errorf("failed to parse PSBT for validation: %w", err)
	}

	// Validate that packet is not nil
	if parsed.Packet == nil {
		return errors.New("PSBT packet is nil")
	}

	// Validate that unsigned transaction exists
	if parsed.Packet.UnsignedTx == nil {
		return errors.New("PSBT unsigned transaction is nil")
	}

	// Validate input count matches
	if len(parsed.Packet.Inputs) != len(parsed.Packet.UnsignedTx.TxIn) {
		return fmt.Errorf("PSBT input count mismatch: %d inputs vs %d TxIn",
			len(parsed.Packet.Inputs), len(parsed.Packet.UnsignedTx.TxIn))
	}

	// Validate output count matches
	if len(parsed.Packet.Outputs) != len(parsed.Packet.UnsignedTx.TxOut) {
		return fmt.Errorf("PSBT output count mismatch: %d outputs vs %d TxOut",
			len(parsed.Packet.Outputs), len(parsed.Packet.UnsignedTx.TxOut))
	}

	// Validate each input has witness UTXO for SegWit/Taproot
	for i, input := range parsed.Packet.Inputs {
		if input.WitnessUtxo == nil {
			return fmt.Errorf("input %d missing witness UTXO (required for SegWit/Taproot)", i)
		}
	}

	logger.Debug("PSBT validation passed",
		"inputs", parsed.InputCount,
		"outputs", parsed.OutputCount)

	return nil
}

// SignPSBTWithKey signs a PSBT with provided private keys (offline).
// This function works completely offline without Bitcoin Core RPC.
// Used by Keygen and Sign wallets for air-gapped signing.
//
// Returns:
//   - psbtBase64: The signed PSBT in base64 format
//   - isComplete: true if all signatures are collected (ready for finalization)
//   - error: any error that occurred
func (b *Bitcoin) SignPSBTWithKey(psbtBase64 string, wifs []string) (string, bool, error) {
	// Parse PSBT
	parsed, err := b.ParsePSBT(psbtBase64)
	if err != nil {
		return "", false, fmt.Errorf("failed to parse PSBT for signing: %w", err)
	}

	// Decode WIF private keys
	privKeys := make([]*btcutil.WIF, 0, len(wifs))
	for _, wif := range wifs {
		privKey, err := btcutil.DecodeWIF(wif)
		if err != nil {
			return "", false, fmt.Errorf("failed to decode WIF private key: %w", err)
		}
		privKeys = append(privKeys, privKey)
	}

	// Create updater for signing
	updater, err := psbt.NewUpdater(parsed.Packet)
	if err != nil {
		return "", false, fmt.Errorf("failed to create updater for signing: %w", err)
	}

	// Sign each input with each provided key
	signedCount := 0
	for i := range parsed.Packet.UnsignedTx.TxIn {
		// Get witness UTXO for this input
		witnessUtxo := parsed.Packet.Inputs[i].WitnessUtxo
		if witnessUtxo == nil {
			logger.Warn("Skipping input without witness UTXO", "input", i)
			continue
		}

		// Try signing with each private key
		for _, privKey := range privKeys {
			// Create signature hash
			sigHashes := txscript.NewTxSigHashes(parsed.Packet.UnsignedTx, nil)
			hash, err := txscript.CalcWitnessSigHash(
				witnessUtxo.PkScript,
				sigHashes,
				txscript.SigHashAll,
				parsed.Packet.UnsignedTx,
				i,
				witnessUtxo.Value,
			)
			if err != nil {
				logger.Warn("Failed to calculate signature hash", "input", i, "error", err)
				continue
			}

			// Sign the hash
			signature := ecdsa.Sign(privKey.PrivKey, hash)

			// Serialize signature with sighash type
			sigBytes := append(signature.Serialize(), byte(txscript.SigHashAll))

			// Add partial signature to PSBT
			pubKey := privKey.PrivKey.PubKey().SerializeCompressed()
			outcome, err := updater.Sign(i, sigBytes, pubKey, nil, nil)
			if err != nil {
				// This may fail if the key doesn't match this input, which is normal
				logger.Debug("Signature not applicable for this input", "input", i, "error", err)
				continue
			}

			// Check outcome: 0 = success, 1 = already finalized, -1 = invalid
			if outcome == psbt.SignSuccesful {
				signedCount++
				logger.Debug("Added signature to input", "input", i)
			} else {
				logger.Debug("Signature not added", "input", i, "outcome", outcome)
			}
		}
	}

	if signedCount == 0 {
		return "", false, errors.New("no signatures were added (keys may not match PSBT inputs)")
	}

	// Check if PSBT is now complete
	isComplete := parsed.Packet.IsComplete()

	// Serialize signed PSBT to base64
	signedPSBT, err := b.serializePSBT(parsed.Packet)
	if err != nil {
		return "", false, fmt.Errorf("failed to serialize signed PSBT: %w", err)
	}

	logger.Debug("PSBT signing completed",
		"signedCount", signedCount,
		"isComplete", isComplete)

	return signedPSBT, isComplete, nil
}

// FinalizePSBT finalizes a fully signed PSBT, converting partial signatures to final scriptSig/witness.
// This function should only be called when PSBT is complete (all signatures collected).
// Used by Watch wallet before extracting the final transaction.
func (b *Bitcoin) FinalizePSBT(psbtBase64 string) (string, error) {
	// Parse PSBT
	parsed, err := b.ParsePSBT(psbtBase64)
	if err != nil {
		return "", fmt.Errorf("failed to parse PSBT for finalization: %w", err)
	}

	// Check if PSBT is complete
	if !parsed.IsComplete {
		return "", errors.New("cannot finalize incomplete PSBT (missing signatures)")
	}

	// Finalize all inputs
	for i := range parsed.Packet.UnsignedTx.TxIn {
		if err := psbt.Finalize(parsed.Packet, i); err != nil {
			return "", fmt.Errorf("failed to finalize input %d: %w", i, err)
		}
	}

	// Serialize finalized PSBT to base64
	finalizedPSBT, err := b.serializePSBT(parsed.Packet)
	if err != nil {
		return "", fmt.Errorf("failed to serialize finalized PSBT: %w", err)
	}

	logger.Debug("PSBT finalization completed",
		"inputs", len(parsed.Packet.UnsignedTx.TxIn))

	return finalizedPSBT, nil
}

// ExtractTransaction extracts the final signed transaction from a finalized PSBT.
// This should only be called after FinalizePSBT.
// Used by Watch wallet to get the transaction ready for broadcasting.
func (b *Bitcoin) ExtractTransaction(psbtBase64 string) (*wire.MsgTx, error) {
	// Parse PSBT
	parsed, err := b.ParsePSBT(psbtBase64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse PSBT for extraction: %w", err)
	}

	// Extract final transaction from finalized PSBT
	finalTx, err := psbt.Extract(parsed.Packet)
	if err != nil {
		return nil, fmt.Errorf("failed to extract transaction from PSBT: %w", err)
	}

	logger.Debug("Extracted final transaction from PSBT",
		"txid", finalTx.TxHash().String(),
		"size", finalTx.SerializeSize())

	return finalTx, nil
}

// IsPSBTComplete checks if a PSBT has all required signatures.
// Used to determine if PSBT is ready for finalization.
func (b *Bitcoin) IsPSBTComplete(psbtBase64 string) (bool, error) {
	parsed, err := b.ParsePSBT(psbtBase64)
	if err != nil {
		return false, fmt.Errorf("failed to parse PSBT: %w", err)
	}

	return parsed.IsComplete, nil
}

// GetPSBTFee calculates the transaction fee from a PSBT.
// Used by Watch wallet to verify fee before broadcasting.
func (b *Bitcoin) GetPSBTFee(psbtBase64 string) (int64, error) {
	parsed, err := b.ParsePSBT(psbtBase64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse PSBT: %w", err)
	}

	// Calculate total input value
	var totalInput int64
	for _, input := range parsed.Packet.Inputs {
		if input.WitnessUtxo != nil {
			totalInput += input.WitnessUtxo.Value
		}
	}

	// Calculate total output value
	var totalOutput int64
	for _, output := range parsed.Packet.UnsignedTx.TxOut {
		totalOutput += output.Value
	}

	// Fee is the difference
	fee := totalInput - totalOutput
	if fee < 0 {
		return 0, errors.New("invalid PSBT: outputs exceed inputs (fee would be negative)")
	}

	return fee, nil
}

// serializePSBT serializes a PSBT packet to base64 string
func (*Bitcoin) serializePSBT(packet *psbt.Packet) (string, error) {
	var buf bytes.Buffer
	if err := packet.Serialize(&buf); err != nil {
		return "", fmt.Errorf("failed to serialize PSBT packet: %w", err)
	}

	psbtBase64 := base64.StdEncoding.EncodeToString(buf.Bytes())
	return psbtBase64, nil
}

// hasPartialSignatures checks if a PSBT has any partial signatures
func (*Bitcoin) hasPartialSignatures(packet *psbt.Packet) bool {
	for _, input := range packet.Inputs {
		if len(input.PartialSigs) > 0 {
			return true
		}
	}
	return false
}

// decodeHexScript decodes a hex-encoded script to bytes
func (*Bitcoin) decodeHexScript(hexScript string) ([]byte, error) {
	if hexScript == "" {
		return nil, errors.New("empty hex script")
	}
	// Remove "0x" prefix if present
	if len(hexScript) >= 2 && hexScript[:2] == "0x" {
		hexScript = hexScript[2:]
	}

	script := make([]byte, len(hexScript)/2)
	_, err := fmt.Sscanf(hexScript, "%x", &script)
	if err != nil {
		return nil, fmt.Errorf("failed to decode hex script: %w", err)
	}

	return script, nil
}
