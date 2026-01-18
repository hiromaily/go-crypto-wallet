// Package btc provides Bitcoin (BTC) infrastructure implementations.
//
// This file contains PSBT finalization functions.
// See psbt.go for the main PSBT documentation.
package btc

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcutil/psbt"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"

	"github.com/hiromaily/go-crypto-wallet/pkg/logger"
)

// FinalizePSBT finalizes a PSBT after all signatures are collected.
// This converts partial signatures into final scriptSig/scriptWitness.
// Used by Watch wallet after collecting all signatures from Keygen/Sign wallets.
func (b *Bitcoin) FinalizePSBT(psbtBase64 string) (string, error) {
	logger.Info("========== FinalizePSBT CALLED ==========")

	// Parse PSBT
	parsed, err := b.parsePSBTInternal(psbtBase64)
	if err != nil {
		logger.Error("Failed to parse PSBT in FinalizePSBT", "error", err)
		return "", fmt.Errorf("failed to parse PSBT for finalization: %w", err)
	}

	logger.Info("PSBT parsed for finalization",
		"inputs", len(parsed.Packet.Inputs),
		"outputs", len(parsed.Packet.Outputs))

	// Log input states before finalization
	for i, input := range parsed.Packet.Inputs {
		logger.Info("Input state before finalization",
			"input", i,
			"hasPartialSigs", len(input.PartialSigs) > 0,
			"partialSigsCount", len(input.PartialSigs),
			"hasRedeemScript", len(input.RedeemScript) > 0,
			"redeemScriptLen", len(input.RedeemScript),
			"hasWitnessScript", len(input.WitnessScript) > 0,
			"hasFinalScriptSig", input.FinalScriptSig != nil,
			"hasFinalScriptWitness", input.FinalScriptWitness != nil)
	}

	// Check if PSBT is complete using custom multisig completion check
	isComplete, err := b.isMultisigPSBTComplete(parsed.Packet)
	if err != nil {
		return "", fmt.Errorf("failed to check PSBT completion: %w", err)
	}
	logger.Info("PSBT completion check result", "isComplete", isComplete)
	if !isComplete {
		logger.Error("PSBT is incomplete - cannot finalize")
		return "", errors.New("cannot finalize incomplete PSBT (missing signatures)")
	}

	// Auto-generate missing RedeemScripts for P2SH-P2WSH inputs before finalization
	for i, input := range parsed.Packet.Inputs {
		logger.Debug("Checking RedeemScript auto-generation",
			"input", i,
			"hasRedeemScript", len(input.RedeemScript) > 0,
			"hasWitnessScript", len(input.WitnessScript) > 0,
			"hasWitnessUtxo", input.WitnessUtxo != nil)

		// If RedeemScript is missing but we have WitnessScript and witnessUtxo is P2SH
		if len(input.RedeemScript) == 0 && len(input.WitnessScript) > 0 && input.WitnessUtxo != nil {
			isP2SH := txscript.IsPayToScriptHash(input.WitnessUtxo.PkScript)
			logger.Debug("RedeemScript check details",
				"input", i,
				"isP2SH", isP2SH)

			if isP2SH {
				// Auto-generate RedeemScript for P2SH-P2WSH
				witnessScriptHash := sha256.Sum256(input.WitnessScript)
				redeemScript, err := txscript.NewScriptBuilder().
					AddOp(txscript.OP_0).
					AddData(witnessScriptHash[:]).
					Script()
				if err != nil {
					return "", fmt.Errorf("failed to build redeemScript for input %d: %w", i, err)
				}
				input.RedeemScript = redeemScript
				logger.Debug("Auto-generated RedeemScript for P2SH-P2WSH input during finalization", "input", i)
			}
		}
	}

	// Finalize all inputs
	for i := range parsed.Packet.UnsignedTx.TxIn {
		// Check if this is a P2SH-P2WSH multisig input
		input := parsed.Packet.Inputs[i]
		hasRedeemScript := len(input.RedeemScript) > 0
		hasWitnessScript := len(input.WitnessScript) > 0
		hasPartialSigs := len(input.PartialSigs) > 0

		// Detect script type for proper finalization
		var scriptType string
		if input.WitnessUtxo != nil {
			if txscript.IsPayToPubKeyHash(input.WitnessUtxo.PkScript) {
				scriptType = "P2PKH"
			} else if txscript.IsPayToScriptHash(input.WitnessUtxo.PkScript) {
				if hasWitnessScript {
					scriptType = "P2SH-P2WSH"
				} else {
					scriptType = "P2SH"
				}
			} else if txscript.IsPayToWitnessPubKeyHash(input.WitnessUtxo.PkScript) {
				scriptType = "P2WPKH"
			} else if txscript.IsPayToWitnessScriptHash(input.WitnessUtxo.PkScript) {
				scriptType = "P2WSH"
			} else if txscript.IsPayToTaproot(input.WitnessUtxo.PkScript) {
				scriptType = "P2TR"
			}
		}

		logger.Info("Checking finalization method for input",
			"input", i,
			"scriptType", scriptType,
			"hasRedeemScript", hasRedeemScript,
			"hasWitnessScript", hasWitnessScript,
			"hasPartialSigs", hasPartialSigs,
			"partialSigsCount", len(input.PartialSigs),
			"hasFinalScriptWitness", input.FinalScriptWitness != nil)

		// Skip already-finalized inputs (e.g., P2TR which sets FinalScriptWitness during signing)
		if input.FinalScriptWitness != nil || input.FinalScriptSig != nil {
			logger.Info("Input already finalized, skipping", "input", i, "scriptType", scriptType)
			continue
		}

		if scriptType == "P2PKH" && hasPartialSigs {
			// P2PKH - use custom finalization (creates scriptSig, no witness)
			logger.Info("Using custom P2PKH finalization", "input", i)
			if err := b.finalizeP2PKHInput(parsed.Packet, i); err != nil {
				return "", fmt.Errorf("failed to finalize P2PKH input %d: %w", i, err)
			}
		} else if scriptType == "P2SH" && hasRedeemScript && !hasWitnessScript && hasPartialSigs {
			// Check if redeemScript is P2WPKH (BIP49 Nested SegWit single-sig)
			// or multisig script (BIP44 P2SH multisig)
			isP2WPKHRedeemScript := txscript.IsPayToWitnessPubKeyHash(input.RedeemScript)

			if isP2WPKHRedeemScript {
				// P2SH-P2WPKH (BIP49 Nested SegWit single-sig) - use custom finalization
				logger.Info("Using custom P2SH-P2WPKH finalization", "input", i)
				if err := b.finalizeP2SHP2WPKHInput(parsed.Packet, i); err != nil {
					return "", fmt.Errorf("failed to finalize P2SH-P2WPKH input %d: %w", i, err)
				}
			} else {
				// P2SH multisig (non-SegWit, BIP44) - use custom finalization
				logger.Info("Using custom P2SH multisig finalization", "input", i)
				if err := b.finalizeMultisigInput(parsed.Packet, i); err != nil {
					return "", fmt.Errorf("failed to finalize P2SH multisig input %d: %w", i, err)
				}
			}
		} else if hasRedeemScript && hasWitnessScript && hasPartialSigs {
			// P2SH-P2WSH multisig - use custom finalization
			logger.Info("Using custom P2SH-P2WSH multisig finalization", "input", i)
			if err := b.finalizeMultisigInput(parsed.Packet, i); err != nil {
				return "", fmt.Errorf("failed to finalize P2SH-P2WSH multisig input %d: %w", i, err)
			}
		} else {
			// Other script types - use btcd's default finalization
			logger.Info("Using btcd default finalization", "input", i, "scriptType", scriptType,
				"reason", "no matching custom finalization condition")
			if err := psbt.Finalize(parsed.Packet, i); err != nil {
				logger.Error("btcd finalization failed", "input", i, "error", err)
				return "", fmt.Errorf("failed to finalize input %d: %w", i, err)
			}
		}
	}

	// Log input states after finalization
	for i, input := range parsed.Packet.Inputs {
		logger.Info("Input state after finalization",
			"input", i,
			"hasPartialSigs", len(input.PartialSigs) > 0,
			"hasFinalScriptSig", input.FinalScriptSig != nil,
			"finalScriptSigLen", len(input.FinalScriptSig),
			"hasFinalScriptWitness", input.FinalScriptWitness != nil)
	}

	// Serialize finalized PSBT to base64
	finalizedPSBT, err := b.serializePSBT(parsed.Packet)
	if err != nil {
		return "", fmt.Errorf("failed to serialize finalized PSBT: %w", err)
	}

	logger.Info("PSBT finalization completed",
		"inputs", len(parsed.Packet.UnsignedTx.TxIn))

	return finalizedPSBT, nil
}

