// Package btc provides BTC-specific keygen signing use cases.
//
// WARNING: This file is for BTC (Bitcoin) ONLY.
// If cointype is BCH (Bitcoin Cash), DO NOT modify this file.
// Use the BCH-specific implementation instead:
//
//	internal/application/usecase/keygen/bch/sign_transaction.go
//
// Key differences:
//   - BTC: Signs PSBT files with PSBTSigner interface, supports Taproot/Schnorr
//   - BCH: Signs raw transaction hex with SIGHASH_FORKID
package btc

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/btcsuite/btcd/btcutil"

	dtobtc "github.com/hiromaily/go-crypto-wallet/internal/application/dto/btc"
	apibtc "github.com/hiromaily/go-crypto-wallet/internal/application/ports/api/btc"
	file "github.com/hiromaily/go-crypto-wallet/internal/application/ports/file"
	repocold "github.com/hiromaily/go-crypto-wallet/internal/application/ports/repository/cold"
	keygenusecase "github.com/hiromaily/go-crypto-wallet/internal/application/usecase/keygen"
	usecaseshared "github.com/hiromaily/go-crypto-wallet/internal/application/usecase/shared"
	domainAccount "github.com/hiromaily/go-crypto-wallet/internal/domain/account"
	domainAddress "github.com/hiromaily/go-crypto-wallet/internal/domain/address"
	domainBitcoin "github.com/hiromaily/go-crypto-wallet/internal/domain/bitcoin"
	domainTx "github.com/hiromaily/go-crypto-wallet/internal/domain/transaction"
	infraKey "github.com/hiromaily/go-crypto-wallet/internal/infrastructure/wallet/key"
	"github.com/hiromaily/go-crypto-wallet/pkg/logger"
)

// signTxBTCClient defines the minimal interface needed for keygen PSBT signing.
// This follows the Interface Segregation Principle - depend only on methods actually used.
type signTxBTCClient interface {
	apibtc.ChainConfigProvider // GetChainConf
	apibtc.PSBTSigner          // SignPSBTWithKey
	apibtc.PSBTHandler         // WalletProcessPsbt, ParsePSBT
}

type signTransactionUseCase struct {
	btc             signTxBTCClient
	accountKeyRepo  repocold.BTCAccountKeyRepositorier
	txFileRepo      file.TransactionFileRepositorier
	multisigAccount *domainAccount.MultisigConfig
}

// NewSignTransactionUseCase creates a new SignTransactionUseCase for BTC keygen
func NewSignTransactionUseCase(
	btc apibtc.Bitcoiner,
	accountKeyRepo repocold.BTCAccountKeyRepositorier,
	txFileRepo file.TransactionFileRepositorier,
	multisigAccount *domainAccount.MultisigConfig,
) keygenusecase.SignTransactionUseCase {
	return &signTransactionUseCase{
		btc:             btc,
		accountKeyRepo:  accountKeyRepo,
		txFileRepo:      txFileRepo,
		multisigAccount: multisigAccount,
	}
}

