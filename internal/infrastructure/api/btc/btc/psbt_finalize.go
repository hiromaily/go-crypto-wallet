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

	// Log and check completion
	b.logInputStates(parsed.Packet, "before finalization")

	isComplete, err := b.isMultisigPSBTComplete(parsed.Packet)
	if err != nil {
		return "", fmt.Errorf("failed to check PSBT completion: %w", err)
	}
	logger.Info("PSBT completion check result", "isComplete", isComplete)
	if !isComplete {
		logger.Error("PSBT is incomplete - cannot finalize")
		return "", errors.New("cannot finalize incomplete PSBT (missing signatures)")
	}

	// Auto-generate missing RedeemScripts for P2SH-P2WSH inputs
	if err := b.autoGenerateMissingRedeemScripts(parsed.Packet); err != nil {
		return "", err
	}

	// Finalize all inputs
	if err := b.finalizeAllInputs(parsed.Packet); err != nil {
		return "", err
	}

	// Log and serialize
	b.logInputStates(parsed.Packet, "after finalization")

	finalizedPSBT, err := b.serializePSBT(parsed.Packet)
	if err != nil {
		return "", fmt.Errorf("failed to serialize finalized PSBT: %w", err)
	}

	logger.Info("PSBT finalization completed", "inputs", len(parsed.Packet.UnsignedTx.TxIn))
	return finalizedPSBT, nil
}

// logInputStates logs the state of all PSBT inputs.
func (*Bitcoin) logInputStates(packet *psbt.Packet, phase string) {
	for i, input := range packet.Inputs {
		logger.Info("Input state "+phase,
			"input", i,
			"hasPartialSigs", len(input.PartialSigs) > 0,
			"partialSigsCount", len(input.PartialSigs),
			"hasRedeemScript", len(input.RedeemScript) > 0,
			"redeemScriptLen", len(input.RedeemScript),
			"hasWitnessScript", len(input.WitnessScript) > 0,
			"hasFinalScriptSig", input.FinalScriptSig != nil,
			"hasFinalScriptWitness", input.FinalScriptWitness != nil)
	}
}

// autoGenerateMissingRedeemScripts generates RedeemScripts for P2SH-P2WSH inputs that are missing them.
func (*Bitcoin) autoGenerateMissingRedeemScripts(packet *psbt.Packet) error {
	for i, input := range packet.Inputs {
		logger.Debug("Checking RedeemScript auto-generation",
			"input", i,
			"hasRedeemScript", len(input.RedeemScript) > 0,
			"hasWitnessScript", len(input.WitnessScript) > 0,
			"hasWitnessUtxo", input.WitnessUtxo != nil)

		// If RedeemScript is missing but we have WitnessScript and witnessUtxo is P2SH
		if len(input.RedeemScript) == 0 && len(input.WitnessScript) > 0 && input.WitnessUtxo != nil {
			isP2SH := txscript.IsPayToScriptHash(input.WitnessUtxo.PkScript)
			logger.Debug("RedeemScript check details", "input", i, "isP2SH", isP2SH)

			if isP2SH {
				witnessScriptHash := sha256.Sum256(input.WitnessScript)
				redeemScript, err := txscript.NewScriptBuilder().
					AddOp(txscript.OP_0).
					AddData(witnessScriptHash[:]).
					Script()
				if err != nil {
					return fmt.Errorf("failed to build redeemScript for input %d: %w", i, err)
				}
				input.RedeemScript = redeemScript
				logger.Debug("Auto-generated RedeemScript for P2SH-P2WSH input during finalization", "input", i)
			}
		}
	}
	return nil
}

// finalizeAllInputs finalizes all inputs in the PSBT.
func (b *Bitcoin) finalizeAllInputs(packet *psbt.Packet) error {
	for i := range packet.UnsignedTx.TxIn {
		if err := b.finalizeInput(packet, i); err != nil {
			return err
		}
	}
	return nil
}

// finalizeInput finalizes a single PSBT input based on its script type.
func (b *Bitcoin) finalizeInput(packet *psbt.Packet, inputIndex int) error {
	input := packet.Inputs[inputIndex]

	// Skip already-finalized inputs
	if input.FinalScriptWitness != nil || input.FinalScriptSig != nil {
		scriptType := b.detectInputScriptType(&input)
		logger.Info("Input already finalized, skipping", "input", inputIndex, "scriptType", scriptType)
		return nil
	}

	scriptType := b.detectInputScriptType(&input)
	hasRedeemScript := len(input.RedeemScript) > 0
	hasWitnessScript := len(input.WitnessScript) > 0
	hasPartialSigs := len(input.PartialSigs) > 0

	logger.Info("Checking finalization method for input",
		"input", inputIndex,
		"scriptType", scriptType,
		"hasRedeemScript", hasRedeemScript,
		"hasWitnessScript", hasWitnessScript,
		"hasPartialSigs", hasPartialSigs,
		"partialSigsCount", len(input.PartialSigs),
		"hasFinalScriptWitness", input.FinalScriptWitness != nil)

	return b.dispatchFinalization(packet, inputIndex, scriptType, hasRedeemScript, hasWitnessScript, hasPartialSigs)
}

