package xrp

import (
	"context"
	"errors"
	"fmt"

	apixrp "github.com/hiromaily/go-crypto-wallet/internal/application/ports/api/xrp"
	file "github.com/hiromaily/go-crypto-wallet/internal/application/ports/file"
	repocold "github.com/hiromaily/go-crypto-wallet/internal/application/ports/repository/cold"
	signusecase "github.com/hiromaily/go-crypto-wallet/internal/application/usecase/sign"
	domainAccount "github.com/hiromaily/go-crypto-wallet/internal/domain/account"
	domainTx "github.com/hiromaily/go-crypto-wallet/internal/domain/transaction"
	domainWallet "github.com/hiromaily/go-crypto-wallet/internal/domain/wallet"
	"github.com/hiromaily/go-crypto-wallet/pkg/logger"
)

// xrpSignClient defines the interface for XRP operations needed by signTransactionUseCase.
// This follows the Interface Segregation Principle - depend only on what you need.
type xrpSignClient interface {
	apixrp.TransactionSigner
}

type signTransactionUseCase struct {
	xrp               xrpSignClient
	xrpAccountKeyRepo repocold.XRPAccountKeyRepositorier
	txFileRepo        file.TransactionFileRepositorier
	wtype             domainWallet.WalletType
}

// NewSignTransactionUseCase creates a new SignTransactionUseCase for sign wallet.
// The xrpAPI parameter accepts any type that implements xrpSignClient (TransactionSigner).
// Typically, apixrp.XRPer is passed which implements all required methods.
func NewSignTransactionUseCase(
	xrpAPI xrpSignClient,
	xrpAccountKeyRepo repocold.XRPAccountKeyRepositorier,
	txFileRepo file.TransactionFileRepositorier,
	wtype domainWallet.WalletType,
) signusecase.SignTransactionUseCase {
	return &signTransactionUseCase{
		xrp:               xrpAPI,
		xrpAccountKeyRepo: xrpAccountKeyRepo,
		txFileRepo:        txFileRepo,
		wtype:             wtype,
	}
}

func (u *signTransactionUseCase) Sign(
	ctx context.Context,
	input signusecase.SignTransactionInput,
) (signusecase.SignTransactionOutput, error) {
	// Step 1: Validate file path and extract metadata
	actionType, _, txID, signedCount, err := u.txFileRepo.ValidateFilePath(input.FilePath, domainTx.TxTypeUnsigned)
	if err != nil {
		return signusecase.SignTransactionOutput{}, fmt.Errorf("failed to validate file path: %w", err)
	}

	// Step 2: Read JSON transaction file
	txFile, err := u.txFileRepo.ReadXRPJSONFile(input.FilePath)
	if err != nil {
		return signusecase.SignTransactionOutput{}, fmt.Errorf("failed to read JSON transaction file: %w", err)
	}

	// Step 3: Validate transaction file has entries
	if len(txFile.Transactions) == 0 {
		return signusecase.SignTransactionOutput{},
			errors.New("transaction file contains no transactions")
	}

	// Step 4: Process each transaction entry
	var (
		signaturesAdded        = false // Track if we added any signatures
		hasIncompleteAfterSign = false // Track if any transactions remain incomplete
	)

	for i := range txFile.Transactions {
		tx := &txFile.Transactions[i]

		// Skip if transaction is already complete
		if tx.IsComplete {
			logger.Debug("transaction already complete, skipping",
				"uuid", tx.UUID,
				"signatureCount", tx.SignatureCount,
				"requiredSignatures", tx.RequiredSignatures)
			continue
		}

		// Get signing secret from database
		senderAccountType := domainAccount.AccountType(tx.SenderAccountType)
		secret, err := u.xrpAccountKeyRepo.GetSecret(ctx, senderAccountType, tx.SenderAccount)
		if err != nil {
			return signusecase.SignTransactionOutput{},
				fmt.Errorf("failed to get secret for transaction %s (account type: %s): %w",
					tx.UUID, senderAccountType, err)
		}

		// Determine if multi-signature is required
		isMultiSig := tx.RequiredSignatures > 1

		// Sign transaction using native Go implementation
		signedTxID, txBlob, err := u.xrp.SignTransactionNative(ctx, &tx.UnsignedData, secret, isMultiSig)
		if err != nil {
			return signusecase.SignTransactionOutput{},
				fmt.Errorf("failed to sign transaction %s (multi-sig: %t): %w",
					tx.UUID, isMultiSig, err)
		}

		// Update transaction entry with signature
		tx.SignedBlob = &txBlob
		tx.SignatureCount++
		signaturesAdded = true

		// Determine if signing is complete for this transaction
		if tx.SignatureCount >= tx.RequiredSignatures {
			tx.IsComplete = true
			logger.Debug("transaction signing complete",
				"uuid", tx.UUID,
				"signatureCount", tx.SignatureCount,
				"requiredSignatures", tx.RequiredSignatures)
		} else {
			hasIncompleteAfterSign = true
			logger.Debug("transaction partially signed",
				"uuid", tx.UUID,
				"signatureCount", tx.SignatureCount,
				"requiredSignatures", tx.RequiredSignatures,
				"remainingSignatures", tx.RequiredSignatures-tx.SignatureCount)
		}

		// Log transaction hash (never log secret)
		logger.Debug("transaction signed successfully",
			"uuid", tx.UUID,
			"signedTxID", signedTxID)
	}

	// Step 5: Determine overall completion status
	// File is complete only if ALL transactions are complete
	allComplete := !hasIncompleteAfterSign

	// Step 6: Write updated JSON file
	// Only increment signed count if we actually added signatures
	newSignedCount := signedCount
	if signaturesAdded {
		newSignedCount++
	}

	path := u.txFileRepo.CreateFilePath(actionType, domainTx.TxTypeSigned, txID, newSignedCount)
	generatedFileName, err := u.txFileRepo.WriteXRPJSONFile(path, txFile)
	if err != nil {
		return signusecase.SignTransactionOutput{},
			fmt.Errorf("failed to write signed JSON file to %s: %w", path, err)
	}

	logger.Debug("signing operation completed",
		"inputFile", input.FilePath,
		"outputFile", generatedFileName,
		"signaturesAdded", signaturesAdded,
		"allComplete", allComplete)

	return signusecase.SignTransactionOutput{
		SignedData:   "",
		IsComplete:   allComplete,
		NextFilePath: generatedFileName,
	}, nil
}