// finalizeP2PKHInput finalizes a P2PKH (Pay-to-Public-Key-Hash) PSBT input.
// This creates a scriptSig in the format: <signature> <pubkey>
// P2PKH does NOT use witness data - all data goes in scriptSig.
//
//nolint:revive // receiver unused but method belongs to Bitcoin type
func (b *Bitcoin) finalizeP2PKHInput(packet *psbt.Packet, inputIndex int) error {
	// IMPORTANT: Get pointer to input, not a copy
	input := &packet.Inputs[inputIndex]

	// Ensure we have a partial signature
	if len(input.PartialSigs) == 0 {
		return errors.New("no signatures to finalize for P2PKH input")
	}

	// P2PKH single-sig should have exactly one signature
	if len(input.PartialSigs) != 1 {
		return fmt.Errorf("P2PKH input has %d signatures, expected 1", len(input.PartialSigs))
	}

	partialSig := input.PartialSigs[0]

	logger.Debug("Finalizing P2PKH input",
		"input", inputIndex,
		"sigLen", len(partialSig.Signature),
		"pubKeyLen", len(partialSig.PubKey))

	// Build scriptSig: <signature> <pubkey>
	scriptSig, err := txscript.NewScriptBuilder().
		AddData(partialSig.Signature).
		AddData(partialSig.PubKey).
		Script()
	if err != nil {
		return fmt.Errorf("failed to build P2PKH scriptSig: %w", err)
	}

	// Set the final scriptSig
	input.FinalScriptSig = scriptSig

	// Clear witness data (P2PKH doesn't use witness)
	input.FinalScriptWitness = nil

	// Clear partial sigs and other metadata (no longer needed after finalization)
	input.PartialSigs = nil
	input.SighashType = 0
	input.Bip32Derivation = nil

	logger.Debug("P2PKH input finalized",
		"input", inputIndex,
		"scriptSigLen", len(scriptSig),
		"hasWitness", input.FinalScriptWitness != nil)

	return nil
}

