// Package btc provides Bitcoin (BTC) infrastructure implementations.
//
// This file contains PSBT signing functions.
// See psbt.go for the main PSBT documentation.
package btc

import (
	"bytes"
	"encoding/hex"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcec/v2/ecdsa"
	"github.com/btcsuite/btcd/btcec/v2/schnorr"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/btcutil/psbt"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"

	"github.com/hiromaily/go-crypto-wallet/pkg/logger"
)

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
// For key-path spending (BIP86), the private key must be tweaked before signing.
// Note: Since btcd's psbt package doesn't have TaprootKeyPathSig field,
// we add the signature directly to FinalScriptWitness for single-sig scenarios.
func (*Bitcoin) signTaprootInput(
	_ *psbt.Updater,
	msgTx *wire.MsgTx,
	inputIndex int,
	psbtInput *psbt.PInput,
	privKey *btcutil.WIF,
	prevOutputFetcher *txscript.MultiPrevOutFetcher,
) bool {
	// Extract the x-only output key from the scriptPubKey
	// P2TR scriptPubKey format: OP_1 (0x51) + OP_DATA_32 (0x20) + 32-byte x-only pubkey
	pkScript := psbtInput.WitnessUtxo.PkScript
	if len(pkScript) != 34 {
		logger.Warn("Invalid P2TR scriptPubKey length", "input", inputIndex, "length", len(pkScript))
		return false
	}
	outputKeyBytes := pkScript[2:34]

	// Get the internal public key (before tweaking)
	internalPubKey := privKey.PrivKey.PubKey()
	internalPubKeyXOnly := schnorr.SerializePubKey(internalPubKey)

	logger.Warn("Taproot signing attempt",
		"input", inputIndex,
		"outputKeyInScript", hex.EncodeToString(outputKeyBytes),
		"internalPubKeyCompressed", hex.EncodeToString(internalPubKey.SerializeCompressed()),
		"internalPubKeyXOnly", hex.EncodeToString(internalPubKeyXOnly))

	// For BIP86 key-path spending, tweak the private key with empty merkle root
	// The output key = internal_key + H(internal_key || empty) * G
	tweakedPrivKey := txscript.TweakTaprootPrivKey(*privKey.PrivKey, nil)

	// Verify that the tweaked public key matches the output key in scriptPubKey
	tweakedPubKey := tweakedPrivKey.PubKey()
	tweakedPubKeyBytes := schnorr.SerializePubKey(tweakedPubKey)

	logger.Warn("Taproot key comparison",
		"input", inputIndex,
		"tweakedPubKey", hex.EncodeToString(tweakedPubKeyBytes),
		"outputKey", hex.EncodeToString(outputKeyBytes),
		"match", bytes.Equal(tweakedPubKeyBytes, outputKeyBytes))

	if !bytes.Equal(tweakedPubKeyBytes, outputKeyBytes) {
		// Also try matching without tweaking (in case output key IS the internal key)
		if bytes.Equal(internalPubKeyXOnly, outputKeyBytes) {
			logger.Info("Output key matches internal key without tweaking - using untweaked key",
				"input", inputIndex)
			// Sign without tweaking (unusual but possible for some descriptor setups)
			return signTaprootWithKey(msgTx, inputIndex, psbtInput, privKey.PrivKey, prevOutputFetcher)
		}
		logger.Debug("Tweaked pubkey does not match output key, key not applicable for this input",
			"input", inputIndex)
		return false
	}

	// Sign with the tweaked private key
	return signTaprootWithKey(msgTx, inputIndex, psbtInput, tweakedPrivKey, prevOutputFetcher)
}

// signTaprootWithKey signs a Taproot input with a given private key (already tweaked if needed)
func signTaprootWithKey(
	msgTx *wire.MsgTx,
	inputIndex int,
	psbtInput *psbt.PInput,
	privKey *btcec.PrivateKey,
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
	signature, err := schnorr.Sign(privKey, hash)
	if err != nil {
		logger.Warn("Failed to create Schnorr signature", "input", inputIndex, "error", err)
		return false
	}

	// For Taproot with SIGHASH_DEFAULT, signature is just 64 bytes (no sighash type appended)
	sigBytes := signature.Serialize()

	// For P2TR key-path spending, the witness is just the signature (64 bytes)
	// Since btcd's psbt package doesn't have TaprootKeyPathSig field,
	// we serialize the witness directly to FinalScriptWitness format
	// Witness format: [varint count=1][varint len=64][64-byte signature]
	witnessData := make([]byte, 0, 66)
	witnessData = append(witnessData, 0x01)        // witness stack count = 1
	witnessData = append(witnessData, 0x40)        // signature length = 64 (0x40)
	witnessData = append(witnessData, sigBytes...) // 64-byte Schnorr signature

	psbtInput.FinalScriptWitness = witnessData

	logger.Debug("Added Taproot key-path signature to input",
		"input", inputIndex,
		"sigLen", len(sigBytes),
		"witnessLen", len(witnessData))
	return true
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
		return nil, newPSBTError(inputIndex, "invalid P2WPKH redeemScript length: expected %d, got %d",
			p2wpkhRedeemScriptLen, len(redeemScript))
	}

	// Validate format: must be OP_0 <20-byte-hash>
	// Format: [OP_0 (0x00), length byte (0x14 = 20 decimal), 20-byte pubkey hash]
	if redeemScript[0] != txscript.OP_0 {
		return nil, newPSBTError(inputIndex,
			"invalid P2WPKH redeemScript format: must start with OP_0 (0x00), got 0x%02x", redeemScript[0])
	}

	// Validate length byte (0x14 = 20 bytes for pubkey hash)
	const pubKeyHashLength = 20
	if redeemScript[1] != pubKeyHashLength {
		return nil, newPSBTError(inputIndex,
			"invalid P2WPKH redeemScript format: invalid length byte, expected 0x14, got 0x%02x", redeemScript[1])
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
		return nil, newPSBTError(inputIndex, "failed to build P2PKH scriptCode: %v", err)
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

// newPSBTError creates a formatted PSBT error
func newPSBTError(inputIndex int, format string, args ...any) error {
	return &psbtError{inputIndex: inputIndex, message: format, args: args}
}

// psbtError represents a PSBT-related error with input context
type psbtError struct {
	inputIndex int
	message    string
	args       []any
}

func (e *psbtError) Error() string {
	return "input " + string(rune(e.inputIndex+'0')) + ": " + e.message
}
