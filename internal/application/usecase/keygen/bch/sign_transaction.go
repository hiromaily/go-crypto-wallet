package bch

import (
	"context"
	"fmt"

	apibtc "github.com/hiromaily/go-crypto-wallet/internal/application/ports/api/btc"
	file "github.com/hiromaily/go-crypto-wallet/internal/application/ports/file"
	repocold "github.com/hiromaily/go-crypto-wallet/internal/application/ports/repository/cold"
	keygenusecase "github.com/hiromaily/go-crypto-wallet/internal/application/usecase/keygen"
	usecaseshared "github.com/hiromaily/go-crypto-wallet/internal/application/usecase/keygen/shared"
	bchshared "github.com/hiromaily/go-crypto-wallet/internal/application/usecase/shared/bch"
	domainAccount "github.com/hiromaily/go-crypto-wallet/internal/domain/account"
	domainAddress "github.com/hiromaily/go-crypto-wallet/internal/domain/address"
	domainTx "github.com/hiromaily/go-crypto-wallet/internal/domain/transaction"
	"github.com/hiromaily/go-crypto-wallet/pkg/logger"
)

// signTransactionUseCase implements BCH-specific transaction signing.
// Unlike BTC, BCH uses Raw Transaction Hex format with SIGHASH_FORKID (0x41).
type signTransactionUseCase struct {
	bch             apibtc.BCHTransactionSigner
	accountKeyRepo  repocold.BTCAccountKeyRepositorier
	txFileRepo      file.TransactionFileRepositorier
	multisigAccount *domainAccount.MultisigConfig
}