// serializeWitness serializes a TxWitness into the raw wire format.
// This is used to convert witness data into the format expected by PSBT FinalScriptWitness.
func serializeWitness(witness wire.TxWitness) ([]byte, error) {
	var witnessBuf bytes.Buffer
	if err := wire.WriteVarInt(&witnessBuf, 0, uint64(len(witness))); err != nil {
		return nil, fmt.Errorf("failed to write witness count: %w", err)
	}
	for _, elem := range witness {
		if err := wire.WriteVarBytes(&witnessBuf, 0, elem); err != nil {
			return nil, fmt.Errorf("failed to write witness element: %w", err)
		}
	}
	return witnessBuf.Bytes(), nil
}

// finalizeP2SHP2WPKHInput finalizes a P2SH-P2WPKH (BIP49 Nested SegWit single-sig) PSBT input.
// This creates a scriptSig containing the redeemScript (P2WPKH script) and witness data.
//
// For P2SH-P2WPKH:
//   - scriptSig contains: <redeemScript> (which is the P2WPKH script: OP_0 <20-byte-pubkey-hash>)
//   - witness contains: [<signature>, <pubkey>]
//
//nolint:revive // receiver unused but method belongs to Bitcoin type
func (b *Bitcoin) finalizeP2SHP2WPKHInput(packet *psbt.Packet, inputIndex int) error {
	// IMPORTANT: Get pointer to input, not a copy
	input := &packet.Inputs[inputIndex]

	// Ensure we have a partial signature
	if len(input.PartialSigs) == 0 {
		return errors.New("no signatures to finalize for P2SH-P2WPKH input")
	}

	// P2SH-P2WPKH single-sig should have exactly one signature
	if len(input.PartialSigs) != 1 {
		return fmt.Errorf("P2SH-P2WPKH input has %d signatures, expected 1", len(input.PartialSigs))
	}

	// Ensure we have a redeemScript (required for P2SH-P2WPKH)
	if len(input.RedeemScript) == 0 {
		return errors.New("missing redeem script for P2SH-P2WPKH input")
	}

	partialSig := input.PartialSigs[0]

	logger.Info("===== FINALIZING P2SH-P2WPKH INPUT =====",
		"input", inputIndex,
		"signature_hex", hex.EncodeToString(partialSig.Signature),
		"signature_len", len(partialSig.Signature),
		"pubkey_hex", hex.EncodeToString(partialSig.PubKey),
		"pubkey_len", len(partialSig.PubKey),
		"redeemScript_hex", hex.EncodeToString(input.RedeemScript),
		"redeemScript_len", len(input.RedeemScript))

	// For P2SH-P2WPKH, scriptSig contains the redeemScript
	// The redeemScript is the P2WPKH script: OP_0 <20-byte-pubkey-hash>
	// Use ScriptBuilder to properly serialize with length prefix
	scriptSigBuilder := txscript.NewScriptBuilder()
	scriptSigBuilder.AddData(input.RedeemScript)
	scriptSig, err := scriptSigBuilder.Script()
	if err != nil {
		return fmt.Errorf("failed to build P2SH-P2WPKH scriptSig: %w", err)
	}
	input.FinalScriptSig = scriptSig

	// Build witness: [<signature>, <pubkey>]
	witness := wire.TxWitness{
		partialSig.Signature,
		partialSig.PubKey,
	}

	// Serialize witness to FinalScriptWitness format
	witnessBytes, err := serializeWitness(witness)
	if err != nil {
		return fmt.Errorf("failed to serialize witness: %w", err)
	}
	input.FinalScriptWitness = witnessBytes

	// Clear partial sigs and other metadata (no longer needed after finalization)
	input.PartialSigs = nil
	input.SighashType = 0
	input.Bip32Derivation = nil
	input.RedeemScript = nil // Clear redeemScript as it's now in FinalScriptSig

	logger.Info("===== P2SH-P2WPKH INPUT FINALIZED =====",
		"input", inputIndex,
		"final_scriptSig_hex", hex.EncodeToString(input.FinalScriptSig),
		"final_scriptSig_len", len(input.FinalScriptSig),
		"final_witness_hex", hex.EncodeToString(input.FinalScriptWitness),
		"final_witness_len", len(input.FinalScriptWitness),
		"witness_count", len(witness))

	return nil
}

