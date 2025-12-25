package eth

import (
	"context"
	"errors"
	"fmt"

	keygenusecase "github.com/hiromaily/go-crypto-wallet/pkg/application/usecase/keygen"
	domainTx "github.com/hiromaily/go-crypto-wallet/pkg/domain/transaction"
	"github.com/hiromaily/go-crypto-wallet/pkg/infrastructure/api/ethereum"
	"github.com/hiromaily/go-crypto-wallet/pkg/infrastructure/api/ethereum/eth"
	"github.com/hiromaily/go-crypto-wallet/pkg/infrastructure/api/ethereum/ethtx"
	"github.com/hiromaily/go-crypto-wallet/pkg/infrastructure/storage/file"
	"github.com/hiromaily/go-crypto-wallet/pkg/serial"
)

type signTransactionUseCase struct {
	eth        ethereum.Ethereumer
	txFileRepo file.TransactionFileRepositorier
}

// NewSignTransactionUseCase creates a new SignTransactionUseCase for ETH keygen
func NewSignTransactionUseCase(
	eth ethereum.Ethereumer,
	txFileRepo file.TransactionFileRepositorier,
) keygenusecase.SignTransactionUseCase {
	return &signTransactionUseCase{
		eth:        eth,
		txFileRepo: txFileRepo,
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

	// Get hex tx from file
	data, err := u.txFileRepo.ReadFileSlice(input.FilePath)
	if err != nil {
		return keygenusecase.SignTransactionOutput{}, fmt.Errorf("fail to call txFileRepo.ReadFileSlice(): %w", err)
	}
	if len(data) <= 1 {
		return keygenusecase.SignTransactionOutput{}, errors.New("file is invalid")
	}
	serializedTxs := data[1:]

	txHexs := make([]string, 0, len(serializedTxs))
	for _, serializedTx := range serializedTxs {
		var rawTx ethtx.RawTx
		if err = serial.DecodeFromString(serializedTx, &rawTx); err != nil {
			return keygenusecase.SignTransactionOutput{}, fmt.Errorf("fail to call serial.DecodeFromString(): %w", err)
		}

		// Sign
		var signedRawTx *ethtx.RawTx
		signedRawTx, err = u.eth.SignOnRawTransaction(&rawTx, eth.Password)
		if err != nil {
			return keygenusecase.SignTransactionOutput{}, fmt.Errorf("fail to call eth.SignOnRawTransaction(): %w", err)
		}
		txHexs = append(txHexs, fmt.Sprintf("%s,%s", rawTx.UUID, signedRawTx.TxHex))
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
		SignedCount:   1, // ETH signs one transaction at a time
		UnsignedCount: 0,
	}, nil
}