func (u *signTransactionUseCase) Sign(
	ctx context.Context,
	input keygenusecase.SignTransactionInput,
) (keygenusecase.SignTransactionOutput, error) {
	// Get tx_deposit_id from tx file name
	//  if payment_5_unsigned_0_1534466246366489473.psbt, 5 is target
	actionType, _, txID, signedCount, err := u.txFileRepo.ValidateFilePath(input.FilePath, domainTx.TxTypeUnsigned)
	if err != nil {
		return keygenusecase.SignTransactionOutput{}, err
	}

	// Read PSBT from file
	psbtBase64, err := u.txFileRepo.ReadPSBTFile(input.FilePath)
	if err != nil {
		return keygenusecase.SignTransactionOutput{}, fmt.Errorf("fail to read PSBT file: %w", err)
	}

	// Sign PSBT (passing actionType to infer sender account)
	signedPSBT, isSigned, err := u.sign(psbtBase64, actionType)
	if err != nil {
		return keygenusecase.SignTransactionOutput{}, err
	}

	// If sign is not finished because of multisig, signedCount should be increment
	txType := domainTx.TxTypeSigned
	if !isSigned {
		txType = domainTx.TxTypeUnsigned
		signedCount++ // Increment for multisig partial signature
	}

	// Write signed PSBT file
	path := u.txFileRepo.CreateFilePath(actionType, txType, txID, signedCount)
	generatedFileName, err := u.txFileRepo.WritePSBTFile(path, signedPSBT)
	if err != nil {
		return keygenusecase.SignTransactionOutput{}, fmt.Errorf("fail to write signed PSBT file: %w", err)
	}

	logger.Debug("signed PSBT",
		"action", actionType.String(),
		"txID", txID,
		"signedCount", signedCount,
		"isSigned", isSigned,
		"fileName", generatedFileName,
	)

	return keygenusecase.SignTransactionOutput{
		FilePath:      generatedFileName,
		IsDone:        isSigned,
		SignedCount:   1, // BTC signs one transaction at a time
		UnsignedCount: 0, // BTC doesn't track unsigned separately in this interface
	}, nil
}

// sign signs a PSBT using offline signing with btcd library (no Bitcoin Core RPC).
// This function supports all Bitcoin address types including Taproot:
//   - Legacy (P2PKH): ECDSA signature
//   - SegWit (P2WPKH, P2SH-SegWit): ECDSA signature with witness data
//   - Taproot (P2TR): Schnorr signature (BIP340) with witness data
//
// The signature algorithm is automatically selected based on the PSBT input's scriptPubKey type.
// For Taproot addresses, Schnorr signatures (BIP340) are used.
//
// Transaction flow:
//   - [actionType:deposit]  [from] client [to] deposit (not multisig addr)
//   - [actionType:payment]  [from] payment [to] unknown (multisig addr)
//   - [actionType:transfer] [from] account [to] account (multisig addr)
//
// Note: This operates OFFLINE - no Bitcoin Core RPC required.
func (u *signTransactionUseCase) sign(
	psbtBase64 string,
	actionType domainTx.ActionType,
) (string, bool, error) {
	// Infer sender account from action type
	// This is a simplified approach since PSBT doesn't store the account concept
	senderAccount, err := usecaseshared.InferSenderAccount(actionType)
	if err != nil {
		return "", false, err
	}

	// Sign PSBT with keys from sender account
	signedPSBT, isSigned, err := u.signWithAccount(psbtBase64, senderAccount)
	if err != nil {
		return "", false, err
	}

	logger.Debug("PSBT signing completed",
		"action", actionType.String(),
		"sender_account", senderAccount.String(),
		"isSigned", isSigned,
	)

	return signedPSBT, isSigned, nil
}

