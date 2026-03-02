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
// │   - docs/chains/bch/README.md                                               │
// └─────────────────────────────────────────────────────────────────────────────┘
//
// PSBT is defined in BIP174 and provides a standard format for unsigned/partially
// signed transactions that can be passed between parties for multi-signature workflows.
//
// BCH Protocol Differences (DO NOT apply PSBT to BCH):
//   - BCH: Raw TX Hex + ECDSA + SIGHASH_FORKID (0x41)
//   - BTC: PSBT (BIP174) + ECDSA/Schnorr + standard sighash types
//
// File organization:
//   - psbt.go: Core PSBT types and creation (CreatePSBT, Parse, Validate, SignPSBTWithKey)
//   - psbt_sign.go: Signing functions (signInputWithKey, signTaprootInput, signSegWitInput, signLegacyInput)
//   - psbt_finalize.go: Finalization functions (FinalizePSBT, finalizeP2PKHInput, finalizeMultisigInput)
//   - psbt_utils.go: Utility functions (ExtractTransaction, IsPSBTComplete, GetPSBTFee, WalletProcessPsbt)
//   - psbt_bip32.go: BIP32 derivation functions (addBIP32DerivationForInput, parseDerivationPath)
//   - psbt_descriptor.go: Descriptor derivation functions (deriveRedeemScriptForAddress, deriveNativeSegWitAddress)
package btc

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"

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
	infraPrevTxs := fromDTOPreviousTxToLocal(prevTxs)
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

	return ToParsedPSBT(infraParsed)
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
		// Log detailed information about why signing failed
		logger.Error("PSBT signing failed - no signatures added",
			"inputCount", len(parsed.Packet.UnsignedTx.TxIn),
			"keyCount", len(privKeys))
		for i, input := range parsed.Packet.Inputs {
			if input.WitnessUtxo != nil {
				isTaproot := txscript.IsPayToTaproot(input.WitnessUtxo.PkScript)
				logger.Error("Input details",
					"input", i,
					"isTaproot", isTaproot,
					"pkScript", hex.EncodeToString(input.WitnessUtxo.PkScript))
			}
		}
		for j, privKey := range privKeys {
			pubKey := privKey.PrivKey.PubKey()
			logger.Error("Key details",
				"keyIndex", j,
				"compressedPubKey", hex.EncodeToString(pubKey.SerializeCompressed()),
				"xOnlyPubKey", hex.EncodeToString(schnorr.SerializePubKey(pubKey)))
		}
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
