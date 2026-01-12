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
	keygenusecase "github.com/hiromaily/go-crypto-wallet/internal/application/usecase/keygen"
	domainAccount "github.com/hiromaily/go-crypto-wallet/internal/domain/account"
	domainAddress "github.com/hiromaily/go-crypto-wallet/internal/domain/address"
	domainBitcoin "github.com/hiromaily/go-crypto-wallet/internal/domain/bitcoin"
	domainTx "github.com/hiromaily/go-crypto-wallet/internal/domain/transaction"
	infraKey "github.com/hiromaily/go-crypto-wallet/internal/infrastructure/wallet/key"
	"github.com/hiromaily/go-crypto-wallet/pkg/logger"
)

type signTransactionUseCase struct {
	btc             portsBtc.Bitcoiner
	accountKeyRepo  persistence.BTCAccountKeyRepositorier
	txFileRepo      portsStorage.TransactionFileRepositorier
	multisigAccount *domainAccount.MultisigConfig
}

// NewSignTransactionUseCase creates a new SignTransactionUseCase for BTC keygen
func NewSignTransactionUseCase(
	btc portsBtc.Bitcoiner,
	accountKeyRepo persistence.BTCAccountKeyRepositorier,
	txFileRepo portsStorage.TransactionFileRepositorier,
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
	senderAccount, err := inferSenderAccount(actionType)
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

// inferSenderAccount infers the sender account from the transaction action type.
// This is a pragmatic approach since PSBT doesn't encode the account concept.
func inferSenderAccount(actionType domainTx.ActionType) (domainAccount.AccountType, error) {
	switch actionType {
	case domainTx.ActionTypeDeposit:
		return domainAccount.AccountTypeClient, nil
	case domainTx.ActionTypePayment:
		return domainAccount.AccountTypePayment, nil
	case domainTx.ActionTypeTransfer:
		// Transfer could be from various accounts
		// Default to payment for now
		return domainAccount.AccountTypePayment, nil
	default:
		return "", fmt.Errorf("unsupported action type: %s", actionType)
	}
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
	// Try to get generated keys for this account
	// For keygen wallet, keys are available after generation and import (AddrStatusPrivKeyImported)
	// We don't need to wait for export to watch wallet to use them for signing
	accountKeys, err := u.accountKeyRepo.GetAllAddrStatus(
		senderAccount,
		domainAddress.AddrStatusPrivKeyImported,
	)
	if err != nil {
		return "", false, fmt.Errorf("fail to get account keys for %s: %w", senderAccount.String(), err)
	}

	// Check if we have account xpriv (descriptor-based workflow)
	// We only need to check one key since all keys for an account share the same xpriv
	logger.Debug("checking for account xpriv",
		"account", senderAccount.String(),
		"key_count", len(accountKeys),
		"has_xpriv", len(accountKeys) > 0 && accountKeys[0].AccountExtendedPrivkey != nil && *accountKeys[0].AccountExtendedPrivkey != "",
	)
	if len(accountKeys) > 0 &&
		accountKeys[0].AccountExtendedPrivkey != nil &&
		*accountKeys[0].AccountExtendedPrivkey != "" {
		// Approach 1: Descriptor-based derivation
		logger.Debug("using descriptor-based key derivation with account xpriv",
			"account", senderAccount.String(),
			"key_count", len(accountKeys),
		)

		// Derive WIF for signing based on PSBT address index
		wif, err := u.deriveWIFForPSBT(psbtBase64, accountKeys[0])
		if err != nil {
			return "", false, fmt.Errorf("fail to derive WIF for PSBT signing: %w", err)
		}

		// Sign PSBT with derived WIF
		signedPSBT, isSigned, err := u.btc.SignPSBTWithKey(psbtBase64, []string{wif})
		if err != nil {
			return "", false, fmt.Errorf("fail to sign PSBT with derived key: %w", err)
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
	accountKey *domainBitcoin.BtcAccountKey,
) (string, error) {
	// Legacy workflow: Use stored WIF directly if no account xpriv available
	if accountKey.AccountExtendedPrivkey == nil || *accountKey.AccountExtendedPrivkey == "" {
		logger.Debug("using stored WIF (legacy workflow, no account xpriv)")
		return accountKey.WalletImportFormat, nil
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
		return accountKey.WalletImportFormat, nil
	}

	// Parse BIP32 derivation path to extract address index and change
	// Path format: m/purpose'/coin'/account'/change/addressIndex
	firstDeriv := parsed.Inputs[0].BIP32Derivation[0]
	pathComponents := strings.Split(strings.TrimPrefix(firstDeriv.Path, "m/"), "/")
	if len(pathComponents) < 5 {
		return "", fmt.Errorf("invalid BIP32 path format: %s", firstDeriv.Path)
	}

	// Parse address index (last component)
	addressIndexStr := strings.TrimSuffix(pathComponents[len(pathComponents)-1], "'")
	addrIdx, err := strconv.ParseUint(addressIndexStr, 10, 32)
	if err != nil {
		return "", fmt.Errorf("failed to parse address index from path %s: %w", firstDeriv.Path, err)
	}
	addressIndex := uint32(addrIdx)

	// Parse change index (second to last component)
	changeStr := strings.TrimSuffix(pathComponents[len(pathComponents)-2], "'")
	chgIdx, err := strconv.ParseUint(changeStr, 10, 32)
	if err != nil {
		return "", fmt.Errorf("failed to parse change index from path %s: %w", firstDeriv.Path, err)
	}
	change := uint32(chgIdx)

	logger.Debug("deriving child key from account xpriv",
		"address_index", addressIndex,
		"change", change,
		"derivation_path", firstDeriv.Path)

	// Derive child private key at the correct address index
	childKey, err := infraKey.DeriveChildPrivateKey(*accountKey.AccountExtendedPrivkey, change, addressIndex)
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