// finalizeMultisigInput finalizes a P2SH or P2SH-P2WSH multisig PSBT input.
// This is a custom implementation to work around btcd's Finalize() limitations with multisig.
//
// For P2SH multisig (non-SegWit, BIP44):
//   - scriptSig contains: [OP_0, sig1, sig2, ..., redeemScript]
//   - No witness data
//
// For P2SH-P2WSH multisig (SegWit):
//   - scriptSig contains the redeemScript (P2WSH: OP_0 <witnessScriptHash>)
//   - witness contains: [OP_0, sig1, sig2, ..., witnessScript]
func (b *Bitcoin) finalizeMultisigInput(packet *psbt.Packet, inputIndex int) error {
	input := &packet.Inputs[inputIndex]

	logger.Info("Starting finalizeMultisigInput",
		"input", inputIndex,
		"hasRedeemScript", len(input.RedeemScript) > 0,
		"redeemScriptLen", len(input.RedeemScript),
		"hasWitnessScript", len(input.WitnessScript) > 0,
		"partialSigsCount", len(input.PartialSigs))

	// Ensure we have the required scripts
	if len(input.RedeemScript) == 0 {
		logger.Error("Missing redeemScript in finalizeMultisigInput", "input", inputIndex)
		return errors.New("missing redeem script for P2SH multisig input")
	}
	if len(input.PartialSigs) == 0 {
		logger.Error("No partial signatures in finalizeMultisigInput", "input", inputIndex)
		return errors.New("no signatures to finalize")
	}

	// Determine if this is P2SH (non-SegWit) or P2SH-P2WSH (SegWit)
	isSegWit := len(input.WitnessScript) > 0
	logger.Info("Determined multisig type", "input", inputIndex, "isSegWit", isSegWit)

	// Extract public keys from the multisig script to determine signature order
	// For multisig, signatures MUST be ordered according to pubkey order in the script
	var scriptToExtract []byte
	if isSegWit {
		scriptToExtract = input.WitnessScript
	} else {
		scriptToExtract = input.RedeemScript
	}

	pubKeys, err := b.extractPubKeysFromScript(scriptToExtract)
	if err != nil {
		return fmt.Errorf("failed to extract public keys from multisig script: %w", err)
	}

	// Log redeemScript public keys
	logger.Debug("RedeemScript public keys",
		"input", inputIndex,
		"count", len(pubKeys))
	for i, pk := range pubKeys {
		logger.Debug("RedeemScript pubkey",
			"input", inputIndex,
			"index", i,
			"pubkey", hex.EncodeToString(pk),
			"len", len(pk))
	}

	// Log PartialSigs public keys
	logger.Debug("PartialSigs public keys",
		"input", inputIndex,
		"count", len(input.PartialSigs))
	for i, partialSig := range input.PartialSigs {
		logger.Debug("PartialSig pubkey",
			"input", inputIndex,
			"index", i,
			"pubkey", hex.EncodeToString(partialSig.PubKey),
			"len", len(partialSig.PubKey))
	}

	// Normalize public keys to compressed format for matching
	// PartialSigs always use compressed public keys (33 bytes), but redeem script
	// may contain uncompressed keys (65 bytes). Convert all to compressed for comparison.
	normalizedPubKeys := make([][]byte, len(pubKeys))
	for i, pk := range pubKeys {
		if len(pk) == 65 {
			// Uncompressed key - convert to compressed
			pubKeyObj, err := btcec.ParsePubKey(pk)
			if err != nil {
				return fmt.Errorf("failed to parse uncompressed public key at index %d: %w", i, err)
			}
			normalizedPubKeys[i] = pubKeyObj.SerializeCompressed()
			logger.Debug("Converted uncompressed pubkey to compressed",
				"input", inputIndex,
				"index", i,
				"uncompressed_len", len(pk),
				"compressed_len", len(normalizedPubKeys[i]))
		} else {
			// Already compressed (33 bytes) or unknown format
			normalizedPubKeys[i] = pk
		}
	}

	// Log public keys from multisig script for debugging
	logger.Debug("Public keys extracted from multisig script",
		"input", inputIndex,
		"count", len(pubKeys),
		"is_segwit", isSegWit)
	for i, pk := range normalizedPubKeys {
		logger.Debug("Multisig script pubkey (normalized)",
			"input", inputIndex,
			"index", i,
			"length", len(pk),
			"pubkey", hex.EncodeToString(pk))
	}

	// Log partial signatures for debugging
	logger.Debug("Partial signatures available",
		"input", inputIndex,
		"count", len(input.PartialSigs))
	for i, ps := range input.PartialSigs {
		logger.Debug("PartialSig pubkey",
			"input", inputIndex,
			"index", i,
			"length", len(ps.PubKey),
			"pubkey", hex.EncodeToString(ps.PubKey))
	}

	// Order signatures according to public key order in multisig script
	// PartialSigs is a slice []*PartialSig where each PartialSig has a PubKey field
	// Use normalizedPubKeys for matching since PartialSigs always use compressed format
	sigs := make([][]byte, 0, len(input.PartialSigs))
	for pkIdx, pubKey := range normalizedPubKeys {
		// Find the signature for this public key in PartialSigs
		found := false
		for psIdx, partialSig := range input.PartialSigs {
			if bytes.Equal(partialSig.PubKey, pubKey) {
				sigs = append(sigs, partialSig.Signature)
				logger.Debug("Matched signature for pubkey",
					"input", inputIndex,
					"multisigPubkeyIndex", pkIdx,
					"partialSigIndex", psIdx,
					"pubkey", hex.EncodeToString(pubKey))
				found = true
				break
			}
		}
		if !found {
			logger.Warn("No signature found for multisig script pubkey",
				"input", inputIndex,
				"multisigPubkeyIndex", pkIdx,
				"pubkey", hex.EncodeToString(pubKey))
		}
	}

	logger.Debug("Signature matching complete",
		"input", inputIndex,
		"multisigScriptPubkeys", len(normalizedPubKeys),
		"partialSigs", len(input.PartialSigs),
		"matchedSigs", len(sigs))

	if len(sigs) == 0 {
		return fmt.Errorf(
			"no matching signatures found: multisig script has %d pubkeys, PSBT has %d partial sigs, but none matched",
			len(normalizedPubKeys), len(input.PartialSigs),
		)
	}

	// Build final scripts based on whether this is SegWit or non-SegWit
	if isSegWit {
		// P2SH-P2WSH (SegWit) multisig:
		// Build witness stack: [OP_0, sig1, sig2, ..., witnessScript]
		witness := wire.TxWitness{
			[]byte{}, // OP_0 (required for CHECKMULTISIG bug)
		}
		for _, sig := range sigs {
			witness = append(witness, sig)
		}
		witness = append(witness, input.WitnessScript)

		// Serialize witness to FinalScriptWitness format
		witnessBytes, err := serializeWitness(witness)
		if err != nil {
			return fmt.Errorf("failed to serialize witness: %w", err)
		}
		input.FinalScriptWitness = witnessBytes

		// Build scriptSig for P2SH wrapping: just the redeemScript (OP_0 <witnessScriptHash>)
		scriptSigBuilder := txscript.NewScriptBuilder()
		scriptSigBuilder.AddData(input.RedeemScript)
		scriptSig, err := scriptSigBuilder.Script()
		if err != nil {
			return fmt.Errorf("failed to build P2SH-P2WSH scriptSig: %w", err)
		}
		input.FinalScriptSig = scriptSig

		logger.Debug("Finalized P2SH-P2WSH multisig input",
			"input", inputIndex,
			"signatures", len(sigs),
			"witnessLen", len(input.FinalScriptWitness),
			"scriptSigLen", len(input.FinalScriptSig))
	} else {
		// P2SH (non-SegWit, BIP44) multisig:
		// Build scriptSig: [OP_0, sig1, sig2, ..., redeemScript]
		scriptSigBuilder := txscript.NewScriptBuilder()
		scriptSigBuilder.AddOp(txscript.OP_0) // OP_0 (required for CHECKMULTISIG bug)
		for _, sig := range sigs {
			scriptSigBuilder.AddData(sig)
		}
		scriptSigBuilder.AddData(input.RedeemScript)
		scriptSig, err := scriptSigBuilder.Script()
		if err != nil {
			return fmt.Errorf("failed to build P2SH scriptSig: %w", err)
		}
		input.FinalScriptSig = scriptSig

		// No witness data for non-SegWit
		input.FinalScriptWitness = nil

		logger.Debug("Finalized P2SH multisig input",
			"input", inputIndex,
			"signatures", len(sigs),
			"scriptSigLen", len(input.FinalScriptSig))
	}

	// Clear partial signatures as they're no longer needed
	input.PartialSigs = nil

	return nil
}