// signWithAccount signs a PSBT with keys from the specified account.
// This function supports three signing approaches:
//
// 1. Descriptor-based derivation (account xpriv available):
//   - Parse PSBT to extract address index
//   - Derive child WIF at correct index from account xpriv
//   - Uses SignPSBTWithKey with derived WIF
//   - Works completely offline
//
// 2. Legacy offline signing (WIF keys in database, no xpriv):
//   - For legacy workflow where keys are exported to database
//   - Uses SignPSBTWithKey with all WIF private keys
//   - Works completely offline
//
// 3. RPC-based signing (descriptor wallets):
//   - For descriptor-based workflow where keys are managed by Bitcoin Core
//   - Uses WalletProcessPsbt RPC to sign with Bitcoin Core's wallet
//   - Requires connection to Bitcoin Core
//
// The function automatically detects which approach to use based on available keys and xpriv.
func (u *signTransactionUseCase) signWithAccount(
	psbtBase64 string,
	senderAccount domainAccount.AccountType,
) (string, bool, error) {
	// Try to get keys for this account, checking multiple statuses
	// For descriptor-based workflows: keys have AddrStatusHDKeyGenerated (don't need WIF import)
	// For legacy workflows: keys have AddrStatusPrivKeyImported (after WIF import to Bitcoin Core)
	statuses := []domainAddress.AddrStatus{
		domainAddress.AddrStatusHDKeyGenerated,
		domainAddress.AddrStatusPrivKeyImported,
	}

	var accountKeys []*domainBitcoin.BtcAccountKey
	var err error
	for _, status := range statuses {
		accountKeys, err = u.accountKeyRepo.GetAllAddrStatus(senderAccount, status)
		if err != nil {
			return "", false, fmt.Errorf("fail to get account keys for %s: %w", senderAccount.String(), err)
		}
		if len(accountKeys) > 0 {
			break
		}
	}

	// Check if we have account xpriv (descriptor-based workflow)
	// We only need to check one key since all keys for an account share the same xpriv
	hasXpriv := len(accountKeys) > 0 &&
		accountKeys[0].AccountExtendedPrivkey != nil &&
		*accountKeys[0].AccountExtendedPrivkey != ""
	logger.Debug("checking for account xpriv",
		"account", senderAccount.String(),
		"key_count", len(accountKeys),
		"has_xpriv", hasXpriv,
	)
	if len(accountKeys) > 0 &&
		accountKeys[0].AccountExtendedPrivkey != nil &&
		*accountKeys[0].AccountExtendedPrivkey != "" {
		// Approach 1: Descriptor-based derivation
		logger.Debug("using descriptor-based key derivation with account xpriv",
			"account", senderAccount.String(),
			"key_count", len(accountKeys),
		)

		// Derive WIFs for signing based on PSBT address indices (handles multi-input transactions)
		wifs, err := u.deriveWIFsForPSBT(psbtBase64, accountKeys[0])
		if err != nil {
			return "", false, fmt.Errorf("fail to derive WIFs for PSBT signing: %w", err)
		}

		// Sign PSBT with derived WIFs
		signedPSBT, isSigned, err := u.btc.SignPSBTWithKey(psbtBase64, wifs)
		if err != nil {
			return "", false, fmt.Errorf("fail to sign PSBT with derived keys: %w", err)
		}

		return signedPSBT, isSigned, nil
	}

	// Check if we have WIF keys in the database (legacy workflow)
	wifs := make([]string, 0)
	if len(accountKeys) > 0 {
		// Extract WIFs from account keys
		for _, key := range accountKeys {
			if key.WalletImportFormat != "" {
				wifs = append(wifs, key.WalletImportFormat)
			}
		}
	}

	// Determine signing approach based on available keys
	if len(wifs) > 0 {
		// Approach 2: Legacy offline signing with WIF keys
		logger.Debug("using legacy offline signing with WIF keys",
			"account", senderAccount.String(),
			"key_count", len(accountKeys),
			"wif_count", len(wifs),
		)

		// Sign PSBT with all WIFs - btcd will automatically use only matching keys
		signedPSBT, isSigned, err := u.btc.SignPSBTWithKey(psbtBase64, wifs)
		if err != nil {
			return "", false, fmt.Errorf("fail to sign PSBT with account keys: %w", err)
		}

		return signedPSBT, isSigned, nil
	}

	// Approach 3: RPC-based signing with descriptor wallet
	logger.Info("no WIF keys found, using descriptor-based signing via RPC",
		"account", senderAccount.String(),
	)

	// Sign PSBT using Bitcoin Core's wallet (descriptor-based)
	signedPSBT, isSigned, err := u.btc.WalletProcessPsbt(psbtBase64, true)
	if err != nil {
		return "", false, fmt.Errorf("fail to sign PSBT with walletprocesspsbt RPC: %w", err)
	}

	logger.Debug("descriptor-based signing completed",
		"account", senderAccount.String(),
		"complete", isSigned,
	)

	return signedPSBT, isSigned, nil
}

