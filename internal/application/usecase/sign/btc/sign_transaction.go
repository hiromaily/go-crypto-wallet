package btc

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/btcsuite/btcd/btcutil"

	portsBtc "github.com/hiromaily/go-crypto-wallet/internal/application/ports/btc"
	"github.com/hiromaily/go-crypto-wallet/internal/application/ports/persistence"
	portsStorage "github.com/hiromaily/go-crypto-wallet/internal/application/ports/storage"
	signusecase "github.com/hiromaily/go-crypto-wallet/internal/application/usecase/sign"
	domainAccount "github.com/hiromaily/go-crypto-wallet/internal/domain/account"
	domainAuth "github.com/hiromaily/go-crypto-wallet/internal/domain/auth"
	domainTx "github.com/hiromaily/go-crypto-wallet/internal/domain/transaction"
	domainWallet "github.com/hiromaily/go-crypto-wallet/internal/domain/wallet"
	infraKey "github.com/hiromaily/go-crypto-wallet/internal/infrastructure/wallet/key"
	"github.com/hiromaily/go-crypto-wallet/pkg/logger"
)

type signTransactionUseCase struct {
	btc             portsBtc.Bitcoiner
	accountKeyRepo  persistence.BTCAccountKeyRepositorier
	authKeyRepo     persistence.AuthAccountKeyRepositorier
	txFileRepo      portsStorage.TransactionFileRepositorier
	multisigAccount *domainAccount.MultisigConfig
	wtype           domainWallet.WalletType
	authType        domainAccount.AuthType
}

// NewSignTransactionUseCase creates a new SignTransactionUseCase for sign wallet
func NewSignTransactionUseCase(
	btcAPI portsBtc.Bitcoiner,
	accountKeyRepo persistence.BTCAccountKeyRepositorier,
	authKeyRepo persistence.AuthAccountKeyRepositorier,
	txFileRepo portsStorage.TransactionFileRepositorier,
	multisigAccount *domainAccount.MultisigConfig,
	wtype domainWallet.WalletType,
	authType domainAccount.AuthType,
) signusecase.SignTransactionUseCase {
	return &signTransactionUseCase{
		btc:             btcAPI,
		accountKeyRepo:  accountKeyRepo,
		authKeyRepo:     authKeyRepo,
		txFileRepo:      txFileRepo,
		multisigAccount: multisigAccount,
		wtype:           wtype,
		authType:        authType,
	}
}

