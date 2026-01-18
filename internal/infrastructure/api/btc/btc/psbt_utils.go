// Package btc provides Bitcoin (BTC) infrastructure implementations.
//
// This file contains PSBT utility functions.
// See psbt.go for the main PSBT documentation.
package btc

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/btcsuite/btcd/btcutil/psbt"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"

	"github.com/hiromaily/go-crypto-wallet/pkg/logger"
)

// ExtractTransaction extracts the final signed transaction from a finalized PSBT.
// This should only be called after FinalizePSBT.
// Used by Watch wallet to get the transaction ready for broadcasting.
func (b *Bitcoin) ExtractTransaction(psbtBase64 string) (*wire.MsgTx, error) {
	// Parse PSBT
	parsed, err := b.parsePSBTInternal(psbtBase64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse PSBT for extraction: %w", err)
	}

	// DIAGNOSTIC: Log state before extraction
	for i, input := range parsed.Packet.Inputs {
		logger.Info("Input state before extraction",
			"input", i,
			"hasPartialSigs", len(input.PartialSigs) > 0,
			"partialSigsCount", len(input.PartialSigs),
			"hasFinalScriptSig", input.FinalScriptSig != nil,
			"finalScriptSigLen", len(input.FinalScriptSig),
			"hasRedeemScript", len(input.RedeemScript) > 0,
			"redeemScriptLen", len(input.RedeemScript))
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
//
// This function uses a custom completion check for multisig PSBTs because btcd's
// Packet.IsComplete() only checks if inputs are finalized, not if enough signatures
// exist for multisig threshold. For P2SH-P2WSH multisig (e.g., 2-of-3), we need to:
// 1. Parse the witness script to determine M-of-N requirements
// 2. Count partial signatures in each input
// 3. Compare signature count against required threshold
func (b *Bitcoin) IsPSBTComplete(psbtBase64 string) (bool, error) {
	parsed, err := b.parsePSBTInternal(psbtBase64)
	if err != nil {
		return false, fmt.Errorf("failed to parse PSBT: %w", err)
	}

	// Use custom multisig completion check
	return b.isMultisigPSBTComplete(parsed.Packet)
}

// GetPSBTFee calculates the transaction fee from a PSBT.
// Used by Watch wallet to verify fee before broadcasting.
func (b *Bitcoin) GetPSBTFee(psbtBase64 string) (int64, error) {
	parsed, err := b.parsePSBTInternal(psbtBase64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse PSBT: %w", err)
	}

	// Calculate total input value
	var totalInput int64
	for i, input := range parsed.Packet.Inputs {
		switch {
		case input.WitnessUtxo != nil:
			// SegWit/Taproot: Use WitnessUtxo
			totalInput += input.WitnessUtxo.Value
		case input.NonWitnessUtxo != nil:
			// Legacy: Use NonWitnessUtxo
			txIn := parsed.Packet.UnsignedTx.TxIn[i]
			if txIn.PreviousOutPoint.Index >= uint32(len(input.NonWitnessUtxo.TxOut)) {
				return 0, fmt.Errorf("input %d: previous output index %d out of range",
					i, txIn.PreviousOutPoint.Index)
			}
			totalInput += input.NonWitnessUtxo.TxOut[txIn.PreviousOutPoint.Index].Value
		default:
			return 0, fmt.Errorf("input %d: missing both WitnessUtxo and NonWitnessUtxo", i)
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

// WalletProcessPsbt signs a PSBT using Bitcoin Core's wallet (online RPC method).
// This is used when working with descriptor-based wallets where private keys are managed by Bitcoin Core.
// Used by Keygen wallet when descriptor-based workflow is enabled.
//
// Parameters:
//   - psbtBase64: Base64-encoded PSBT to sign
//   - sign: Whether to sign (true) or just add metadata (false)
//
// Returns:
//   - signedPSBT: The signed PSBT in base64 format
//   - isComplete: true if all signatures are collected
//   - error: any error that occurred
func (b *Bitcoin) WalletProcessPsbt(psbtBase64 string, sign bool) (string, bool, error) {
	// Call Bitcoin Core RPC
	signPtr := &sign
	// Use SigHashAll as default sighash type
	result, err := b.Client.WalletProcessPsbt(psbtBase64, signPtr, "ALL", nil)
	if err != nil {
		return "", false, fmt.Errorf("failed to call walletprocesspsbt RPC: %w", err)
	}

	logger.Debug("WalletProcessPsbt result",
		"complete", result.Complete,
		"psbt_len", len(result.Psbt))

	return result.Psbt, result.Complete, nil
}

// serializePSBT serializes a PSBT packet to base64 string.
func (b *Bitcoin) serializePSBT(packet *psbt.Packet) (string, error) {
	var buf bytes.Buffer
	if err := packet.Serialize(&buf); err != nil {
		return "", fmt.Errorf("failed to serialize PSBT: %w", err)
	}
	return base64.StdEncoding.EncodeToString(buf.Bytes()), nil
}

// hasPartialSignatures checks if any input in the PSBT has partial signatures.
func (b *Bitcoin) hasPartialSignatures(packet *psbt.Packet) bool {
	for _, input := range packet.Inputs {
		if len(input.PartialSigs) > 0 {
			return true
		}
	}
	return false
}

// isMultisigPSBTComplete checks if a multisig PSBT has enough signatures for each input.
// This is a custom implementation because btcd's Packet.IsComplete() only checks if inputs
// are finalized, not if there are enough partial signatures for the multisig threshold.
//
// Supports multiple multisig types:
// - P2WSH multisig: Uses WitnessScript
// - P2SH-P2WSH multisig: Uses both RedeemScript and WitnessScript
// - Legacy P2SH multisig: Uses RedeemScript
//
// For each input:
// 1. Parse the multisig script (WitnessScript or RedeemScript) to extract M-of-N requirements
// 2. Count partial signatures for the input
// 3. Return true if all inputs have at least M signatures
func (b *Bitcoin) isMultisigPSBTComplete(packet *psbt.Packet) (bool, error) {
	logger.Debug("Checking PSBT completion", "totalInputs", len(packet.Inputs))

	// Check each input
	for i, input := range packet.Inputs {
		logger.Debug("Checking input",
			"input", i,
			"hasFinalScriptSig", input.FinalScriptSig != nil,
			"hasFinalScriptWitness", input.FinalScriptWitness != nil,
			"witnessScriptLen", len(input.WitnessScript),
			"redeemScriptLen", len(input.RedeemScript),
			"partialSigsCount", len(input.PartialSigs))

		// If input is already finalized, it's complete
		if input.FinalScriptSig != nil || input.FinalScriptWitness != nil {
			logger.Debug("Input already finalized", "input", i)
			continue
		}

		// Determine which script to parse for multisig requirements
		// Priority: WitnessScript (P2WSH, P2SH-P2WSH) > RedeemScript (legacy P2SH)
		var scriptToParse []byte
		var scriptType string
		if len(input.WitnessScript) > 0 {
			scriptToParse = input.WitnessScript
			scriptType = "WitnessScript"
		} else if len(input.RedeemScript) > 0 {
			// Also check RedeemScript for legacy P2SH multisig
			scriptToParse = input.RedeemScript
			scriptType = "RedeemScript"
		}

		// If no script is available, assume single-sig and check for at least one signature
		if scriptToParse == nil {
			if len(input.PartialSigs) == 0 {
				logger.Debug("Input has no script and no partial signatures", "input", i)
				return false, nil
			}
			logger.Debug("Input has partial sigs but no multisig script (single-sig?)", "input", i)
			continue
		}

		// Log script for debugging
		logger.Debug("Input has multisig script",
			"input", i,
			"scriptType", scriptType,
			"scriptHex", hex.EncodeToString(scriptToParse))

		// Parse script to determine required signatures (M-of-N)
		requiredSigs, totalSigs, err := b.parseMultisigScript(scriptToParse)
		if err != nil {
			// If we can't parse it as multisig, fall back to checking for at least one partial sig
			// This handles non-multisig scripts
			logger.Debug("Failed to parse script as multisig, falling back to checking for any partial sig",
				"input", i, "scriptType", scriptType, "error", err, "partialSigs", len(input.PartialSigs))
			if len(input.PartialSigs) == 0 {
				return false, nil
			}
			continue
		}

		// Count partial signatures for this input
		sigCount := len(input.PartialSigs)

		logger.Debug("Checking multisig completion",
			"input", i,
			"requiredSigs", requiredSigs,
			"totalSigs", totalSigs,
			"currentSigs", sigCount)

		// Check if we have enough signatures
		if sigCount < requiredSigs {
			logger.Debug("Insufficient signatures for input",
				"input", i,
				"required", requiredSigs,
				"have", sigCount)
			return false, nil
		}
	}

	// All inputs have enough signatures
	logger.Debug("PSBT has enough signatures for all inputs")
	return true, nil
}

// parseMultisigScript parses a multisig witness script to extract M-of-N requirements.
// Returns (requiredSigs, totalSigs, error).
//
// Multisig script format: <M> <pubkey1> <pubkey2> ... <pubkeyN> <N> OP_CHECKMULTISIG
// Example 2-of-3: OP_2 <pk1> <pk2> <pk3> OP_3 OP_CHECKMULTISIG
func (*Bitcoin) parseMultisigScript(script []byte) (int, int, error) {
	if len(script) == 0 {
		return 0, 0, errors.New("empty script")
	}

	// Parse script to extract opcodes and data
	tokenizer := txscript.MakeScriptTokenizer(0, script)

	// First opcode should be OP_M (required signatures)
	if !tokenizer.Next() {
		return 0, 0, errors.New("failed to read first opcode")
	}

	firstOp := tokenizer.Opcode()
	if !txscript.IsSmallInt(firstOp) {
		return 0, 0, fmt.Errorf("first opcode is not a small int (M): %v", firstOp)
	}
	requiredSigs := txscript.AsSmallInt(firstOp)

	// Count public keys
	pubKeyCount := 0
	for tokenizer.Next() {
		opcode := tokenizer.Opcode()

		// Check if this is the total signers count (N)
		if txscript.IsSmallInt(opcode) {
			totalSigs := txscript.AsSmallInt(opcode)

			// Next should be OP_CHECKMULTISIG
			if !tokenizer.Next() {
				return 0, 0, errors.New("script ended before OP_CHECKMULTISIG")
			}
			if tokenizer.Opcode() != txscript.OP_CHECKMULTISIG {
				return 0, 0, fmt.Errorf("expected OP_CHECKMULTISIG, got %v", tokenizer.Opcode())
			}

			// Verify counts match
			if pubKeyCount != totalSigs {
				return 0, 0, fmt.Errorf("pubkey count mismatch: found %d pubkeys, script says %d",
					pubKeyCount, totalSigs)
			}

			logger.Debug("Parsed multisig script",
				"required", requiredSigs,
				"total", totalSigs,
				"pubkeys", pubKeyCount)

			return requiredSigs, totalSigs, nil
		}

		// This should be a public key
		if len(tokenizer.Data()) == 33 || len(tokenizer.Data()) == 65 {
			pubKeyCount++
		}
	}

	return 0, 0, errors.New("script does not match multisig format (missing N or OP_CHECKMULTISIG)")
}

// extractPubKeysFromScript extracts public keys from a multisig script in order
func (*Bitcoin) extractPubKeysFromScript(script []byte) ([][]byte, error) {
	if len(script) == 0 {
		return nil, errors.New("empty script")
	}

	// Parse script to extract public keys
	tokenizer := txscript.MakeScriptTokenizer(0, script)

	// First opcode should be OP_M (required signatures) - skip it
	if !tokenizer.Next() {
		return nil, errors.New("failed to read first opcode")
	}

	firstOp := tokenizer.Opcode()
	if !txscript.IsSmallInt(firstOp) {
		return nil, fmt.Errorf("first opcode is not a small int (M): %v", firstOp)
	}

	// Extract public keys in order
	pubKeys := make([][]byte, 0)
	for tokenizer.Next() {
		opcode := tokenizer.Opcode()

		// Check if this is the total signers count (N) - we're done
		if txscript.IsSmallInt(opcode) {
			// Next should be OP_CHECKMULTISIG
			if !tokenizer.Next() {
				return nil, errors.New("script ended before OP_CHECKMULTISIG")
			}
			if tokenizer.Opcode() != txscript.OP_CHECKMULTISIG {
				return nil, fmt.Errorf("expected OP_CHECKMULTISIG, got %v", tokenizer.Opcode())
			}

			logger.Debug("Extracted public keys from multisig script",
				"count", len(pubKeys))

			return pubKeys, nil
		}

		// This should be a public key (33 bytes compressed or 65 bytes uncompressed)
		data := tokenizer.Data()
		if len(data) == 33 || len(data) == 65 {
			// Make a copy of the data since tokenizer reuses the buffer
			pubKey := make([]byte, len(data))
			copy(pubKey, data)
			pubKeys = append(pubKeys, pubKey)
		}
	}

	return nil, errors.New("script does not match multisig format (missing N or OP_CHECKMULTISIG)")
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

	script, err := hex.DecodeString(hexScript)
	if err != nil {
		return nil, fmt.Errorf("failed to decode hex script: %w", err)
	}

	return script, nil
}