// deriveWIFsForPSBT derives all WIF (Wallet Import Format) private keys needed for signing a PSBT.
// This handles multi-input transactions (e.g., coin consolidation) by deriving keys for each input.
//
// For descriptor-based workflows (#320 fix):
//   - If accountExtendedPrivkey is available, parses PSBT to extract address indices from all inputs
//   - Derives child private keys at the correct indices from account-level xpriv
//   - Converts to WIF format for signing
//   - Returns unique WIFs (deduplicates if multiple inputs use the same address)
//
// For legacy workflows:
//   - Returns the stored WIF directly (static key at index 0)
//
// This ensures signatures match the descriptor-derived public keys regardless of address indices.
func (u *signTransactionUseCase) deriveWIFsForPSBT(
	psbtBase64 string,
	accountKey *domainBitcoin.BtcAccountKey,
) ([]string, error) {
	// Legacy workflow: Use stored WIF directly if no account xpriv available
	if accountKey.AccountExtendedPrivkey == nil || *accountKey.AccountExtendedPrivkey == "" {
		logger.Debug("using stored WIF (legacy workflow, no account xpriv)")
		return []string{accountKey.WalletImportFormat}, nil
	}

	// Descriptor workflow: Parse PSBT to extract address indices for all inputs
	parsed, err := u.btc.ParsePSBT(psbtBase64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse PSBT: %w", err)
	}

	if len(parsed.Inputs) == 0 {
		return nil, errors.New("PSBT has no inputs")
	}

	// Use map to deduplicate WIFs (multiple inputs may use same address)
	wifKeys := make(map[string]struct{})

	// Process each input to derive its required WIF
	for i, input := range parsed.Inputs {
		if len(input.BIP32Derivation) == 0 {
			// No BIP32 derivation info - derive fallback keys for P2TR
			u.deriveFallbackWIFsForP2TR(accountKey, i, wifKeys)
			continue
		}

		// Derive WIF from BIP32 derivation path
		if err := u.deriveWIFFromBIP32Path(accountKey, input, i, wifKeys); err != nil {
			return nil, err
		}
	}

	// If no WIFs could be derived, fall back to stored WIF
	if len(wifKeys) == 0 {
		logger.Warn("could not derive any specific WIFs for PSBT inputs, falling back to stored WIF")
		return []string{accountKey.WalletImportFormat}, nil
	}

	// Convert map keys to slice
	derivedWIFs := make([]string, 0, len(wifKeys))
	for wif := range wifKeys {
		derivedWIFs = append(derivedWIFs, wif)
	}

	logger.Debug("derived WIFs for PSBT inputs",
		"total_inputs", len(parsed.Inputs),
		"unique_wifs", len(derivedWIFs))

	return derivedWIFs, nil
}

// deriveFallbackWIFsForP2TR derives WIFs at multiple indices for P2TR inputs without BIP32 derivation.
// P2TR (Taproot) inputs don't include BIP32 derivation info, so we try common indices.
func (u *signTransactionUseCase) deriveFallbackWIFsForP2TR(
	accountKey *domainBitcoin.BtcAccountKey,
	inputIndex int,
	wifKeys map[string]struct{},
) {
	logger.Warn("PSBT input has no BIP32 derivation (likely P2TR), deriving multiple keys", "input_index", inputIndex)

	// Derive keys at multiple indices (both receive and change paths)
	for change := uint32(0); change <= 1; change++ {
		for idx := uint32(0); idx < 10; idx++ {
			childKey, err := infraKey.DeriveChildPrivateKey(*accountKey.AccountExtendedPrivkey, change, idx)
			if err != nil {
				logger.Debug("failed to derive key", "change", change, "index", idx, "error", err)
				continue
			}
			privKey, err := childKey.ECPrivKey()
			if err != nil {
				continue
			}
			wif, err := btcutil.NewWIF(privKey, u.btc.GetChainConf(), true)
			if err != nil {
				continue
			}
			wifKeys[wif.String()] = struct{}{}
		}
	}
}