// detectInputScriptType detects the script type of a PSBT input.
func (*Bitcoin) detectInputScriptType(input *psbt.PInput) string {
	if input.WitnessUtxo == nil {
		return ""
	}

	hasWitnessScript := len(input.WitnessScript) > 0

	switch {
	case txscript.IsPayToPubKeyHash(input.WitnessUtxo.PkScript):
		return "P2PKH"
	case txscript.IsPayToScriptHash(input.WitnessUtxo.PkScript):
		if hasWitnessScript {
			return "P2SH-P2WSH"
		}
		return "P2SH"
	case txscript.IsPayToWitnessPubKeyHash(input.WitnessUtxo.PkScript):
		return "P2WPKH"
	case txscript.IsPayToWitnessScriptHash(input.WitnessUtxo.PkScript):
		return "P2WSH"
	case txscript.IsPayToTaproot(input.WitnessUtxo.PkScript):
		return "P2TR"
	default:
		return ""
	}
}

// dispatchFinalization dispatches to the appropriate finalization method based on script type.
func (b *Bitcoin) dispatchFinalization(
	packet *psbt.Packet,
	inputIndex int,
	scriptType string,
	hasRedeemScript, hasWitnessScript, hasPartialSigs bool,
) error {
	input := &packet.Inputs[inputIndex]

	switch {
	case scriptType == "P2PKH" && hasPartialSigs:
		logger.Info("Using custom P2PKH finalization", "input", inputIndex)
		if err := b.finalizeP2PKHInput(packet, inputIndex); err != nil {
			return fmt.Errorf("failed to finalize P2PKH input %d: %w", inputIndex, err)
		}

	case scriptType == "P2SH" && hasRedeemScript && !hasWitnessScript && hasPartialSigs:
		return b.finalizeP2SHInput(packet, inputIndex, input)

	case hasRedeemScript && hasWitnessScript && hasPartialSigs:
		logger.Info("Using custom P2SH-P2WSH multisig finalization", "input", inputIndex)
		if err := b.finalizeMultisigInput(packet, inputIndex); err != nil {
			return fmt.Errorf("failed to finalize P2SH-P2WSH multisig input %d: %w", inputIndex, err)
		}

	default:
		logger.Info("Using btcd default finalization", "input", inputIndex, "scriptType", scriptType,
			"reason", "no matching custom finalization condition")
		if err := psbt.Finalize(packet, inputIndex); err != nil {
			logger.Error("btcd finalization failed", "input", inputIndex, "error", err)
			return fmt.Errorf("failed to finalize input %d: %w", inputIndex, err)
		}
	}

	return nil
}

// finalizeP2SHInput finalizes a P2SH input (either P2SH-P2WPKH or P2SH multisig).
func (b *Bitcoin) finalizeP2SHInput(packet *psbt.Packet, inputIndex int, input *psbt.PInput) error {
	// Check if redeemScript is P2WPKH (BIP49 Nested SegWit single-sig)
	isP2WPKHRedeemScript := txscript.IsPayToWitnessPubKeyHash(input.RedeemScript)

	if isP2WPKHRedeemScript {
		logger.Info("Using custom P2SH-P2WPKH finalization", "input", inputIndex)
		if err := b.finalizeP2SHP2WPKHInput(packet, inputIndex); err != nil {
			return fmt.Errorf("failed to finalize P2SH-P2WPKH input %d: %w", inputIndex, err)
		}
	} else {
		logger.Info("Using custom P2SH multisig finalization", "input", inputIndex)
		if err := b.finalizeMultisigInput(packet, inputIndex); err != nil {
			return fmt.Errorf("failed to finalize P2SH multisig input %d: %w", inputIndex, err)
		}
	}

	return nil
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

	// Validate required data
	if err := validateMultisigInputData(input, inputIndex); err != nil {
		return err
	}

	// Determine script type and extract public keys
	isSegWit := len(input.WitnessScript) > 0
	logger.Info("Determined multisig type", "input", inputIndex, "isSegWit", isSegWit)

	scriptToExtract := input.RedeemScript
	if isSegWit {
		scriptToExtract = input.WitnessScript
	}

	pubKeys, err := b.extractPubKeysFromScript(scriptToExtract)
	if err != nil {
		return fmt.Errorf("failed to extract public keys from multisig script: %w", err)
	}

	// Log and normalize public keys
	logPubKeysDebugInfo(pubKeys, input.PartialSigs, inputIndex)
	normalizedPubKeys, err := normalizePubKeysToCompressed(pubKeys, inputIndex)
	if err != nil {
		return err
	}

	// Order signatures by public key order
	sigs := orderSignaturesByPubKeyOrder(normalizedPubKeys, input.PartialSigs, inputIndex)
	if len(sigs) == 0 {
		return fmt.Errorf(
			"no matching signatures found: multisig script has %d pubkeys, PSBT has %d partial sigs, but none matched",
			len(normalizedPubKeys), len(input.PartialSigs),
		)
	}

	// Build final scripts
	if isSegWit {
		if err := buildSegWitMultisigFinalScripts(input, sigs, inputIndex); err != nil {
			return err
		}
	} else {
		if err := buildLegacyMultisigFinalScripts(input, sigs, inputIndex); err != nil {
			return err
		}
	}

	// Clear partial signatures as they're no longer needed
	input.PartialSigs = nil
	return nil
}

