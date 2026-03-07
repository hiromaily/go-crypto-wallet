package xrp

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	dtoxrp "github.com/hiromaily/go-crypto-wallet/internal/application/dto/xrp"
	apixrp "github.com/hiromaily/go-crypto-wallet/internal/application/ports/api/xrp"
	file "github.com/hiromaily/go-crypto-wallet/internal/application/ports/file"
	repocold "github.com/hiromaily/go-crypto-wallet/internal/application/ports/repository/cold"
	keygenusecase "github.com/hiromaily/go-crypto-wallet/internal/application/usecase/keygen"
	domainAccount "github.com/hiromaily/go-crypto-wallet/internal/domain/account"
	domainTx "github.com/hiromaily/go-crypto-wallet/internal/domain/transaction"
	"github.com/hiromaily/go-crypto-wallet/pkg/logger"
)

// xrpKeygenSignClient defines the interface for XRP operations needed by signTransactionUseCase in keygen.
// This follows the Interface Segregation Principle - depend only on what you need.
type xrpKeygenSignClient interface {
	apixrp.TransactionSigner
}

type signTransactionUseCase struct {
	xrp               xrpKeygenSignClient
	xrpAccountKeyRepo repocold.XRPAccountKeyRepositorier
	txFileRepo        file.TransactionFileRepositorier
}

// NewSignTransactionUseCase creates a new SignTransactionUseCase for XRP keygen.
// The xrp parameter accepts any type that implements xrpKeygenSignClient (TransactionSigner).
// Typically, apixrp.XRPer is passed which implements all required methods.
func NewSignTransactionUseCase(
	xrp xrpKeygenSignClient,
	xrpAccountKeyRepo repocold.XRPAccountKeyRepositorier,
	txFileRepo file.TransactionFileRepositorier,
) keygenusecase.SignTransactionUseCase {
	return &signTransactionUseCase{
		xrp:               xrp,
		xrpAccountKeyRepo: xrpAccountKeyRepo,
		txFileRepo:        txFileRepo,
	}
}

func (u *signTransactionUseCase) Sign(
	ctx context.Context,
	input keygenusecase.SignTransactionInput,
) (keygenusecase.SignTransactionOutput, error) {
	// Get tx_deposit_id from tx file name
	actionType, _, txID, signedCount, err := u.txFileRepo.ValidateFilePath(input.FilePath, domainTx.TxTypeUnsigned)
	if err != nil {
		return keygenusecase.SignTransactionOutput{}, err
	}

	var senderAccount domainAccount.AccountType

	// Get hex tx from file
	data, err := u.txFileRepo.ReadFileSlice(input.FilePath)
	if err != nil {
		return keygenusecase.SignTransactionOutput{}, fmt.Errorf("fail to call txFileRepo.ReadFileSlice(): %w", err)
	}
	if len(data) > 1 {
		senderAccount = domainAccount.AccountType(data[0])
	} else {
		return keygenusecase.SignTransactionOutput{}, errors.New("file is invalid")
	}
	serializedTxs := data[1:]

	txHexs := make([]string, 0, len(serializedTxs))
	for _, serializedTx := range serializedTxs {
		// uuid, txJSON
		tmp := strings.SplitAfterN(serializedTx, ",", 2)
		if len(tmp) != 2 {
			return keygenusecase.SignTransactionOutput{}, errors.New("data format is invalid in file")
		}
		uuid := strings.TrimRight(tmp[0], ",")
		txJSON := tmp[1]

		var txInput dtoxrp.TxInput
		if err = json.Unmarshal([]byte(txJSON), &txInput); err != nil {
			return keygenusecase.SignTransactionOutput{}, fmt.Errorf("fail to call json.Unmarshal(txJSON): %w", err)
		}
		// Populate RawTxFields to preserve transaction-type-specific fields
		// (e.g. SignerQuorum/SignerEntries for SignerListSet) not captured in TxInput.
		var rawFields map[string]any
		if err = json.Unmarshal([]byte(txJSON), &rawFields); err != nil {
			return keygenusecase.SignTransactionOutput{}, fmt.Errorf("fail to unmarshal raw tx fields: %w", err)
		}
		txInput.RawTxFields = rawFields

		// TODO: get secret from database by txInput.Account
		// master_seed from xrp_account_key table
		var secret string
		secret, err = u.xrpAccountKeyRepo.GetSecret(ctx, senderAccount, txInput.Account)
		if err != nil {
			return keygenusecase.SignTransactionOutput{},
				fmt.Errorf("fail to call xrpAccountKeyRepo.GetSecret(): %w", err)
		}

		// Sign
		var signedTxID string
		var txBlob string
		signedTxID, txBlob, err = u.xrp.SignTransaction(ctx, &txInput, secret)
		if err != nil {
			return keygenusecase.SignTransactionOutput{}, fmt.Errorf("fail to call xrp.SignTransaction(): %w", err)
		}
		logger.Debug("signed_tx",
			"uuid", uuid, "signed_tx_id", signedTxID, "signed_tx_blob", txBlob)
		txHexs = append(txHexs, fmt.Sprintf("%s,%s,%s", uuid, signedTxID, txBlob))
	}

	// Write file
	path := u.txFileRepo.CreateFilePath(actionType, domainTx.TxTypeSigned, txID, signedCount+1)
	generatedFileName, err := u.txFileRepo.WriteFileSlice(path, txHexs)
	if err != nil {
		return keygenusecase.SignTransactionOutput{}, fmt.Errorf("fail to call txFileRepo.WriteFileSlice(): %w", err)
	}

	return keygenusecase.SignTransactionOutput{
		FilePath:      generatedFileName,
		IsDone:        true,
		SignedCount:   1, // XRP signs one transaction at a time
		UnsignedCount: 0,
	}, nil
}