// NewBCHSignTransactionUseCase creates a new BCH SignTransactionUseCase.
// BCH uses Raw TX Hex with SIGHASH_FORKID instead of PSBT.
func NewBCHSignTransactionUseCase(
	bch apibtc.BCHTransactionSigner,
	accountKeyRepo repocold.BTCAccountKeyRepositorier,
	txFileRepo file.TransactionFileRepositorier,
	multisigAccount *domainAccount.MultisigConfig,
) keygenusecase.SignTransactionUseCase {
	return &signTransactionUseCase{
		bch:             bch,
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
	// BCH uses .hex extension, so validate accordingly
	actionType, _, txID, signedCount, err := u.txFileRepo.ValidateFilePath(input.FilePath, domainTx.TxTypeUnsigned)
	if err != nil {
		return keygenusecase.SignTransactionOutput{}, err
	}

	// Read Raw TX Hex from file (BCH format)
	txContent, err := u.txFileRepo.ReadFile(input.FilePath)
	if err != nil {
		return keygenusecase.SignTransactionOutput{}, fmt.Errorf("fail to read tx file: %w", err)
	}

	// Parse the BCH raw tx content (hex + prevtx metadata)
	txHex, prevTxs, err := bchshared.ParseRawTxContent(txContent)
	if err != nil {
		return keygenusecase.SignTransactionOutput{}, fmt.Errorf("fail to parse raw tx content: %w", err)
	}

	// Sign raw transaction using SignRawTransactionWithKey
	// BCH signing uses SIGHASH_FORKID (0x41) internally
	signedHex, isSigned, err := u.sign(txHex, prevTxs, actionType)
	if err != nil {
		return keygenusecase.SignTransactionOutput{}, err
	}

	// Determine output file type based on signature completion
	txType := domainTx.TxTypeSigned
	if !isSigned {
		txType = domainTx.TxTypeUnsigned
		signedCount++
	}

	// Write signed transaction file
	// CreateFilePath returns base path without extension: {action}_{txID}_{type}_{count}_
	// WriteHexFile adds timestamp and .hex extension (following WritePSBTFile pattern)
	// Final format: payment_1_signed_0_1768819476518045000.hex
	basePath := u.txFileRepo.CreateFilePath(actionType, txType, txID, signedCount)

	var generatedFileName string
	// Workaround for BCH Node RPC bug in signrawtransactionwithkey.
	//
	// Issue: BCH Node incorrectly reports complete=true after adding the 2nd signature
	// to a 3-of-3 multisig transaction (should require all 3 signatures).
	//
	// Impact: When complete=true is reported prematurely, the transaction file is marked
	// as "signed" and prevTx metadata (TXID, vout, scriptPubKey, redeemScript, amount)
	// is omitted. Without this metadata, subsequent signers cannot add their signatures,
	// resulting in error: "Input not found or already spent"
	//
	// Solution: For multisig transactions, always include prevTx metadata regardless of
	// the complete flag. This allows Sign1 and Sign2 wallets to add their signatures.
	//
	// Detection: Use non-nil multisigAccount to determine if transaction is multisig.
	//
	// Reference: docs/chains/bch/README.md - "Known Issues and Workarounds"
	// Related: Issue #485, Issue #433
	isMultisig := u.multisigAccount != nil
	if bchshared.ShouldIncludePrevTxMetadata(isSigned, isMultisig) {
		// Include prevTx metadata for next signer (multisig or not yet fully signed)
		content := bchshared.FormatSignedTxContent(signedHex, prevTxs)
		generatedFileName, err = u.txFileRepo.WriteHexFile(basePath, content)
	} else {
		// For fully signed single-sig transactions, just write the hex
		generatedFileName, err = u.txFileRepo.WriteHexFile(basePath, signedHex)
	}
	if err != nil {
		return keygenusecase.SignTransactionOutput{}, fmt.Errorf("fail to write signed tx file: %w", err)
	}

	logger.Debug("signed BCH transaction",
		"action", actionType.String(),
		"txID", txID,
		"signedCount", signedCount,
		"isSigned", isSigned,
		"fileName", generatedFileName,
	)

	return keygenusecase.SignTransactionOutput{
		FilePath:      generatedFileName,
		IsDone:        isSigned,
		SignedCount:   1,
		UnsignedCount: 0,
	}, nil
}

// sign signs a BCH raw transaction using SignRawTransactionWithKey.
// BCH uses SIGHASH_FORKID (0x41) which is handled internally by the BCH node.
func (u *signTransactionUseCase) sign(
	txHex string,
	prevTxs []bchshared.PrevTx,
	actionType domainTx.ActionType,
) (string, bool, error) {
	// Infer sender account from action type
	senderAccount, err := usecaseshared.InferSenderAccount(actionType)
	if err != nil {
		return "", false, err
	}

	// Get WIF keys for this account
	wifs, err := u.getWIFsForAccount(senderAccount)
	if err != nil {
		return "", false, err
	}
	if len(wifs) == 0 {
		return "", false, fmt.Errorf("no WIF keys found for account: %s", senderAccount)
	}

	// Convert txHex to MsgTx
	msgTx, err := u.bch.ToMsgTx(txHex)
	if err != nil {
		return "", false, fmt.Errorf("fail to convert hex to MsgTx: %w", err)
	}

	// Convert prevTxs to DTO format
	dtoPrevTxs := bchshared.ConvertPrevTxsToDTO(prevTxs)

	// Sign raw transaction with WIF keys
	// BCH nodes handle SIGHASH_FORKID internally
	signedTx, isSigned, err := u.bch.SignRawTransactionWithKey(msgTx, wifs, dtoPrevTxs)
	if err != nil {
		return "", false, fmt.Errorf("fail to sign raw transaction: %w", err)
	}

	// Convert signed transaction back to hex
	signedHex, err := u.bch.ToHex(signedTx)
	if err != nil {
		return "", false, fmt.Errorf("fail to convert signed tx to hex: %w", err)
	}

	logger.Debug("BCH transaction signing completed",
		"action", actionType.String(),
		"sender_account", senderAccount.String(),
		"isSigned", isSigned,
	)

	return signedHex, isSigned, nil
}

// getWIFsForAccount gets WIF private keys for the specified account
func (u *signTransactionUseCase) getWIFsForAccount(
	senderAccount domainAccount.AccountType,
) ([]string, error) {
	// Get keys with various statuses
	// After export address command, status becomes AddrStatusAddressExported (3)
	// So we need to query for all relevant statuses
	statuses := []domainAddress.AddrStatus{
		domainAddress.AddrStatusHDKeyGenerated,           // 0
		domainAddress.AddrStatusPrivKeyImported,          // 1
		domainAddress.AddrStatusMultisigAddressGenerated, // 2 (for multisig)
		domainAddress.AddrStatusAddressExported,          // 3 (after export address)
	}

	wifs := make([]string, 0)
	for _, status := range statuses {
		accountKeys, err := u.accountKeyRepo.GetAllAddrStatus(senderAccount, status)
		if err != nil {
			return nil, fmt.Errorf("fail to get account keys for %s: %w", senderAccount.String(), err)
		}

		for _, key := range accountKeys {
			if key.WalletImportFormat != "" {
				wifs = append(wifs, key.WalletImportFormat)
			}
		}
	}

	return wifs, nil
}