// validateMultisigInputData validates that required data is present for multisig finalization.
func validateMultisigInputData(input *psbt.PInput, inputIndex int) error {
	if len(input.RedeemScript) == 0 {
		logger.Error("Missing redeemScript in finalizeMultisigInput", "input", inputIndex)
		return errors.New("missing redeem script for P2SH multisig input")
	}
	if len(input.PartialSigs) == 0 {
		logger.Error("No partial signatures in finalizeMultisigInput", "input", inputIndex)
		return errors.New("no signatures to finalize")
	}
	return nil
}

// logPubKeysDebugInfo logs debug information about public keys.
func logPubKeysDebugInfo(pubKeys [][]byte, partialSigs []*psbt.PartialSig, inputIndex int) {
	logger.Debug("RedeemScript public keys", "input", inputIndex, "count", len(pubKeys))
	for i, pk := range pubKeys {
		logger.Debug("RedeemScript pubkey",
			"input", inputIndex, "index", i,
			"pubkey", hex.EncodeToString(pk), "len", len(pk))
	}

	logger.Debug("PartialSigs public keys", "input", inputIndex, "count", len(partialSigs))
	for i, partialSig := range partialSigs {
		logger.Debug("PartialSig pubkey",
			"input", inputIndex, "index", i,
			"pubkey", hex.EncodeToString(partialSig.PubKey), "len", len(partialSig.PubKey))
	}
}

// normalizePubKeysToCompressed normalizes all public keys to compressed format.
func normalizePubKeysToCompressed(pubKeys [][]byte, inputIndex int) ([][]byte, error) {
	normalizedPubKeys := make([][]byte, len(pubKeys))
	for i, pk := range pubKeys {
		if len(pk) == 65 {
			pubKeyObj, err := btcec.ParsePubKey(pk)
			if err != nil {
				return nil, fmt.Errorf("failed to parse uncompressed public key at index %d: %w", i, err)
			}
			normalizedPubKeys[i] = pubKeyObj.SerializeCompressed()
			logger.Debug("Converted uncompressed pubkey to compressed",
				"input", inputIndex, "index", i,
				"uncompressed_len", len(pk), "compressed_len", len(normalizedPubKeys[i]))
		} else {
			normalizedPubKeys[i] = pk
		}
	}
	return normalizedPubKeys, nil
}

// orderSignaturesByPubKeyOrder orders signatures according to public key order in multisig script.
func orderSignaturesByPubKeyOrder(
	normalizedPubKeys [][]byte,
	partialSigs []*psbt.PartialSig,
	inputIndex int,
) [][]byte {
	sigs := make([][]byte, 0, len(partialSigs))
	for pkIdx, pubKey := range normalizedPubKeys {
		found := false
		for psIdx, partialSig := range partialSigs {
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
		"partialSigs", len(partialSigs),
		"matchedSigs", len(sigs))

	return sigs
}

// buildSegWitMultisigFinalScripts builds final scripts for P2SH-P2WSH (SegWit) multisig.
func buildSegWitMultisigFinalScripts(input *psbt.PInput, sigs [][]byte, inputIndex int) error {
	// Build witness stack: [OP_0, sig1, sig2, ..., witnessScript]
	witness := wire.TxWitness{[]byte{}} // OP_0 (required for CHECKMULTISIG bug)
	for _, sig := range sigs {
		witness = append(witness, sig)
	}
	witness = append(witness, input.WitnessScript)

	witnessBytes, err := serializeWitness(witness)
	if err != nil {
		return fmt.Errorf("failed to serialize witness: %w", err)
	}
	input.FinalScriptWitness = witnessBytes

	// Build scriptSig for P2SH wrapping
	scriptSigBuilder := txscript.NewScriptBuilder()
	scriptSigBuilder.AddData(input.RedeemScript)
	scriptSig, err := scriptSigBuilder.Script()
	if err != nil {
		return fmt.Errorf("failed to build P2SH-P2WSH scriptSig: %w", err)
	}
	input.FinalScriptSig = scriptSig

	logger.Debug("Finalized P2SH-P2WSH multisig input",
		"input", inputIndex, "signatures", len(sigs),
		"witnessLen", len(input.FinalScriptWitness), "scriptSigLen", len(input.FinalScriptSig))
	return nil
}

// buildLegacyMultisigFinalScripts builds final scripts for P2SH (non-SegWit) multisig.
func buildLegacyMultisigFinalScripts(input *psbt.PInput, sigs [][]byte, inputIndex int) error {
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
	input.FinalScriptWitness = nil // No witness data for non-SegWit

	logger.Debug("Finalized P2SH multisig input",
		"input", inputIndex, "signatures", len(sigs), "scriptSigLen", len(input.FinalScriptSig))
	return nil
}
