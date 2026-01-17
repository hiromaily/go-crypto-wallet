// Package btc provides Bitcoin (BTC) infrastructure implementations.
//
// IMPORTANT: PSBT (Partially Signed Bitcoin Transaction) - BTC ONLY
//
// This file implements PSBT functionality which is ONLY applicable to Bitcoin (BTC).
// BCH (Bitcoin Cash) does NOT support PSBT.
//
// ┌─────────────────────────────────────────────────────────────────────────────┐
// │ WARNING: DO NOT USE THIS FILE FOR BCH (Bitcoin Cash) IMPLEMENTATIONS        │
// │                                                                             │
// │ BCH uses Raw Transaction Hex format, NOT PSBT.                              │
// │ BCH signing requires SIGHASH_FORKID (0x40) which is not handled here.       │
// │                                                                             │
// │ For BCH implementations, see:                                               │
// │   - internal/infrastructure/api/btc/bch/                                    │
// │   - docs/task-contexts/chains/bch.md                                        │
// └─────────────────────────────────────────────────────────────────────────────┘
//
// PSBT is defined in BIP174 and provides a standard format for unsigned/partially
// signed transactions that can be passed between parties for multi-signature workflows.
//
// BCH Protocol Differences (DO NOT apply PSBT to BCH):
//   - BCH: Raw TX Hex + ECDSA + SIGHASH_FORKID (0x41)
//   - BTC: PSBT (BIP174) + ECDSA/Schnorr + standard sighash types
package btc

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcec/v2/ecdsa"
	"github.com/btcsuite/btcd/btcec/v2/schnorr"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/btcutil/psbt"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"

	dtobtc "github.com/hiromaily/go-crypto-wallet/internal/application/dto/btc"
	domainAccount "github.com/hiromaily/go-crypto-wallet/internal/domain/account"
	domainCoin "github.com/hiromaily/go-crypto-wallet/internal/domain/coin"
	domainWallet "github.com/hiromaily/go-crypto-wallet/internal/domain/wallet"
	"github.com/hiromaily/go-crypto-wallet/pkg/logger"
)

