package btc

import (
	"context"
	"fmt"

	portsBtc "github.com/hiromaily/go-crypto-wallet/internal/application/ports/btc"
	"github.com/hiromaily/go-crypto-wallet/internal/application/ports/persistence"
	portsStorage "github.com/hiromaily/go-crypto-wallet/internal/application/ports/storage"
	keygenusecase "github.com/hiromaily/go-crypto-wallet/internal/application/usecase/keygen"
	domainAccount "github.com/hiromaily/go-crypto-wallet/internal/domain/account"
	domainAddress "github.com/hiromaily/go-crypto-wallet/internal/domain/address"
	domainTx "github.com/hiromaily/go-crypto-wallet/internal/domain/transaction"
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
// This function supports two signing approaches:
//
// 1. Offline signing (WIF keys in database):
//   - For legacy workflow where keys are exported to database
//   - Uses SignPSBTWithKey with WIF private keys
//   - Works completely offline
//
// 2. RPC-based signing (descriptor wallets):
//   - For descriptor-based workflow where keys are managed by Bitcoin Core
//   - Uses WalletProcessPsbt RPC to sign with Bitcoin Core's wallet
//   - Requires connection to Bitcoin Core
//
// The function automatically detects which approach to use:
// - If WIF keys are found in database → offline signing
// - If no WIF keys found → RPC-based signing
func (u *signTransactionUseCase) signWithAccount(
	psbtBase64 string,
	senderAccount domainAccount.AccountType,
) (string, bool, error) {
	// Try to get exported keys for this account
	// Using AddrStatusAddressExported ensures keys are ready and have been exported to watch wallet
	accountKeys, err := u.accountKeyRepo.GetAllAddrStatus(
		senderAccount,
		domainAddress.AddrStatusAddressExported,
	)
	if err != nil {
		return "", false, fmt.Errorf("fail to get account keys for %s: %w", senderAccount.String(), err)
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
		// Approach 1: Offline signing with WIF keys
		logger.Debug("using offline signing with WIF keys",
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

	// Approach 2: RPC-based signing with descriptor wallet
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
