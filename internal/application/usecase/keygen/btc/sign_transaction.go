package btc

import (
	"context"
	"fmt"

	keygenusecase "github.com/hiromaily/go-crypto-wallet/internal/application/usecase/keygen"
	domainAccount "github.com/hiromaily/go-crypto-wallet/internal/domain/account"
	domainTx "github.com/hiromaily/go-crypto-wallet/internal/domain/transaction"
	"github.com/hiromaily/go-crypto-wallet/internal/infrastructure/api/bitcoin"
	"github.com/hiromaily/go-crypto-wallet/internal/infrastructure/api/bitcoin/btc"
	"github.com/hiromaily/go-crypto-wallet/internal/infrastructure/config/account"
	"github.com/hiromaily/go-crypto-wallet/internal/infrastructure/repository/cold"
	"github.com/hiromaily/go-crypto-wallet/internal/infrastructure/storage/file"
	"github.com/hiromaily/go-crypto-wallet/internal/infrastructure/storage/file/address"
	"github.com/hiromaily/go-crypto-wallet/pkg/logger"
)

type signTransactionUseCase struct {
	btc             bitcoin.Bitcoiner
	accountKeyRepo  cold.AccountKeyRepositorier
	txFileRepo      file.TransactionFileRepositorier
	multisigAccount account.MultisigAccounter
}

// NewSignTransactionUseCase creates a new SignTransactionUseCase for BTC keygen
func NewSignTransactionUseCase(
	btc bitcoin.Bitcoiner,
	accountKeyRepo cold.AccountKeyRepositorier,
	txFileRepo file.TransactionFileRepositorier,
	multisigAccount account.MultisigAccounter,
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

	// Parse PSBT to extract metadata
	parsedPSBT, err := u.btc.ParsePSBT(psbtBase64)
	if err != nil {
		return keygenusecase.SignTransactionOutput{}, fmt.Errorf("fail to parse PSBT: %w", err)
	}

	// Sign PSBT (passing actionType to infer sender account)
	signedPSBT, isSigned, err := u.sign(psbtBase64, parsedPSBT, actionType)
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
	parsedPSBT *btc.ParsedPSBT,
	actionType domainTx.ActionType,
) (string, bool, error) {
	// Infer sender account from action type
	// This is a simplified approach since PSBT doesn't store the account concept
	senderAccount := inferSenderAccount(actionType)

	// Sign PSBT with keys from sender account
	signedPSBT, isSigned, err := u.signWithAccount(psbtBase64, senderAccount)
	if err != nil {
		return "", false, err
	}

	logger.Debug("PSBT signing completed",
		"action", actionType.String(),
		"sender_account", senderAccount.String(),
		"isSigned", isSigned,
		"inputs", len(parsedPSBT.Packet.Inputs),
	)

	return signedPSBT, isSigned, nil
}

// inferSenderAccount infers the sender account from the transaction action type.
// This is a pragmatic approach since PSBT doesn't encode the account concept.
func inferSenderAccount(actionType domainTx.ActionType) domainAccount.AccountType {
	switch actionType {
	case domainTx.ActionTypeDeposit:
		return domainAccount.AccountTypeClient
	case domainTx.ActionTypePayment:
		return domainAccount.AccountTypePayment
	case domainTx.ActionTypeTransfer:
		// Transfer could be from various accounts
		// Default to payment for now
		return domainAccount.AccountTypePayment
	default:
		return domainAccount.AccountTypeClient
	}
}

// signWithAccount signs a PSBT with keys from the specified account.
// This is a simplified MVP approach that works for both single-sig and multisig:
// - For single-sig: Signs completely if the key matches
// - For multisig: Adds first signature (Keygen wallet signature)
//
// The PSBT signing operation will only apply signatures where keys match input requirements.
// This approach works offline without needing to extract addresses from PSBT scriptPubKeys.
//
// Implementation:
// 1. Get all exported keys for the sender account
// 2. Extract WIFs from retrieved keys
// 3. Pass all WIFs to SignPSBTWithKey - btcd will use only matching keys
func (u *signTransactionUseCase) signWithAccount(
	psbtBase64 string,
	senderAccount domainAccount.AccountType,
) (string, bool, error) {
	// Get all exported keys for this account
	// Using AddrStatusAddressExported ensures keys are ready and have been exported to watch wallet
	accountKeys, err := u.accountKeyRepo.GetAllAddrStatus(
		senderAccount,
		address.AddrStatusAddressExported,
	)
	if err != nil {
		return "", false, fmt.Errorf("fail to get account keys for %s: %w", senderAccount.String(), err)
	}

	if len(accountKeys) == 0 {
		return "", false, fmt.Errorf("no exported keys found for account %s", senderAccount.String())
	}

	// Extract WIFs from account keys
	wifs := make([]string, 0, len(accountKeys))
	for _, key := range accountKeys {
		if key.WalletImportFormat != "" {
			wifs = append(wifs, key.WalletImportFormat)
		}
	}

	if len(wifs) == 0 {
		return "", false, fmt.Errorf("no valid WIFs found for account %s", senderAccount.String())
	}

	logger.Debug("signing PSBT with account keys",
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