const (
	// p2wpkhRedeemScriptLen is the expected length of a P2WPKH redeemScript
	// Format: OP_0 (1 byte) + length byte (1 byte) + pubkey hash (20 bytes) = 22 bytes
	p2wpkhRedeemScriptLen = 22
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
func (b *Bitcoin) CreatePSBT(
	msgTx *wire.MsgTx, prevTxs []dtobtc.PreviousTx, senderAccount domainAccount.AccountType,
) (string, error) {
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
		logger.Info("Processing input for PSBT",
			"input", i,
			"has_redeem_script", prevTx.RedeemScript != "",
			"redeem_script_len", len(prevTx.RedeemScript),
			"has_witness_script", prevTx.WitnessScript != "")

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
		logger.Debug("Checking redeemScript for input",
			"input", i,
			"has_redeem_script", prevTx.RedeemScript != "",
			"has_witness_script", len(witnessScript) > 0,
			"is_p2sh", txscript.IsPayToScriptHash(scriptPubKey))

		if prevTx.RedeemScript != "" {
			redeemScript, err := b.decodeHexScript(prevTx.RedeemScript)
			if err != nil {
				return "", fmt.Errorf("failed to decode redeemScript for input %d: %w", i, err)
			}
			if err := updater.AddInRedeemScript(redeemScript, i); err != nil {
				return "", fmt.Errorf("failed to add redeem script for input %d: %w", i, err)
			}
			logger.Debug("Added redeemScript from prevTx", "input", i, "script_len", len(redeemScript))
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
		} else if txscript.IsPayToScriptHash(scriptPubKey) && len(witnessScript) == 0 {
			// P2SH (non-SegWit) multisig - derive redeemScript from descriptor
			// Bitcoin Core's listunspent doesn't return redeemScript for descriptor-based addresses

			// Extract address from scriptPubKey
			_, addrs, _, err := txscript.ExtractPkScriptAddrs(scriptPubKey, b.chainConf)
			if err != nil || len(addrs) == 0 {
				return "", fmt.Errorf("failed to extract address from scriptPubKey for input %d: %w", i, err)
			}
			address := addrs[0].EncodeAddress()

			logger.Info("RedeemScript missing for P2SH input, attempting to derive from descriptor",
				"input", i,
				"address", address)

			redeemScript, err := b.deriveRedeemScriptForAddress(address, senderAccount)
			if err != nil {
				return "", fmt.Errorf("failed to derive redeemScript for input %d: %w", i, err)
			}

			if err := updater.AddInRedeemScript(redeemScript, i); err != nil {
				return "", fmt.Errorf("failed to add derived redeem script for input %d: %w", i, err)
			}
			logger.Info("Derived and added RedeemScript for P2SH multisig input",
				"input", i,
				"address", address,
				"script_len", len(redeemScript))
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
			return "", fmt.Errorf(
				"failed to add BIP32 derivation for input %d (required for descriptor-based signing): %w",
				i, err,
			)
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
			"partialSigsCount", len(input.PartialSigs))

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
		(txscript.IsPayToScriptHash(witnessUtxo.PkScript) && len(psbtInput.WitnessScript) > 0) ||
		(txscript.IsPayToScriptHash(witnessUtxo.PkScript) && len(psbtInput.RedeemScript) > 0 &&
			txscript.IsPayToWitnessPubKeyHash(psbtInput.RedeemScript))

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

// buildP2PKHScriptCodeForP2SHWPKH validates a P2WPKH redeemScript and builds
// the P2PKH scriptCode for signature hash calculation per BIP143.
//
// For P2SH-P2WPKH, the redeemScript is: OP_0 <20-byte-pubkey-hash>
// This function validates the format and builds the P2PKH scriptCode:
// OP_DUP OP_HASH160 <20-byte-hash> OP_EQUALVERIFY OP_CHECKSIG
func buildP2PKHScriptCodeForP2SHWPKH(redeemScript []byte, inputIndex int, utxoValue int64) ([]byte, error) {
	// Validate length (must be 22 bytes: OP_0 + length byte + 20-byte hash)
	if len(redeemScript) != p2wpkhRedeemScriptLen {
		return nil, fmt.Errorf("invalid P2WPKH redeemScript length: expected %d, got %d",
			p2wpkhRedeemScriptLen, len(redeemScript))
	}

	// Validate format: must be OP_0 <20-byte-hash>
	// Format: [OP_0 (0x00), length byte (0x14 = 20 decimal), 20-byte pubkey hash]
	if redeemScript[0] != txscript.OP_0 {
		return nil, fmt.Errorf("invalid P2WPKH redeemScript format: must start with OP_0 (0x00), got 0x%02x",
			redeemScript[0])
	}

	// Validate length byte (0x14 = 20 bytes for pubkey hash)
	const pubKeyHashLength = 20
	if redeemScript[1] != pubKeyHashLength {
		return nil, fmt.Errorf("invalid P2WPKH redeemScript format: invalid length byte, expected 0x14, got 0x%02x",
			redeemScript[1])
	}

	// Extract pubkey hash (now safe after validation)
	pubKeyHash := redeemScript[2:]

	// Build P2PKH scriptCode per BIP143
	// Format: OP_DUP OP_HASH160 <20-byte-hash> OP_EQUALVERIFY OP_CHECKSIG
	scriptCodeBuilder := txscript.NewScriptBuilder()
	scriptCodeBuilder.AddOp(txscript.OP_DUP)
	scriptCodeBuilder.AddOp(txscript.OP_HASH160)
	scriptCodeBuilder.AddData(pubKeyHash)
	scriptCodeBuilder.AddOp(txscript.OP_EQUALVERIFY)
	scriptCodeBuilder.AddOp(txscript.OP_CHECKSIG)

	scriptForHash, err := scriptCodeBuilder.Script()
	if err != nil {
		return nil, fmt.Errorf("failed to build P2PKH scriptCode: %w", err)
	}

	logger.Info("P2SH-P2WPKH signing - constructed P2PKH scriptCode",
		"input", inputIndex,
		"redeemScript_hex", hex.EncodeToString(redeemScript),
		"pubKeyHash_hex", hex.EncodeToString(pubKeyHash),
		"scriptCode_hex", hex.EncodeToString(scriptForHash),
		"witnessUtxo_value", utxoValue)

	return scriptForHash, nil
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

	// Determine the script to use for hash calculation per BIP143:
	// - P2WSH multisig: use the witness script
	// - P2WPKH (native): use the scriptPubKey (witness program)
	// - P2SH-P2WPKH (nested): use the redeem script (which IS the witness program)
	// CalcWitnessSigHash will convert the witness program to scriptCode internally
	scriptForHash := witnessUtxo.PkScript
	if len(psbtInput.WitnessScript) > 0 {
		// P2WSH multisig
		scriptForHash = psbtInput.WitnessScript
		logger.Debug("Using witness script for P2WSH multisig hash calculation",
			"input", inputIndex,
			"witnessScriptLen", len(psbtInput.WitnessScript))
	} else if len(psbtInput.RedeemScript) > 0 && txscript.IsPayToWitnessPubKeyHash(psbtInput.RedeemScript) {
		// P2SH-P2WPKH: Build P2PKH scriptCode from redeemScript per BIP143
		var err error
		scriptForHash, err = buildP2PKHScriptCodeForP2SHWPKH(psbtInput.RedeemScript, inputIndex, witnessUtxo.Value)
		if err != nil {
			logger.Error("Failed to build P2PKH scriptCode for P2SH-P2WPKH",
				"input", inputIndex,
				"error", err)
			return false
		}
	}

	logger.Info("Calculating witness signature hash",
		"input", inputIndex,
		"scriptForHash_hex", hex.EncodeToString(scriptForHash),
		"scriptForHash_len", len(scriptForHash),
		"amount_satoshis", witnessUtxo.Value)

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

	logger.Info("Witness signature hash calculated",
		"input", inputIndex,
		"hash_hex", hex.EncodeToString(hash),
		"hash_len", len(hash))

	// Sign with ECDSA
	signature := ecdsa.Sign(privKey.PrivKey, hash)
	sigBytes := append(signature.Serialize(), byte(txscript.SigHashAll))
	pubKeyObj := privKey.PrivKey.PubKey()
	pubKey := pubKeyObj.SerializeCompressed()

	logger.Info("ECDSA signature created",
		"input", inputIndex,
		"signature_hex", hex.EncodeToString(sigBytes),
		"signature_len", len(sigBytes),
		"pubkey_hex", hex.EncodeToString(pubKey),
		"pubkey_len", len(pubKey))

	// Verify signature locally before adding
	if !signature.Verify(hash, pubKeyObj) {
		logger.Error("SegWit signature verification FAILED locally", "input", inputIndex)
		return false
	}
	logger.Info("SegWit signature verified locally - signature is VALID",
		"input", inputIndex,
		"pubkey_hash", hex.EncodeToString(btcutil.Hash160(pubKey)))

	// For P2SH-P2WPKH, add signature directly (btcd's updater may have issues)
	// For other SegWit types, use updater.Sign()
	isP2SHP2WPKH := len(psbtInput.RedeemScript) > 0 && txscript.IsPayToWitnessPubKeyHash(psbtInput.RedeemScript)
	if isP2SHP2WPKH {
		logger.Debug("Adding P2SH-P2WPKH signature directly to PartialSigs", "input", inputIndex)

		// Initialize PartialSigs if needed
		if psbtInput.PartialSigs == nil {
			psbtInput.PartialSigs = make([]*psbt.PartialSig, 0)
		}

		// Check if signature already exists
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
		logger.Debug("Added P2SH-P2WPKH signature to PartialSigs",
			"input", inputIndex,
			"pubKey", hex.EncodeToString(pubKey))
		return true
	}

	// For other SegWit types (P2WSH, native P2WPKH), use btcd's updater
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

// addBIP32DerivationFromDescriptor adds BIP32 derivation paths to PSBT from descriptor.
// This is required for descriptor-based signing when Bitcoin Core doesn't return HDKeyPath.
func (b *Bitcoin) addBIP32DerivationFromDescriptor(
	updater *psbt.Updater,
	address string,
	inputIndex int,
	senderAccount domainAccount.AccountType,
) error {
	logger.Debug("Deriving BIP32 paths from descriptor",
		"address", address,
		"account", senderAccount.String())

	// List all descriptors to find the original descriptor with xpubs
	// Note: getaddressinfo returns a descriptor with raw public keys, not xpubs
	descriptorList, err := b.ListDescriptors(false)
	if err != nil {
		return fmt.Errorf("failed to list descriptors: %w", err)
	}

	// Get address info to check which descriptor this address belongs to
	addressInfo, err := b.GetAddressInfo(address)
	if err != nil {
		return fmt.Errorf("failed to get address info: %w", err)
	}

	if addressInfo.Desc == "" {
		return errors.New("address info does not contain descriptor (not a descriptor wallet address?)")
	}

	logger.Debug("Got descriptor from getaddressinfo (with raw keys)",
		"address", address,
		"descriptor_len", len(addressInfo.Desc))

	// Find the matching wallet descriptor (with xpubs) by comparing the structure
	// We need to match descriptors based on their fingerprints
	var walletDescriptor string
	var isInternal bool
	var addressIndex uint32

	for _, desc := range descriptorList.Descriptors {
		// Try to find address in this descriptor
		foundIndex, err := b.findAddressIndexInDescriptor(desc.Desc, address)
		if err == nil {
			// Found it!
			walletDescriptor = desc.Desc
			addressIndex = foundIndex
			isInternal = desc.Internal != nil && *desc.Internal
			logger.Info("Found matching wallet descriptor",
				"address", address,
				"address_index", addressIndex,
				"is_internal", isInternal,
				"descriptor_prefix", desc.Desc[:80]+"...")

			// DIAGNOSTIC: Verify what Bitcoin Core returns for this descriptor at this index
			addresses, err := b.deriveAddressesFromDescriptor(desc.Desc, addressIndex, addressIndex)
			if err == nil && len(addresses) > 0 {
				logger.Info("Bitcoin Core deriveaddresses verification",
					"descriptor_index", addressIndex,
					"bitcoin_core_address", addresses[0],
					"target_address", address,
					"match", addresses[0] == address)
			} else {
				logger.Warn("Failed to verify with Bitcoin Core deriveaddresses", "error", err)
			}

			// Parse the wallet descriptor (which has xpubs)
			parser := NewDescriptorParser()
			parsed, err := parser.Parse(walletDescriptor)
			if err != nil {
				return fmt.Errorf("failed to parse wallet descriptor: %w", err)
			}

			// Support all descriptor types: P2SH-wrapped and Native SegWit
			// P2SH-wrapped: sh (P2SH multisig), sh(wpkh) (P2SH-P2WPKH), sh(wsh) (P2SH-P2WSH)
			// Native SegWit: wpkh (P2WPKH), wsh (P2WSH)
			if parsed.Type != domainWallet.DescriptorTypeSH &&
				parsed.Type != domainWallet.DescriptorTypeSHWPKH &&
				parsed.Type != domainWallet.DescriptorTypeSHWSH &&
				parsed.Type != domainWallet.DescriptorTypeWPKH &&
				parsed.Type != domainWallet.DescriptorTypeWSH {
				return fmt.Errorf(
					"unsupported descriptor type: %s (only sh/sh(wpkh)/sh(wsh)/wpkh/wsh supported)",
					parsed.Type,
				)
			}

			// Verify descriptor by deriving the address
			// Different logic for P2SH vs Native SegWit
			var derivedAddr string

			if parsed.Type == domainWallet.DescriptorTypeWSH || parsed.Type == domainWallet.DescriptorTypeWPKH {
				// Native SegWit - derive Bech32 address from witness script/program
				derivedAddr, err = b.deriveNativeSegWitAddress(parsed, addressIndex)
				if err != nil {
					return fmt.Errorf("failed to derive native SegWit address: %w", err)
				}
			} else {
				// P2SH - derive P2SH address from redeem script
				redeemScript, err := b.DeriveRedeemScriptFromDescriptor(walletDescriptor, address, addressIndex)
				if err != nil {
					return fmt.Errorf("failed to derive redeemScript: %w", err)
				}

				derivedAddr, err = b.deriveP2SHAddressFromRedeemScript(redeemScript)
				if err != nil {
					return fmt.Errorf("failed to derive P2SH address from redeemScript: %w", err)
				}
			}

			if derivedAddr != address {
				return fmt.Errorf("derived address %s does not match target %s", derivedAddr, address)
			}

			logger.Info("Verified address matches descriptor, adding BIP32 derivation for all keys",
				"address", address,
				"descriptor_type", parsed.Type,
				"descriptor_index", addressIndex,
				"num_keys", len(parsed.Keys))

			// Add BIP32 derivation for each key in the descriptor
			// For sh(wpkh(...)), there is only 1 key
			// For sh(multi(...)), there are multiple keys
			for keyIdx, keyInfo := range parsed.Keys {
				// Derive the public key at this index
				pubKey, err := b.derivePublicKeyFromDescriptorKey(keyInfo, addressIndex)
				if err != nil {
					return fmt.Errorf("failed to derive public key %d: %w", keyIdx, err)
				}

				// Parse the derivation path to get the full path
				// Format: fingerprint from descriptor + derivation path
				fingerprint, err := hex.DecodeString(keyInfo.Fingerprint)
				if err != nil {
					return fmt.Errorf("invalid fingerprint: %w", err)
				}

				// Build the full derivation path from master key
				// Format: OriginPath + DerivationPath (with wildcard replaced)
				// Example: "/44'/1'/1'" + "/0/*" = "/44'/1'/1'/0/0" (for addressIndex=0)
				relativePath := keyInfo.DerivationPath
				if relativePath == "" {
					relativePath = fmt.Sprintf("/%d", addressIndex)
				} else {
					// Replace wildcard with actual index
					relativePath = strings.ReplaceAll(relativePath, "/*", fmt.Sprintf("/%d", addressIndex))
				}

				// Combine origin and relative paths for full BIP32 path
				fullPath := keyInfo.OriginPath + relativePath

				// Parse full derivation path into []uint32
				pathIndices, err := b.parseDerivationPath(fullPath)
				if err != nil {
					return fmt.Errorf("failed to parse derivation path %s: %w", fullPath, err)
				}

				// Add BIP32 derivation to PSBT
				// AddInBip32Derivation signature: (fingerprint uint32, path []uint32, pubkey []byte, inputIndex int)
				fingerprintUint32 := binary.LittleEndian.Uint32(fingerprint)

				if err := updater.AddInBip32Derivation(fingerprintUint32, pathIndices, pubKey, inputIndex); err != nil {
					return fmt.Errorf("failed to add BIP32 derivation for key %d: %w", keyIdx, err)
				}

				logger.Debug("Added BIP32 derivation for descriptor key",
					"input", inputIndex,
					"descriptor_type", parsed.Type,
					"key_index", keyIdx,
					"pubkey", hex.EncodeToString(pubKey),
					"fingerprint", keyInfo.Fingerprint,
					"full_path", fullPath,
					"origin", keyInfo.OriginPath,
					"relative", relativePath)
			}

			logger.Info("Successfully added BIP32 derivation for all descriptor keys",
				"address", address,
				"descriptor_type", parsed.Type,
				"input", inputIndex,
				"num_keys", len(parsed.Keys))
			return nil
		}
	}

	return fmt.Errorf("no matching wallet descriptor found for address %s", address)
}

// deriveRedeemScriptForAddress derives a redeemScript for a P2SH address by:
// 1. Finding the matching descriptor for the account
// 2. Searching for the address index within the descriptor range
// 3. Deriving the redeemScript at that index
//
// Performance characteristics:
//   - Lists all wallet descriptors (typically < 10)
//   - Iterates through up to 1,000 address indices per descriptor
//   - For each index, derives redeemScript and checks address match
//   - Time complexity: O(n * m) where n = descriptors, m = search range (1000)
//   - Typical time: < 1 second for addresses in first 100 indices
//   - Worst case: Several seconds for high-index addresses
//
// This is used as a fallback when Bitcoin Core's listunspent doesn't return
// redeemScript for descriptor-based addresses. It's called once per PSBT creation
// for P2SH multisig inputs.
func (b *Bitcoin) deriveRedeemScriptForAddress(
	address string, senderAccount domainAccount.AccountType,
) ([]byte, error) {
	// List all descriptors
	descriptorList, err := b.ListDescriptors(false)
	if err != nil {
		return nil, fmt.Errorf("failed to list descriptors: %w", err)
	}

	// Determine expected account index
	var expectedAccountIndex uint32
	switch senderAccount {
	case domainAccount.AccountTypeDeposit:
		expectedAccountIndex = 0
	case domainAccount.AccountTypePayment:
		expectedAccountIndex = 1
	case domainAccount.AccountTypeStored:
		expectedAccountIndex = 2
	case domainAccount.AccountTypeClient,
		domainAccount.AccountTypeAuthorization,
		domainAccount.AccountTypeAuth1, domainAccount.AccountTypeAuth2, domainAccount.AccountTypeAuth3,
		domainAccount.AccountTypeAuth4, domainAccount.AccountTypeAuth5, domainAccount.AccountTypeAuth6,
		domainAccount.AccountTypeAuth7, domainAccount.AccountTypeAuth8, domainAccount.AccountTypeAuth9,
		domainAccount.AccountTypeAuth10, domainAccount.AccountTypeAuth11, domainAccount.AccountTypeAuth12,
		domainAccount.AccountTypeAuth13, domainAccount.AccountTypeAuth14, domainAccount.AccountTypeAuth15,
		domainAccount.AccountTypeAnonymous, domainAccount.AccountTypeTest:
		return nil, fmt.Errorf("unsupported account type for this operation: %s", senderAccount.String())
	default:
		return nil, fmt.Errorf("unknown account type: %s", senderAccount.String())
	}

	logger.Debug("Searching for descriptor matching address",
		"address", address,
		"account", senderAccount.String(),
		"account_index", expectedAccountIndex)

	// Find matching descriptor and derive redeemScript
	for _, desc := range descriptorList.Descriptors {
		// Skip internal (change) descriptors
		if desc.Internal != nil && *desc.Internal {
			continue
		}

		// Check if descriptor matches the account
		if !b.descriptorMatchesAccountIndex(desc.Desc, expectedAccountIndex) {
			continue
		}

		logger.Debug("Found matching descriptor, searching for address",
			"descriptor_len", len(desc.Desc),
			"account_index", expectedAccountIndex)

		// Try to find the address by deriving up to 1000 addresses (default range)
		for i := uint32(0); i < 1000; i++ {
			redeemScript, err := b.DeriveRedeemScriptFromDescriptor(desc.Desc, address, i)
			if err != nil {
				// This index doesn't work, try next
				continue
			}

			// Verify the derived redeemScript matches the address
			derivedAddr, err := b.deriveP2SHAddressFromRedeemScript(redeemScript)
			if err != nil {
				continue
			}

			if derivedAddr == address {
				logger.Info("Successfully derived redeemScript for address",
					"address", address,
					"descriptor_index", i,
					"script_len", len(redeemScript))
				return redeemScript, nil
			}
		}
	}

	return nil, fmt.Errorf("no matching descriptor found for address %s (account=%s)", address, senderAccount.String())
}

// deriveP2SHAddressFromRedeemScript derives a P2SH address from a redeemScript
func (b *Bitcoin) deriveP2SHAddressFromRedeemScript(redeemScript []byte) (string, error) {
	// Hash the redeemScript
	scriptHash := btcutil.Hash160(redeemScript)

	// Create P2SH address from hash
	// Note: Use NewAddressScriptHashFromHash since we already hashed the script
	address, err := btcutil.NewAddressScriptHashFromHash(scriptHash, b.chainConf)
	if err != nil {
		return "", fmt.Errorf("failed to create P2SH address: %w", err)
	}

	return address.EncodeAddress(), nil
}

// deriveNativeSegWitAddress derives a Bech32 address from a Native SegWit descriptor.
// Supports both P2WPKH (single-sig) and P2WSH (multisig) descriptors.
//
// For P2WPKH:
//   - Derives public key at addressIndex
//   - Hashes public key (HASH160)
//   - Creates Bech32 address from 20-byte witness program
//   - Format: bc1q<20-byte-hash> (42 chars mainnet), bcrt1q<20-byte-hash> (regtest)
//
// For P2WSH:
//   - Derives witness script (multisig script) at addressIndex
//   - Hashes witness script (SHA256)
//   - Creates Bech32 address from 32-byte witness program
//   - Format: bc1q<32-byte-hash> (62 chars mainnet), bcrt1q<32-byte-hash> (regtest)
func (b *Bitcoin) deriveNativeSegWitAddress(parsed *domainWallet.Descriptor, addressIndex uint32) (string, error) {
	//nolint:exhaustive // Only WPKH and WSH are native SegWit types
	switch parsed.Type {
	case domainWallet.DescriptorTypeWPKH:
		// P2WPKH: single-sig Native SegWit
		if len(parsed.Keys) != 1 {
			return "", fmt.Errorf("P2WPKH descriptor must have exactly 1 key, got %d", len(parsed.Keys))
		}

		// Derive the public key at this index
		pubKey, err := b.derivePublicKeyFromDescriptorKey(parsed.Keys[0], addressIndex)
		if err != nil {
			return "", fmt.Errorf("failed to derive public key: %w", err)
		}

		// Hash the public key (HASH160 = RIPEMD160(SHA256(pubkey)))
		pubKeyHash := btcutil.Hash160(pubKey)

		// Create P2WPKH address
		address, err := btcutil.NewAddressWitnessPubKeyHash(pubKeyHash, b.chainConf)
		if err != nil {
			return "", fmt.Errorf("failed to create P2WPKH address: %w", err)
		}

		return address.EncodeAddress(), nil

	case domainWallet.DescriptorTypeWSH:
		// P2WSH: multisig Native SegWit
		// Derive the witness script (multisig script)
		witnessScript, err := b.deriveWitnessScriptFromDescriptor(parsed, addressIndex)
		if err != nil {
			return "", fmt.Errorf("failed to derive witness script: %w", err)
		}

		// Hash the witness script (SHA256)
		witnessScriptHash := sha256.Sum256(witnessScript)

		// Create P2WSH address
		address, err := btcutil.NewAddressWitnessScriptHash(witnessScriptHash[:], b.chainConf)
		if err != nil {
			return "", fmt.Errorf("failed to create P2WSH address: %w", err)
		}

		return address.EncodeAddress(), nil

	default:
		return "", fmt.Errorf("unsupported native SegWit descriptor type: %s", parsed.Type)
	}
}

// deriveWitnessScriptFromDescriptor derives a witness script (multisig script) from a P2WSH descriptor.
// The witness script defines the spending conditions (e.g., 2-of-3 multisig).
//
// Format for sortedmulti(M, key1, key2, ...):
//   - M (threshold as OP_N)
//   - For each key: <pubkey>
//   - N (total keys as OP_N)
//   - OP_CHECKMULTISIG
//
// Example for 2-of-3:
//
//	OP_2 <pubkey1> <pubkey2> <pubkey3> OP_3 OP_CHECKMULTISIG
func (b *Bitcoin) deriveWitnessScriptFromDescriptor(
	parsed *domainWallet.Descriptor,
	addressIndex uint32,
) ([]byte, error) {
	if len(parsed.Keys) < 2 {
		return nil, fmt.Errorf("multisig descriptor must have at least 2 keys, got %d", len(parsed.Keys))
	}

	// Derive all public keys at this index
	pubKeys := make([][]byte, len(parsed.Keys))
	for i, keyInfo := range parsed.Keys {
		pubKey, err := b.derivePublicKeyFromDescriptorKey(keyInfo, addressIndex)
		if err != nil {
			return nil, fmt.Errorf("failed to derive public key %d: %w", i, err)
		}
		pubKeys[i] = pubKey
	}

	// Sort public keys lexicographically (BIP67 sortedmulti)
	// This ensures deterministic key ordering
	sortedPubKeys := make([][]byte, len(pubKeys))
	copy(sortedPubKeys, pubKeys)
	sort.Slice(sortedPubKeys, func(i, j int) bool {
		return bytes.Compare(sortedPubKeys[i], sortedPubKeys[j]) < 0
	})

	// Extract threshold from descriptor script
	// Format: wsh(sortedmulti(M, key1, key2, ...))
	// We need to parse "sortedmulti(M," to extract M
	threshold := 2 // Default for 2-of-3
	re := regexp.MustCompile(`sortedmulti\((\d+),`)
	matches := re.FindStringSubmatch(parsed.Script)
	if len(matches) > 1 {
		parsedThreshold, err := strconv.Atoi(matches[1])
		if err != nil {
			return nil, fmt.Errorf("failed to parse threshold from descriptor: %w", err)
		}
		threshold = parsedThreshold
	}

	// Build witness script using txscript
	builder := txscript.NewScriptBuilder()
	builder.AddInt64(int64(threshold)) // M (threshold)

	for _, pubKey := range sortedPubKeys {
		builder.AddData(pubKey)
	}

	builder.AddInt64(int64(len(sortedPubKeys))) // N (total keys)
	builder.AddOp(txscript.OP_CHECKMULTISIG)

	witnessScript, err := builder.Script()
	if err != nil {
		return nil, fmt.Errorf("failed to build witness script: %w", err)
	}

	return witnessScript, nil
}

// findAddressIndexInDescriptor uses Bitcoin Core's deriveaddresses RPC to find
// the index of an address within a descriptor's range.
//
// Performance characteristics:
//   - Time complexity: O(n) where n is the address index
//   - Searches in chunks of 100 addresses at a time
//   - Maximum search range: 10,000 addresses
//   - Early termination when address is found
//   - Average case: ~50 RPC calls for index 5000 (50 chunks * 100 addresses)
//   - Best case: 1 RPC call (address in first chunk)
//   - Worst case: 100 RPC calls (address at index 9,999)
//
// Optimization opportunities:
//   - Use Bitcoin Core's native address index lookup if available
//   - Cache descriptor derivation results
//   - Binary search if descriptor supports random access
func (b *Bitcoin) findAddressIndexInDescriptor(descriptor string, targetAddress string) (uint32, error) {
	// Search in chunks to avoid deriving too many addresses at once
	const chunkSize = 100
	const maxSearchRange = 10000 // Maximum range to search

	for startIdx := uint32(0); startIdx < maxSearchRange; startIdx += chunkSize {
		endIdx := startIdx + chunkSize - 1

		// Call deriveaddresses with range
		addresses, err := b.deriveAddressesFromDescriptor(descriptor, startIdx, endIdx)
		if err != nil {
			return 0, fmt.Errorf("failed to derive addresses [%d,%d]: %w", startIdx, endIdx, err)
		}

		// Search for target address in the chunk
		for i, addr := range addresses {
			if addr == targetAddress {
				return startIdx + uint32(i), nil
			}
		}
	}

	return 0, fmt.Errorf("address %s not found in descriptor range [0,%d]", targetAddress, maxSearchRange)
}

// deriveAddressesFromDescriptor calls Bitcoin Core's deriveaddresses RPC.
// It derives addresses from a descriptor within the specified range.
func (b *Bitcoin) deriveAddressesFromDescriptor(descriptor string, startIdx, endIdx uint32) ([]string, error) {
	// Build the RPC parameters
	// Format: deriveaddresses "descriptor" [start, end]
	rangeParam := fmt.Sprintf("[%d,%d]", startIdx, endIdx)

	// Call RPC
	params := []json.RawMessage{
		json.RawMessage(fmt.Sprintf(`"%s"`, descriptor)),
		json.RawMessage(rangeParam),
	}

	rawResult, err := b.Client.RawRequest("deriveaddresses", params)
	if err != nil {
		return nil, fmt.Errorf("deriveaddresses RPC failed: %w", err)
	}

	// Parse result
	var addresses []string
	if err := json.Unmarshal(rawResult, &addresses); err != nil {
		return nil, fmt.Errorf("failed to parse deriveaddresses result: %w", err)
	}

	logger.Debug("Derived addresses from descriptor",
		"range", fmt.Sprintf("[%d,%d]", startIdx, endIdx),
		"count", len(addresses))

	return addresses, nil
}

// parseDerivationPath parses a BIP32 derivation path string into a slice of indices.
// Format: "/0/5" or "/0/*" (wildcard should be replaced before calling)
//
//nolint:revive // receiver unused but method belongs to Bitcoin type
func (b *Bitcoin) parseDerivationPath(path string) ([]uint32, error) {
	if path == "" {
		return []uint32{}, nil
	}

	// Remove leading slash
	path = strings.TrimPrefix(path, "/")
	if path == "" {
		return []uint32{}, nil
	}

	// Split by slash
	parts := strings.Split(path, "/")
	indices := make([]uint32, 0, len(parts))

	for _, part := range parts {
		if part == "" {
			continue
		}

		// Check for hardened derivation (')
		hardened := false
		if strings.HasSuffix(part, "'") || strings.HasSuffix(part, "h") {
			hardened = true
			part = strings.TrimSuffix(part, "'")
			part = strings.TrimSuffix(part, "h")
		}

		// Parse the index
		var index uint32
		_, err := fmt.Sscanf(part, "%d", &index)
		if err != nil {
			return nil, fmt.Errorf("invalid path component %s: %w", part, err)
		}

		// Apply hardened bit if needed
		if hardened {
			index += 0x80000000
		}

		indices = append(indices, index)
	}

	return indices, nil
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
		"pubkey_len", len(addressInfo.PubKey),
		"is_watch_only", addressInfo.IsWatchOnly,
		"is_script", addressInfo.IsScript)

	// Check if address has derivation path
	// For watch-only multisig descriptors, HDKeyPath and PubKey may be empty
	// Bitcoin Core doesn't populate HDKeyPath for ranged multisig descriptors
	// In this case, derive BIP32 information from the descriptor
	if addressInfo.HDKeyPath == "" {
		if addressInfo.IsScript {
			// BCH uses legacy (non-descriptor) wallets and doesn't support descriptor-based signing
			// Skip BIP32 derivation from descriptor for BCH
			if b.coinTypeCode == domainCoin.BCH {
				logger.Debug("BCH: Skipping BIP32 derivation from descriptor (legacy wallet, not supported)",
					"input_index", inputIndex,
					"address", addressStr,
					"sender_account", senderAccount.String())
				return nil
			}

			logger.Info("Deriving BIP32 information from descriptor for multisig address",
				"address", addressStr,
				"input_index", inputIndex)

			// Derive BIP32 derivation paths from descriptor
			err := b.addBIP32DerivationFromDescriptor(updater, addressStr, inputIndex, senderAccount)
			if err != nil {
				logger.Error("Failed to derive BIP32 from descriptor",
					"address", addressStr,
					"error", err)
				return fmt.Errorf("failed to derive BIP32 from descriptor: %w", err)
			}
			logger.Info("Successfully added BIP32 derivation from descriptor",
				"address", addressStr,
				"input_index", inputIndex)
			return nil
		}
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
func (b *Bitcoin) getDescriptorInfoForAddress(
	address, fullPath string, senderAccount domainAccount.AccountType,
) (uint32, string, error) {
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

	return 0, "", fmt.Errorf(
		"no matching descriptor found for address %s with path %s and account %s (expected account index %d)",
		address, fullPath, senderAccount.String(), expectedAccountIndex,
	)
}