func (u *signTransactionUseCase) Sign(
	ctx context.Context,
	input signusecase.SignTransactionInput,
) (signusecase.SignTransactionOutput, error) {
	// Get tx_deposit_id from tx file name
	//  if payment_5_unsigned_1_1534466246366489473.psbt, 5 is target
	actionType, _, txID, signedCount, err := u.txFileRepo.ValidateFilePath(input.FilePath, domainTx.TxTypeUnsigned)
	if err != nil {
		return signusecase.SignTransactionOutput{}, err
	}

	// Read PSBT from file
	psbtBase64, err := u.txFileRepo.ReadPSBTFile(input.FilePath)
	if err != nil {
		return signusecase.SignTransactionOutput{}, fmt.Errorf("fail to read PSBT file: %w", err)
	}

	// Sign PSBT (add second signature for multisig)
	signedPSBT, isSigned, err := u.sign(psbtBase64, actionType)
	if err != nil {
		return signusecase.SignTransactionOutput{}, err
	}

	// If sign is not finished because of multisig, signedCount should be increment
	txType := domainTx.TxTypeSigned
	if !isSigned {
		txType = domainTx.TxTypeUnsigned
		signedCount++ // Increment for additional signatures needed
	}

	// Write signed PSBT file
	path := u.txFileRepo.CreateFilePath(actionType, txType, txID, signedCount)
	generatedFileName, err := u.txFileRepo.WritePSBTFile(path, signedPSBT)
	if err != nil {
		return signusecase.SignTransactionOutput{}, fmt.Errorf("fail to write signed PSBT file: %w", err)
	}

	logger.Debug("signed PSBT",
		"action", actionType.String(),
		"txID", txID,
		"signedCount", signedCount,
		"isSigned", isSigned,
		"fileName", generatedFileName,
	)

	return signusecase.SignTransactionOutput{
		SignedData:   signedPSBT, // PSBT base64
		IsComplete:   isSigned,
		NextFilePath: generatedFileName,
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
// The Sign wallet acts as the second (or subsequent) signer in a multisig setup, using keys from
// the auth_account_key table. This adds signatures to a partially signed PSBT from Keygen wallet.
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
	// Sign wallet always signs multisig transactions
	// Add second signature to PSBT using auth key
	// Determine which multisig account to use based on the action type
	accountType := actionType.ToAccountType()
	signedPSBT, isSigned, err := u.signMultisigPSBT(psbtBase64, accountType)
	if err != nil {
		return "", false, err
	}

	logger.Debug("PSBT signing completed",
		"action", actionType.String(),
		"wallet_type", u.wtype.String(),
		"isSigned", isSigned,
	)

	return signedPSBT, isSigned, nil
}

// signMultisigPSBT adds the second signature to a partially signed PSBT for multisig transactions.
// The Sign wallet uses the auth key from auth_account_key table to add its signature.
// This function works offline using btcd library, supporting all address types:
//   - P2SH (Legacy multisig): ECDSA signature
//   - P2SH-SegWit: ECDSA signature with witness
//   - P2WSH (Native SegWit multisig): ECDSA signature with witness
//   - P2TR (Taproot): Schnorr signature (BIP340) or script path
//
// For 2-of-2 multisig, this signature typically completes the transaction.
// For 2-of-N multisig (N>2), the transaction is complete once 2 signatures are present.
//
// CRITICAL FIX (#320): For descriptor-based multisig, derives child keys at the correct
// address index from account-level extended private key, ensuring signatures match
// the descriptor-derived public keys.
func (u *signTransactionUseCase) signMultisigPSBT(
	psbtBase64 string,
	accountType domainAccount.AccountType,
) (string, bool, error) {
	// Get auth key from auth_account_key table (Sign wallet's key)
	// Using explicit authType and accountType for robust key selection
	// The accountType ensures we use keys derived at the correct BIP44 account index
	// to match the descriptor public keys (e.g., m/49'/1'/1' for payment account)
	authKey, err := u.authKeyRepo.GetByAccount(u.authType, accountType)
	if err != nil {
		return "", false, fmt.Errorf("fail to get auth key for authType %s, account %s: %w",
			u.authType, accountType, err)
	}

	logger.Debug("signing PSBT with auth key",
		"wallet_type", u.wtype.String(),
		"auth_type", u.authType.String(),
		"account_type", accountType.String(),
		"has_account_xpriv", authKey.AccountExtendedPrivkey != nil,
	)

	// Derive WIF for signing based on PSBT address index
	wif, err := u.deriveWIFForPSBT(psbtBase64, authKey)
	if err != nil {
		return "", false, fmt.Errorf("fail to derive WIF for PSBT signing: %w", err)
	}

	// Sign PSBT with Sign wallet's private key (offline, using btcd)
	// This adds the second signature to the partially signed PSBT
	signedPSBT, isSigned, err := u.btc.SignPSBTWithKey(psbtBase64, []string{wif})
	if err != nil {
		return "", false, fmt.Errorf("fail to sign PSBT with auth key: %w", err)
	}

	return signedPSBT, isSigned, nil
}

// deriveWIFForPSBT derives the appropriate WIF (Wallet Import Format) private key for signing a PSBT.
//
// For descriptor-based workflows (#320 fix):
//   - If accountExtendedPrivkey is available, parses PSBT to extract address index
//   - Derives child private key at the correct index from account-level xpriv
//   - Converts to WIF format for signing
//
// For legacy workflows:
//   - Returns the stored WIF directly (static key at index 0)
//
// This ensures signatures match the descriptor-derived public keys regardless of address index.
func (u *signTransactionUseCase) deriveWIFForPSBT(
	psbtBase64 string,
	authKey *domainAuth.AuthAccountKey,
) (string, error) {
	// Legacy workflow: Use stored WIF directly if no account xpriv available
	if authKey.AccountExtendedPrivkey == nil || *authKey.AccountExtendedPrivkey == "" {
		logger.Debug("using stored WIF (legacy workflow, no account xpriv)")
		return authKey.WalletImportFormat, nil
	}

	// Descriptor workflow: Parse PSBT to extract address index
	parsed, err := u.btc.ParsePSBT(psbtBase64)
	if err != nil {
		return "", fmt.Errorf("failed to parse PSBT: %w", err)
	}

	if len(parsed.Inputs) == 0 {
		return "", errors.New("PSBT has no inputs")
	}

	// Extract address index from first input's BIP32 derivation
	// All inputs should have same address index for single-address transactions
	if len(parsed.Inputs[0].BIP32Derivation) == 0 {
		// No BIP32 derivation info - fall back to stored WIF
		logger.Warn("PSBT input has no BIP32 derivation information, using stored WIF")
		return authKey.WalletImportFormat, nil
	}

	// Convert DTO BIP32Derivation to format expected by ExtractAddressIndexFromPSBTInput
	// We need to parse the Path string to get the Bip32Path []uint32
	firstDeriv := parsed.Inputs[0].BIP32Derivation[0]
	addressIndex, err := infraKey.ParseBIP32DerivationPath(firstDeriv.Path)
	if err != nil {
		return "", fmt.Errorf("failed to parse BIP32 derivation path %s: %w", firstDeriv.Path, err)
	}

	// Determine change index from path
	// Path format: m/purpose'/coin'/account'/change/addressIndex
	pathComponents := strings.Split(strings.TrimPrefix(firstDeriv.Path, "m/"), "/")
	if len(pathComponents) < 5 {
		return "", fmt.Errorf("invalid BIP32 path format: %s", firstDeriv.Path)
	}
	change, err := strconv.ParseUint(strings.TrimSuffix(pathComponents[len(pathComponents)-2], "'"), 10, 32)
	if err != nil {
		return "", fmt.Errorf("failed to parse change index from path %s: %w", firstDeriv.Path, err)
	}

	logger.Debug("deriving child key from account xpriv",
		"address_index", addressIndex,
		"change", change,
		"derivation_path", firstDeriv.Path)

	// Derive child private key at the correct address index
	childKey, err := infraKey.DeriveChildPrivateKey(*authKey.AccountExtendedPrivkey, uint32(change), addressIndex)
	if err != nil {
		return "", fmt.Errorf("failed to derive child key at index %d: %w", addressIndex, err)
	}

	// Extract private key
	privKey, err := childKey.ECPrivKey()
	if err != nil {
		return "", fmt.Errorf("failed to get private key from child: %w", err)
	}

	// Convert to WIF (compressed format)
	wif, err := btcutil.NewWIF(privKey, u.btc.GetChainConf(), true)
	if err != nil {
		return "", fmt.Errorf("failed to create WIF from derived key: %w", err)
	}

	logger.Debug("derived WIF for address index",
		"address_index", addressIndex,
		"change", change)

	return wif.String(), nil
}
