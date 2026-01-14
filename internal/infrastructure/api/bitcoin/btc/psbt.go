package btc

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/btcsuite/btcd/btcec/v2/ecdsa"
	"github.com/btcsuite/btcd/btcec/v2/schnorr"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/btcutil/psbt"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"

	dtobtc "github.com/hiromaily/go-crypto-wallet/internal/application/dto/btc"
	domainAccount "github.com/hiromaily/go-crypto-wallet/internal/domain/account"
	"github.com/hiromaily/go-crypto-wallet/pkg/logger"
)

// ParsedPSBT represents a parsed PSBT with metadata
type ParsedPSBT struct {
	Packet       *psbt.Packet
	InputCount   int
	OutputCount  int
	IsComplete   bool
	HasSignature bool
}

// CreatePSBT creates a PSBT from an unsigned transaction with metadata.
// This function adds all necessary metadata (witness UTXO, redeem scripts, BIP32 derivation) for offline signing.
// Used by Watch wallet to create unsigned PSBTs.
//
//nolint:gocyclo // Complex function handling PSBT creation with metadata
func (b *Bitcoin) CreatePSBT(msgTx *wire.MsgTx, prevTxs []dtobtc.PreviousTx, senderAccount domainAccount.AccountType) (string, error) {
	logger.Info("Creating PSBT",
		"inputs", len(msgTx.TxIn),
		"outputs", len(msgTx.TxOut),
		"prevTxs", len(prevTxs),
		"sender_account", senderAccount.String())

	// Convert application DTOs to infrastructure types
	infraPrevTxs, err := FromPreviousTx(prevTxs, b)
	if err != nil {
		return "", fmt.Errorf("failed to convert PreviousTx: %w", err)
	}
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
	for i, prevTx := range infraPrevTxs {
		if i >= len(packet.UnsignedTx.TxIn) {
			return "", fmt.Errorf("prevTxs index %d exceeds number of inputs %d", i, len(packet.UnsignedTx.TxIn))
		}

		// Validate that prevTx matches the transaction input
		txIn := packet.UnsignedTx.TxIn[i]

		// Parse and validate txid
		prevTxHash, err := chainhash.NewHashFromStr(prevTx.Txid)
		if err != nil {
			return "", fmt.Errorf("failed to parse txid for input %d: %w", i, err)
		}
		if !prevTxHash.IsEqual(&txIn.PreviousOutPoint.Hash) {
			return "", fmt.Errorf("input %d: prevTx txid %s does not match transaction input %s",
				i, prevTxHash.String(), txIn.PreviousOutPoint.Hash.String())
		}

		// Validate vout index
		if prevTx.Vout != txIn.PreviousOutPoint.Index {
			return "", fmt.Errorf("input %d: prevTx vout %d does not match transaction input vout %d",
				i, prevTx.Vout, txIn.PreviousOutPoint.Index)
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

		// Add witness script for P2WSH (native SegWit multisig) if provided
		// This is the actual multisig script that defines required signatures
		var witnessScript []byte
		if prevTx.WitnessScript != "" {
			var err error
			witnessScript, err = b.decodeHexScript(prevTx.WitnessScript)
			if err != nil {
				return "", fmt.Errorf("failed to decode witnessScript for input %d: %w", i, err)
			}
			if err := updater.AddInWitnessScript(witnessScript, i); err != nil {
				return "", fmt.Errorf("failed to add witness script for input %d: %w", i, err)
			}
		}

		// Add redeem script for P2SH if provided
		// For P2SH-wrapped SegWit, this is the witness program
		if prevTx.RedeemScript != "" {
			redeemScript, err := b.decodeHexScript(prevTx.RedeemScript)
			if err != nil {
				return "", fmt.Errorf("failed to decode redeemScript for input %d: %w", i, err)
			}
			if err := updater.AddInRedeemScript(redeemScript, i); err != nil {
				return "", fmt.Errorf("failed to add redeem script for input %d: %w", i, err)
			}
		} else if len(witnessScript) > 0 && txscript.IsPayToScriptHash(scriptPubKey) {
			// Auto-generate RedeemScript for P2SH-P2WSH when it's missing
			// RedeemScript is the witness program: OP_0 <witnessScriptHash>
			witnessScriptHash := sha256.Sum256(witnessScript)
			redeemScript, err := txscript.NewScriptBuilder().
				AddOp(txscript.OP_0).
				AddData(witnessScriptHash[:]).
				Script()
			if err != nil {
				return "", fmt.Errorf("failed to build redeemScript for input %d: %w", i, err)
			}
			if err := updater.AddInRedeemScript(redeemScript, i); err != nil {
				return "", fmt.Errorf("failed to add auto-generated redeem script for input %d: %w", i, err)
			}
			logger.Debug("Auto-generated RedeemScript for P2SH-P2WSH input", "input", i)
		}

		// Add sighash type (default to SIGHASH_ALL)
		if err := updater.AddInSighashType(txscript.SigHashAll, i); err != nil {
			return "", fmt.Errorf("failed to add sighash type for input %d: %w", i, err)
		}

		// Add BIP32 derivation path for descriptor-based signing
		// This is REQUIRED for descriptor wallets to sign PSBTs
		if err := b.addBIP32DerivationForInput(updater, scriptPubKey, i, senderAccount); err != nil {
			logger.Error("failed to add BIP32 derivation for input",
				"input_index", i,
				"error", err,
				"sender_account", senderAccount.String())
			return "", fmt.Errorf("failed to add BIP32 derivation for input %d (required for descriptor-based signing): %w", i, err)
		}
		logger.Debug("Added BIP32 derivation for input",
			"input_index", i,
			"sender_account", senderAccount.String())
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

// parsePSBTInternal parses a base64-encoded PSBT and returns infrastructure type.
// Internal helper for methods that need access to the underlying packet.
func (b *Bitcoin) parsePSBTInternal(psbtBase64 string) (*ParsedPSBT, error) {
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

	infraParsed := &ParsedPSBT{
		Packet:       packet,
		InputCount:   len(packet.Inputs),
		OutputCount:  len(packet.Outputs),
		IsComplete:   packet.IsComplete(),
		HasSignature: hasSignature,
	}

	logger.Debug("Parsed PSBT",
		"inputs", infraParsed.InputCount,
		"outputs", infraParsed.OutputCount,
		"complete", infraParsed.IsComplete,
		"hasSignature", infraParsed.HasSignature)

	return infraParsed, nil
}

// ParsePSBT parses a base64-encoded PSBT and returns metadata as application DTO.
// Used by all wallets to read PSBT files.
func (b *Bitcoin) ParsePSBT(psbtBase64 string) (*dtobtc.ParsedPSBT, error) {
	infraParsed, err := b.parsePSBTInternal(psbtBase64)
	if err != nil {
		return nil, err
	}

	return ToParsedPSBT(infraParsed, b)
}

// ValidatePSBT validates a PSBT structure and checks BIP174 compliance.
// Used by all wallets to verify PSBT before processing.
func (b *Bitcoin) ValidatePSBT(psbtBase64 string) error {
	parsed, err := b.parsePSBTInternal(psbtBase64)
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

	// Validate each input has either WitnessUtxo (SegWit/Taproot) or NonWitnessUtxo (Legacy)
	for i, input := range parsed.Packet.Inputs {
		if input.WitnessUtxo == nil && input.NonWitnessUtxo == nil {
			return fmt.Errorf("input %d missing both WitnessUtxo and NonWitnessUtxo", i)
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
	parsed, err := b.parsePSBTInternal(psbtBase64)
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

	// Create PrevOutputFetcher for Taproot signing
	prevOutputFetcher := txscript.NewMultiPrevOutFetcher(nil)
	for i, input := range parsed.Packet.Inputs {
		if input.WitnessUtxo != nil {
			prevOut := parsed.Packet.UnsignedTx.TxIn[i].PreviousOutPoint
			prevOutputFetcher.AddPrevOut(prevOut, input.WitnessUtxo)
		}
	}

	// Sign each input with each provided key
	signedCount := 0
	for i := range parsed.Packet.UnsignedTx.TxIn {
		// Get PSBT input which contains metadata (WitnessUtxo, RedeemScript, WitnessScript)
		psbtInput := &parsed.Packet.Inputs[i]

		logger.Debug("Processing PSBT input for signing",
			"input", i,
			"hasWitnessUtxo", psbtInput.WitnessUtxo != nil,
			"hasNonWitnessUtxo", psbtInput.NonWitnessUtxo != nil,
			"hasBIP32Derivation", len(psbtInput.Bip32Derivation) > 0,
			"hasRedeemScript", len(psbtInput.RedeemScript) > 0,
			"hasWitnessScript", len(psbtInput.WitnessScript) > 0)

		if psbtInput.WitnessUtxo == nil {
			logger.Warn("Skipping input without witness UTXO", "input", i)
			continue
		}

		// Detect script type for logging purposes.
		// NOTE: If this pattern is needed in multiple places, consider extracting to
		// pkg/cryptocurrency/btc_script.go as DetectScriptType(pkScript []byte) string.
		scriptType := "unknown"
		if txscript.IsPayToTaproot(psbtInput.WitnessUtxo.PkScript) {
			scriptType = "P2TR (Taproot)"
		} else if txscript.IsPayToWitnessPubKeyHash(psbtInput.WitnessUtxo.PkScript) {
			scriptType = "P2WPKH (SegWit v0)"
		} else if txscript.IsPayToWitnessScriptHash(psbtInput.WitnessUtxo.PkScript) {
			scriptType = "P2WSH (SegWit v0)"
		} else if txscript.IsPayToScriptHash(psbtInput.WitnessUtxo.PkScript) {
			scriptType = "P2SH"
		} else if txscript.IsPayToPubKeyHash(psbtInput.WitnessUtxo.PkScript) {
			scriptType = "P2PKH (Legacy)"
		}

		logger.Debug("PSBT input script type detected",
			"input", i,
			"scriptType", scriptType,
			"pkScriptHex", hex.EncodeToString(psbtInput.WitnessUtxo.PkScript))

		// Try signing with each private key
		for j, privKey := range privKeys {
			logger.Debug("Attempting to sign with key",
				"input", i,
				"keyIndex", j,
				"pubKey", hex.EncodeToString(privKey.PrivKey.PubKey().SerializeCompressed()))

			if b.signInputWithKey(updater, parsed.Packet.UnsignedTx, i, psbtInput, privKey, prevOutputFetcher) {
				signedCount++
				logger.Info("Successfully signed input",
					"input", i,
					"keyIndex", j)
			}
		}
	}

	if signedCount == 0 {
		return "", false, errors.New("no signatures were added (keys may not match PSBT inputs)")
	}

	// Check if PSBT is now complete using custom multisig completion check
	isComplete, err := b.isMultisigPSBTComplete(parsed.Packet)
	if err != nil {
		logger.Warn("Failed to check PSBT completion, assuming incomplete", "error", err)
		isComplete = false
	}

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
	logger.Info("FinalizePSBT called")

	// Parse PSBT
	parsed, err := b.parsePSBTInternal(psbtBase64)
	if err != nil {
		return "", fmt.Errorf("failed to parse PSBT for finalization: %w", err)
	}

	logger.Debug("PSBT parsed for finalization",
		"inputs", len(parsed.Packet.Inputs),
		"outputs", len(parsed.Packet.Outputs))

	// Log input states before finalization
	for i, input := range parsed.Packet.Inputs {
		logger.Debug("Input state before finalization",
			"input", i,
			"hasPartialSigs", len(input.PartialSigs) > 0,
			"partialSigsCount", len(input.PartialSigs),
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

		logger.Debug("Checking finalization method for input",
			"input", i,
			"scriptType", scriptType,
			"hasRedeemScript", hasRedeemScript,
			"hasWitnessScript", hasWitnessScript,
			"hasPartialSigs", hasPartialSigs,
			"partialSigsCount", len(input.PartialSigs))

		if scriptType == "P2PKH" && hasPartialSigs {
			// P2PKH - use custom finalization (creates scriptSig, no witness)
			logger.Debug("Using custom P2PKH finalization", "input", i)
			if err := b.finalizeP2PKHInput(parsed.Packet, i); err != nil {
				return "", fmt.Errorf("failed to finalize P2PKH input %d: %w", i, err)
			}
		} else if hasRedeemScript && hasWitnessScript && hasPartialSigs {
			// P2SH-P2WSH multisig - use custom finalization
			logger.Debug("Using custom multisig finalization", "input", i)
			if err := b.finalizeMultisigInput(parsed.Packet, i); err != nil {
				return "", fmt.Errorf("failed to finalize multisig input %d: %w", i, err)
			}
		} else {
			// Other script types - use btcd's default finalization
			logger.Debug("Using btcd default finalization", "input", i, "scriptType", scriptType)
			if err := psbt.Finalize(parsed.Packet, i); err != nil {
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

// finalizeMultisigInput finalizes a P2SH-P2WSH multisig PSBT input.
// This is a custom implementation to work around btcd's Finalize() limitations with P2SH-P2WSH.
//
// For P2SH-P2WSH multisig:
//   - scriptSig contains the redeemScript (P2WSH: OP_0 <witnessScriptHash>)
//   - witness contains: [OP_0, sig1, sig2, ..., witnessScript]
func (b *Bitcoin) finalizeMultisigInput(packet *psbt.Packet, inputIndex int) error {
	input := packet.Inputs[inputIndex]

	// Ensure we have the required scripts
	if len(input.RedeemScript) == 0 {
		return errors.New("missing redeem script for P2SH-P2WSH input")
	}
	if len(input.WitnessScript) == 0 {
		return errors.New("missing witness script for P2WSH multisig input")
	}
	if len(input.PartialSigs) == 0 {
		return errors.New("no signatures to finalize")
	}

	// Extract public keys from witness script to determine signature order
	// For multisig, signatures MUST be ordered according to pubkey order in witness script
	pubKeys, err := b.extractPubKeysFromScript(input.WitnessScript)
	if err != nil {
		return fmt.Errorf("failed to extract public keys from witness script: %w", err)
	}

	// Log public keys from witness script for debugging
	logger.Debug("Public keys extracted from witness script",
		"input", inputIndex,
		"count", len(pubKeys))
	for i, pk := range pubKeys {
		logger.Debug("Witness script pubkey",
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

	// Order signatures according to public key order in witness script
	// PartialSigs is a slice []*PartialSig where each PartialSig has a PubKey field
	sigs := make([][]byte, 0, len(input.PartialSigs))
	for pkIdx, pubKey := range pubKeys {
		// Find the signature for this public key in PartialSigs
		found := false
		for psIdx, partialSig := range input.PartialSigs {
			if bytes.Equal(partialSig.PubKey, pubKey) {
				sigs = append(sigs, partialSig.Signature)
				logger.Debug("Matched signature for pubkey",
					"input", inputIndex,
					"witnessPubkeyIndex", pkIdx,
					"partialSigIndex", psIdx,
					"pubkey", hex.EncodeToString(pubKey))
				found = true
				break
			}
		}
		if !found {
			logger.Warn("No signature found for witness script pubkey",
				"input", inputIndex,
				"witnessPubkeyIndex", pkIdx,
				"pubkey", hex.EncodeToString(pubKey))
		}
	}

	logger.Debug("Signature matching complete",
		"input", inputIndex,
		"witnessScriptPubkeys", len(pubKeys),
		"partialSigs", len(input.PartialSigs),
		"matchedSigs", len(sigs))

	if len(sigs) == 0 {
		return fmt.Errorf("no matching signatures found: witness script has %d pubkeys, PSBT has %d partial sigs, but none matched",
			len(pubKeys), len(input.PartialSigs))
	}

	// Build witness stack for P2WSH multisig:
	// [OP_0, sig1, sig2, ..., witnessScript]
	witness := wire.TxWitness{
		[]byte{}, // OP_0 (required for CHECKMULTISIG bug)
	}
	for _, sig := range sigs {
		witness = append(witness, sig)
	}
	witness = append(witness, input.WitnessScript)

	// Serialize witness to FinalScriptWitness format
	// PSBT witness format is the same as transaction witness serialization
	var witnessBuf bytes.Buffer
	// Write witness stack count (always 1 for witness serialization per input)
	if err := wire.WriteVarInt(&witnessBuf, 0, uint64(len(witness))); err != nil {
		return fmt.Errorf("failed to write witness count: %w", err)
	}
	// Write each witness element
	for _, elem := range witness {
		if err := wire.WriteVarBytes(&witnessBuf, 0, elem); err != nil {
			return fmt.Errorf("failed to write witness element: %w", err)
		}
	}
	input.FinalScriptWitness = witnessBuf.Bytes()

	// Build scriptSig for P2SH wrapping:
	// For P2SH-P2WSH, scriptSig is just the redeemScript (P2WSH: OP_0 <hash>)
	scriptSigBuilder := txscript.NewScriptBuilder()
	scriptSigBuilder.AddData(input.RedeemScript)
	scriptSig, err := scriptSigBuilder.Script()
	if err != nil {
		return fmt.Errorf("failed to build scriptSig: %w", err)
	}
	input.FinalScriptSig = scriptSig

	// Clear partial signatures as they're no longer needed
	input.PartialSigs = nil

	logger.Debug("Finalized P2SH-P2WSH multisig input",
		"input", inputIndex,
		"signatures", len(sigs),
		"witnessLen", len(input.FinalScriptWitness),
		"scriptSigLen", len(input.FinalScriptSig))

	return nil
}

// ExtractTransaction extracts the final signed transaction from a finalized PSBT.
// This should only be called after FinalizePSBT.
// Used by Watch wallet to get the transaction ready for broadcasting.
func (b *Bitcoin) ExtractTransaction(psbtBase64 string) (*wire.MsgTx, error) {
	// Parse PSBT
	parsed, err := b.parsePSBTInternal(psbtBase64)
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

	logger.Debug("walletprocesspsbt completed",
		"complete", result.Complete,
		"sign", sign)

	return result.Psbt, result.Complete, nil
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

// signInputWithKey signs a single PSBT input with a private key.
// Returns true if signature was successfully added, false otherwise.
//
// This function acts as a dispatcher, routing to the appropriate signing method
// based on script type:
//   - Taproot (P2TR): signTaprootInput
//   - SegWit v0 (P2WPKH/P2WSH): signSegWitInput
//   - Legacy (P2PKH/P2SH): signLegacyInput
func (b *Bitcoin) signInputWithKey(
	updater *psbt.Updater,
	msgTx *wire.MsgTx,
	inputIndex int,
	psbtInput *psbt.PInput,
	privKey *btcutil.WIF,
	prevOutputFetcher *txscript.MultiPrevOutFetcher,
) bool {
	witnessUtxo := psbtInput.WitnessUtxo

	// Detect script type to determine signing method
	isTaproot := txscript.IsPayToTaproot(witnessUtxo.PkScript)
	isSegWit := txscript.IsPayToWitnessPubKeyHash(witnessUtxo.PkScript) ||
		txscript.IsPayToWitnessScriptHash(witnessUtxo.PkScript) ||
		(txscript.IsPayToScriptHash(witnessUtxo.PkScript) && len(psbtInput.WitnessScript) > 0)

	// Route to appropriate signing function
	if isTaproot {
		return b.signTaprootInput(updater, msgTx, inputIndex, psbtInput, privKey, prevOutputFetcher)
	} else if isSegWit {
		return b.signSegWitInput(updater, msgTx, inputIndex, psbtInput, privKey, prevOutputFetcher)
	}
	return b.signLegacyInput(msgTx, inputIndex, psbtInput, privKey)
}

// signTaprootInput signs a Taproot (P2TR) input using Schnorr signatures.
func (*Bitcoin) signTaprootInput(
	updater *psbt.Updater,
	msgTx *wire.MsgTx,
	inputIndex int,
	psbtInput *psbt.PInput,
	privKey *btcutil.WIF,
	prevOutputFetcher *txscript.MultiPrevOutFetcher,
) bool {
	// Calculate Taproot signature hash
	sigHashes := txscript.NewTxSigHashes(msgTx, prevOutputFetcher)
	hash, err := txscript.CalcTaprootSignatureHash(
		sigHashes,
		txscript.SigHashDefault,
		msgTx,
		inputIndex,
		prevOutputFetcher,
	)
	if err != nil {
		logger.Warn("Failed to calculate Taproot signature hash", "input", inputIndex, "error", err)
		return false
	}

	// Sign with Schnorr
	signature, err := schnorr.Sign(privKey.PrivKey, hash)
	if err != nil {
		logger.Warn("Failed to create Schnorr signature", "input", inputIndex, "error", err)
		return false
	}

	// For Taproot, signature is just the Schnorr signature (no sighash type appended for SIGHASH_DEFAULT)
	sigBytes := signature.Serialize()
	pubKey := privKey.PrivKey.PubKey().SerializeCompressed()

	// Add signature using btcd's updater
	outcome, err := updater.Sign(inputIndex, sigBytes, pubKey, psbtInput.RedeemScript, psbtInput.WitnessScript)
	if err != nil {
		logger.Debug("Signature not applicable for this input", "input", inputIndex, "error", err)
		return false
	}

	if outcome == psbt.SignSuccesful {
		logger.Debug("Added Taproot signature to input", "input", inputIndex)
		return true
	}

	logger.Debug("Taproot signature not added", "input", inputIndex, "outcome", outcome)
	return false
}

// signSegWitInput signs a SegWit v0 (P2WPKH/P2WSH) input using ECDSA signatures.
func (*Bitcoin) signSegWitInput(
	updater *psbt.Updater,
	msgTx *wire.MsgTx,
	inputIndex int,
	psbtInput *psbt.PInput,
	privKey *btcutil.WIF,
	prevOutputFetcher *txscript.MultiPrevOutFetcher,
) bool {
	witnessUtxo := psbtInput.WitnessUtxo

	// Calculate witness signature hash
	sigHashes := txscript.NewTxSigHashes(msgTx, prevOutputFetcher)

	// For P2WSH multisig, use the witness script for hash calculation
	// For P2WPKH (single sig), use the scriptPubKey
	scriptForHash := witnessUtxo.PkScript
	if len(psbtInput.WitnessScript) > 0 {
		scriptForHash = psbtInput.WitnessScript
		logger.Debug("Using witness script for P2WSH multisig hash calculation",
			"input", inputIndex,
			"witnessScriptLen", len(psbtInput.WitnessScript))
	}

	hash, err := txscript.CalcWitnessSigHash(
		scriptForHash,
		sigHashes,
		txscript.SigHashAll,
		msgTx,
		inputIndex,
		witnessUtxo.Value,
	)
	if err != nil {
		logger.Warn("Failed to calculate witness signature hash", "input", inputIndex, "error", err)
		return false
	}

	// Sign with ECDSA
	signature := ecdsa.Sign(privKey.PrivKey, hash)
	sigBytes := append(signature.Serialize(), byte(txscript.SigHashAll))
	pubKey := privKey.PrivKey.PubKey().SerializeCompressed()

	// Add signature using btcd's updater
	outcome, err := updater.Sign(inputIndex, sigBytes, pubKey, psbtInput.RedeemScript, psbtInput.WitnessScript)
	if err != nil {
		logger.Debug("Signature not applicable for this input", "input", inputIndex, "error", err)
		return false
	}

	if outcome == psbt.SignSuccesful {
		logger.Debug("Added SegWit signature to input",
			"input", inputIndex,
			"hasWitnessScript", len(psbtInput.WitnessScript) > 0)
		return true
	}

	logger.Debug("SegWit signature not added", "input", inputIndex, "outcome", outcome)
	return false
}

// signLegacyInput signs a legacy (P2PKH/P2SH) input using ECDSA signatures.
// Note: btcd's updater.Sign() has limitations with P2PKH, so we add the signature directly to PartialSigs.
func (*Bitcoin) signLegacyInput(
	msgTx *wire.MsgTx,
	inputIndex int,
	psbtInput *psbt.PInput,
	privKey *btcutil.WIF,
) bool {
	witnessUtxo := psbtInput.WitnessUtxo

	logger.Debug("Using legacy signature algorithm for P2PKH/P2SH",
		"input", inputIndex,
		"pkScript", hex.EncodeToString(witnessUtxo.PkScript))

	// For legacy P2PKH, the scriptPubKey is used for signature hash
	// For P2SH, use the redeem script if available
	scriptForHash := witnessUtxo.PkScript
	if len(psbtInput.RedeemScript) > 0 {
		scriptForHash = psbtInput.RedeemScript
		logger.Debug("Using redeem script for legacy P2SH signature",
			"input", inputIndex,
			"redeemScriptLen", len(psbtInput.RedeemScript))
	}

	// Calculate legacy signature hash
	hash, err := txscript.CalcSignatureHash(
		scriptForHash,
		txscript.SigHashAll,
		msgTx,
		inputIndex,
	)
	if err != nil {
		logger.Warn("Failed to calculate legacy signature hash", "input", inputIndex, "error", err)
		return false
	}

	logger.Debug("Legacy signature hash calculated",
		"input", inputIndex,
		"hash", hex.EncodeToString(hash),
		"scriptForHashLen", len(scriptForHash))

	// Sign with ECDSA
	signature := ecdsa.Sign(privKey.PrivKey, hash)

	logger.Debug("ECDSA signature created",
		"input", inputIndex,
		"sigR", signature.R().String(),
		"sigS", signature.S().String())

	// Verify signature locally before adding to PSBT
	pubKeyObj := privKey.PrivKey.PubKey()
	if !signature.Verify(hash, pubKeyObj) {
		logger.Error("Signature verification FAILED locally",
			"input", inputIndex,
			"pubKey", hex.EncodeToString(pubKeyObj.SerializeCompressed()))
		return false
	}
	logger.Debug("Signature verified successfully locally", "input", inputIndex)

	// Serialize signature with sighash type
	sigBytes := append(signature.Serialize(), byte(txscript.SigHashAll))
	pubKey := pubKeyObj.SerializeCompressed()

	// Debug: Compare public key hashes for P2PKH
	if len(witnessUtxo.PkScript) == 25 { // P2PKH scriptPubKey length
		expectedPKH := witnessUtxo.PkScript[3:23]
		actualPKH := btcutil.Hash160(pubKey)
		logger.Debug("Comparing public key hashes for P2PKH",
			"input", inputIndex,
			"expectedPKH", hex.EncodeToString(expectedPKH),
			"actualPKH", hex.EncodeToString(actualPKH),
			"match", bytes.Equal(expectedPKH, actualPKH))
	}

	// Log BIP32 derivation for debugging
	if len(psbtInput.Bip32Derivation) > 0 {
		logger.Debug("PSBT BIP32 derivation info",
			"input", inputIndex,
			"derivation_count", len(psbtInput.Bip32Derivation))
		for i, deriv := range psbtInput.Bip32Derivation {
			logger.Debug("BIP32 derivation entry",
				"input", inputIndex,
				"index", i,
				"pubKey", hex.EncodeToString(deriv.PubKey),
				"match", bytes.Equal(deriv.PubKey, pubKey))
		}
	}

	// For P2PKH, btcd's updater.Sign() has issues, so add signature directly
	logger.Debug("Adding P2PKH signature directly to PartialSigs (bypassing updater.Sign)",
		"input", inputIndex)

	// Initialize PartialSigs slice if needed
	if psbtInput.PartialSigs == nil {
		psbtInput.PartialSigs = make([]*psbt.PartialSig, 0)
	}

	// Check if signature already exists for this public key
	for _, partialSig := range psbtInput.PartialSigs {
		if bytes.Equal(partialSig.PubKey, pubKey) {
			logger.Debug("Signature already exists for this pubkey", "input", inputIndex)
			return false
		}
	}

	// Add partial signature directly
	psbtInput.PartialSigs = append(psbtInput.PartialSigs, &psbt.PartialSig{
		PubKey:    pubKey,
		Signature: sigBytes,
	})
	logger.Debug("Added P2PKH signature directly to PartialSigs",
		"input", inputIndex,
		"pubKey", hex.EncodeToString(pubKey))

	return true
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

// addBIP32DerivationForInput adds BIP32 derivation path information to a PSBT input.
// This is required for descriptor-based signing to work correctly.
// It extracts the address from the scriptPubKey, queries Bitcoin Core for derivation info,
// and adds it to the PSBT.
func (b *Bitcoin) addBIP32DerivationForInput(
	updater *psbt.Updater,
	scriptPubKey []byte,
	inputIndex int,
	senderAccount domainAccount.AccountType,
) error {
	// Extract address from scriptPubKey
	_, addrs, _, err := txscript.ExtractPkScriptAddrs(scriptPubKey, b.chainConf)
	if err != nil {
		return fmt.Errorf("failed to extract address from scriptPubKey: %w", err)
	}
	if len(addrs) == 0 {
		return errors.New("no address found in scriptPubKey")
	}

	// Get address info to determine the derivation path
	addressStr := addrs[0].EncodeAddress()
	logger.Info("Adding BIP32 derivation for input",
		"input_index", inputIndex,
		"address", addressStr,
		"sender_account", senderAccount.String())

	addressInfo, err := b.GetAddressInfo(addressStr)
	if err != nil {
		logger.Error("GetAddressInfo failed", "address", addressStr, "error", err)
		return fmt.Errorf("failed to get address info for %s: %w", addressStr, err)
	}

	logger.Info("Got address info",
		"address", addressStr,
		"hd_key_path", addressInfo.HDKeyPath,
		"hd_fingerprint", addressInfo.HDMasterFingerprint,
		"pubkey_len", len(addressInfo.PubKey))

	// Check if address has derivation path
	if addressInfo.HDKeyPath == "" {
		return errors.New("address has no HD key path (not from descriptor wallet?)")
	}
	if addressInfo.PubKey == "" {
		return errors.New("address has no public key")
	}

	// Parse the public key
	pubKeyBytes, err := hex.DecodeString(addressInfo.PubKey)
	if err != nil {
		return fmt.Errorf("failed to decode public key: %w", err)
	}

	// Get the correct fingerprint from imported descriptors, not from Bitcoin Core's internal wallet
	// Query list of descriptors to find the one that matches this address and account
	fingerprint, _, err := b.getDescriptorInfoForAddress(addressStr, addressInfo.HDKeyPath, senderAccount)
	if err != nil {
		return fmt.Errorf("failed to get descriptor info: %w", err)
	}

	// Parse the full derivation path
	// Use addressInfo.HDKeyPath directly as it already contains the full path from getaddressinfo
	derivationPath, err := parseDerivationPath(addressInfo.HDKeyPath)
	if err != nil {
		return fmt.Errorf("failed to parse derivation path: %w", err)
	}

	// Add BIP32 derivation to PSBT
	// updater.AddInBip32Derivation signature: (fingerprint uint32, path []uint32, pubkey []byte, inputIndex int)
	if err := updater.AddInBip32Derivation(fingerprint, derivationPath, pubKeyBytes, inputIndex); err != nil {
		return fmt.Errorf("failed to add BIP32 derivation to PSBT: %w", err)
	}

	logger.Debug("Added BIP32 derivation to PSBT input",
		"input_index", inputIndex,
		"address", addrs[0].EncodeAddress(),
		"path", addressInfo.HDKeyPath,
		"fingerprint", addressInfo.HDMasterFingerprint,
	)

	return nil
}

// parseDerivationPath converts a BIP32 derivation path string to a uint32 array.
// Input format: "m/44h/1h/1h/0/0" or "m/44'/1'/1'/0/0"
// Output format: []uint32 with hardened keys having high bit set (0x80000000)
func parseDerivationPath(path string) ([]uint32, error) {
	// Remove "m/" prefix if present
	if len(path) >= 2 && path[:2] == "m/" {
		path = path[2:]
	}

	// Split by "/"
	parts := make([]string, 0)
	for _, part := range splitPath(path) {
		if part != "" {
			parts = append(parts, part)
		}
	}

	if len(parts) == 0 {
		return nil, errors.New("empty derivation path")
	}

	// Parse each component
	result := make([]uint32, len(parts))
	for i, part := range parts {
		// Check if hardened (ends with 'h' or apostrophe)
		hardened := false
		if len(part) > 0 && (part[len(part)-1] == 'h' || part[len(part)-1] == '\'') {
			hardened = true
			part = part[:len(part)-1]
		}

		// Parse the number
		var num uint32
		_, err := fmt.Sscanf(part, "%d", &num)
		if err != nil {
			return nil, fmt.Errorf("failed to parse path component %q: %w", part, err)
		}

		// Set hardened bit if needed
		if hardened {
			num |= 0x80000000 // Set high bit for hardened derivation
		}

		result[i] = num
	}

	return result, nil
}

// splitPath splits a derivation path by "/" separator
func splitPath(path string) []string {
	parts := make([]string, 0)
	current := ""
	for _, ch := range path {
		if ch == '/' {
			parts = append(parts, current)
			current = ""
		} else {
			current += string(ch)
		}
	}
	if current != "" {
		parts = append(parts, current)
	}
	return parts
}

// getDescriptorInfoForAddress queries loaded descriptors to find the correct fingerprint
// for the given address. This is necessary because getaddressinfo returns Bitcoin Core's
// internal wallet fingerprint, not the imported descriptor's fingerprint.
func (b *Bitcoin) getDescriptorInfoForAddress(address, fullPath string, senderAccount domainAccount.AccountType) (uint32, string, error) {
	// Call listdescriptors RPC with false to get public descriptors
	// Note: Watch-only wallets can't return private descriptors (true parameter would fail)
	rawResult, err := b.Client.RawRequest("listdescriptors", []json.RawMessage{json.RawMessage("false")})
	if err != nil {
		return 0, "", fmt.Errorf("failed to call listdescriptors: %w", err)
	}

	var result struct {
		Descriptors []struct {
			Desc     string `json:"desc"`
			Active   bool   `json:"active"`
			Internal *bool  `json:"internal"`
		} `json:"descriptors"`
	}
	if err := json.Unmarshal(rawResult, &result); err != nil {
		return 0, "", fmt.Errorf("failed to unmarshal listdescriptors result: %w", err)
	}

	// Extract the base path from fullPath (e.g., "m/44h/1h/1h" from "m/44h/1h/1h/0/0")
	// The last two components are change/index
	pathParts := strings.Split(strings.TrimPrefix(fullPath, "m/"), "/")
	if len(pathParts) < 3 {
		return 0, "", fmt.Errorf("invalid derivation path: %s", fullPath)
	}
	// Get the account-level path (e.g., "44h/1h/1h")
	accountPath := strings.Join(pathParts[:len(pathParts)-2], "/")

	// Prepare path variations for matching
	// Both 'h' notation (44h/1h/1h) and apostrophe notation (44'/1'/1') are valid
	// listdescriptors may return either format depending on how descriptors were imported
	accountPathApostrophe := strings.ReplaceAll(accountPath, "h", "'")

	// Get expected BIP44 account index for validation
	expectedAccountIndex := senderAccount.BIP44AccountIndex()
	logger.Debug("Searching for descriptor",
		"account", senderAccount.String(),
		"expected_account_index", expectedAccountIndex,
		"account_path", accountPath,
		"normalized_path", accountPathApostrophe,
	)

	// Find the descriptor that matches this account path
	for _, desc := range result.Descriptors {
		// Parse descriptor to extract fingerprint and path
		// Format: "pkh([fingerprint/path]xpub.../change/*)"
		// Match using both 'h' notation and apostrophe notation
		if !strings.Contains(desc.Desc, accountPath) && !strings.Contains(desc.Desc, accountPathApostrophe) {
			continue
		}

		// Verify this is an active descriptor (not internal change addresses)
		if desc.Internal != nil && *desc.Internal {
			logger.Debug("Skipping internal descriptor", "desc", desc.Desc[:50])
			continue
		}

		// Extract fingerprint from descriptor
		// Format: [fingerprint/44h/1h/1h]
		start := strings.Index(desc.Desc, "[")
		end := strings.Index(desc.Desc, "/")
		if start == -1 || end == -1 || end <= start {
			continue
		}

		fingerprintHex := desc.Desc[start+1 : end]
		fingerprintBytes, err := hex.DecodeString(fingerprintHex)
		if err != nil || len(fingerprintBytes) != 4 {
			continue
		}

		// Convert to uint32 (big-endian)
		fingerprint := uint32(fingerprintBytes[0])<<24 |
			uint32(fingerprintBytes[1])<<16 |
			uint32(fingerprintBytes[2])<<8 |
			uint32(fingerprintBytes[3])

		// Extract base path (e.g., "/44h/1h/1h")
		pathStart := strings.Index(desc.Desc, "/")
		pathEnd := strings.Index(desc.Desc, "]")
		if pathStart == -1 || pathEnd == -1 || pathEnd <= pathStart {
			continue
		}
		basePath := "m" + desc.Desc[pathStart:pathEnd]

		// Verify the account index in the path matches what we expect
		// basePath format: "m/44h/1h/1h" where the last component is the account index
		basePathParts := strings.Split(strings.TrimPrefix(basePath, "m/"), "/")
		if len(basePathParts) >= 3 {
			// Parse the account component (e.g., "1h" -> 1)
			accountComponent := basePathParts[2]
			accountComponent = strings.TrimSuffix(accountComponent, "h")
			accountComponent = strings.TrimSuffix(accountComponent, "'")

			var parsedAccountIndex uint32
			_, err := fmt.Sscanf(accountComponent, "%d", &parsedAccountIndex)
			if err == nil && parsedAccountIndex == expectedAccountIndex {
				logger.Info("Successfully matched descriptor for sender account",
					"address", address,
					"account", senderAccount.String(),
					"fingerprint", fingerprintHex,
					"base_path", basePath,
					"account_index", parsedAccountIndex,
				)
				return fingerprint, basePath, nil
			} else {
				logger.Debug("Account index mismatch",
					"parsed", parsedAccountIndex,
					"expected", expectedAccountIndex,
					"base_path", basePath,
				)
			}
		}
	}

	return 0, "", fmt.Errorf("no matching descriptor found for address %s with path %s and account %s (expected account index %d)",
		address, fullPath, senderAccount.String(), expectedAccountIndex)
}