// deriveWIFFromBIP32Path parses BIP32 derivation path and derives the corresponding WIF.
func (u *signTransactionUseCase) deriveWIFFromBIP32Path(
	accountKey *domainBitcoin.BtcAccountKey,
	input dtobtc.ParsedPSBTInput,
	inputIndex int,
	wifKeys map[string]struct{},
) error {
	// Parse BIP32 derivation path to extract address index and change
	// Path format: m/purpose'/coin'/account'/change/addressIndex
	firstDeriv := input.BIP32Derivation[0]
	pathComponents := strings.Split(strings.TrimPrefix(firstDeriv.Path, "m/"), "/")
	if len(pathComponents) < 5 {
		logger.Warn("invalid BIP32 path format, skipping input",
			"path", firstDeriv.Path,
			"input_index", inputIndex)
		return nil
	}

	// Parse address index and change from path
	addressIndex, change, err := u.parseAddressAndChangeFromPath(pathComponents, firstDeriv.Path, inputIndex)
	if err != nil {
		// Logged as warning, not fatal - skip this input
		return nil
	}

	logger.Debug("deriving child key from account xpriv",
		"wallet_type", "keygen",
		"input_index", inputIndex,
		"address_index", addressIndex,
		"change", change,
		"derivation_path", firstDeriv.Path)

	// Derive child private key at the correct address index
	childKey, err := infraKey.DeriveChildPrivateKey(*accountKey.AccountExtendedPrivkey, change, addressIndex)
	if err != nil {
		return fmt.Errorf("failed to derive child key for input %d at index %d: %w", inputIndex, addressIndex, err)
	}

	// Extract private key
	privKey, err := childKey.ECPrivKey()
	if err != nil {
		return fmt.Errorf("failed to get private key from child for input %d: %w", inputIndex, err)
	}

	// Get public key for verification logging
	pubKey, err := childKey.ECPubKey()
	if err != nil {
		return fmt.Errorf("failed to get public key from child for input %d: %w", inputIndex, err)
	}

	// Convert to WIF (compressed format)
	wif, err := btcutil.NewWIF(privKey, u.btc.GetChainConf(), true)
	if err != nil {
		return fmt.Errorf("failed to create WIF from derived key for input %d: %w", inputIndex, err)
	}

	logger.Info("derived WIF for signing",
		"wallet_type", "keygen",
		"input_index", inputIndex,
		"address_index", addressIndex,
		"change", change,
		"pubkey_hex", hex.EncodeToString(pubKey.SerializeCompressed()))

	wifKeys[wif.String()] = struct{}{}
	return nil
}

// parseAddressAndChangeFromPath extracts address index and change from BIP32 path components.
func (*signTransactionUseCase) parseAddressAndChangeFromPath(
	pathComponents []string,
	fullPath string,
	inputIndex int,
) (addressIndex, change uint32, err error) {
	// Parse address index (last component)
	addressIndexStr := strings.TrimSuffix(pathComponents[len(pathComponents)-1], "'")
	addrIdx, err := strconv.ParseUint(addressIndexStr, 10, 32)
	if err != nil {
		logger.Warn("failed to parse address index from path, skipping input",
			"path", fullPath,
			"input_index", inputIndex,
			"error", err)
		return 0, 0, err
	}

	// Parse change index (second to last component)
	changeStr := strings.TrimSuffix(pathComponents[len(pathComponents)-2], "'")
	chgIdx, err := strconv.ParseUint(changeStr, 10, 32)
	if err != nil {
		logger.Warn("failed to parse change index from path, skipping input",
			"path", fullPath,
			"input_index", inputIndex,
			"error", err)
		return 0, 0, err
	}

	return uint32(addrIdx), uint32(chgIdx), nil
}
