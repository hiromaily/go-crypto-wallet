package xrp

import (
	"context"
	"errors"
	"fmt"

	dtoxrp "github.com/hiromaily/go-crypto-wallet/internal/application/dto/xrp"
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
	// Validate file path and extract transaction metadata
	actionType, _, txID, signedCount, err := u.txFileRepo.ValidateFilePath(input.FilePath, domainTx.TxTypeUnsigned)
	if err != nil {
		return signusecase.SignTransactionOutput{}, fmt.Errorf("failed to validate file path: %w", err)
	}

	// Read and parse JSON transaction file
	jsonData, err := u.txFileRepo.ReadJSONFile(input.FilePath)
	if err != nil {
		return signusecase.SignTransactionOutput{}, fmt.Errorf("failed to read JSON file: %w", err)
	}

	txFile, err := dtoxrp.TransactionFileFromJSON(jsonData)
	if err != nil {
		return signusecase.SignTransactionOutput{}, fmt.Errorf("failed to parse transaction file: %w", err)
	}

	// Validate that transactions are unsigned
	if !txFile.IsUnsigned() {
		return signusecase.SignTransactionOutput{},
			errors.New("transaction file does not contain unsigned transactions")
	}

	// Get sender account type from file metadata
	senderAccount := domainAccount.AccountType(txFile.SenderAccountType)

	// Create signed transaction file
	signedFile := dtoxrp.NewSignedTransactionFile(txFile)

	// Sign each transaction
	for _, txEntry := range txFile.Transactions {
		if txEntry.UnsignedTx == nil {
			logger.Warn("skipping transaction entry with nil unsigned tx", "uuid", txEntry.UUID)
			continue
		}

		// Get secret from database for the sender address
		secret, err := u.xrpAccountKeyRepo.GetSecret(ctx, senderAccount, txEntry.SenderAccount)
		if err != nil {
			return signusecase.SignTransactionOutput{},
				fmt.Errorf("failed to get secret for account %s: %w", txEntry.SenderAccount, err)
		}

		// Sign the transaction
		signedTxID, txBlob, err := u.xrp.SignTransaction(ctx, txEntry.UnsignedTx, secret)
		if err != nil {
			return signusecase.SignTransactionOutput{},
				fmt.Errorf("failed to sign transaction %s: %w", txEntry.UUID, err)
		}

		logger.Debug("signed_tx",
			"uuid", txEntry.UUID,
			"signed_tx_id", signedTxID,
			"signed_tx_blob_len", len(txBlob),
		)

		// Add signed transaction to the output file
		signedFile.AddSignedTransaction(
			txEntry.UUID,
			signedTxID,
			txBlob,
			txEntry.SenderAccount,
			txEntry.SenderAccountType,
			txEntry.ReceiverAccount,
			txEntry.ReceiverAccountType,
			txEntry.Amount,
			txEntry.UnsignedTx.LastLedgerSequence,
			1, // Single signature applied
			txEntry.RequiredSignatures,
		)
	}

	// Serialize and write signed transaction file
	signedJSON, err := signedFile.ToJSON()
	if err != nil {
		return signusecase.SignTransactionOutput{}, fmt.Errorf("failed to serialize signed transactions: %w", err)
	}

	// Write signed JSON file
	path := u.txFileRepo.CreateFilePath(actionType, domainTx.TxTypeSigned, txID, signedCount+1)

	// Note: WriteJSONFile expects path without .json extension, it will add timestamp and .json
	// However, for signed files, we need to ensure the file can be parsed later
	generatedFileName, err := u.txFileRepo.WriteJSONFile(path, signedJSON)
	if err != nil {
		return signusecase.SignTransactionOutput{}, fmt.Errorf("failed to write signed JSON file: %w", err)
	}

	// Check if all transactions are complete (have required signatures)
	isComplete := len(signedFile.Transactions) > 0
	for _, tx := range signedFile.Transactions {
		if tx.SignatureCount < tx.RequiredSignatures {
			isComplete = false
			break
		}
	}

	return signusecase.SignTransactionOutput{
		SignedData:   "",
		IsComplete:   isComplete,
		NextFilePath: generatedFileName,
	}, nil
}
