package eth

import (
	"context"
	"errors"
	"fmt"

	file "github.com/hiromaily/go-crypto-wallet/internal/application/ports/file"
	signusecase "github.com/hiromaily/go-crypto-wallet/internal/application/usecase/sign"
	domainEthereum "github.com/hiromaily/go-crypto-wallet/internal/domain/ethereum"
	domainTx "github.com/hiromaily/go-crypto-wallet/internal/domain/transaction"
	domainWallet "github.com/hiromaily/go-crypto-wallet/internal/domain/wallet"
	ethereum "github.com/hiromaily/go-crypto-wallet/internal/infrastructure/api/eth"
	apiethimpl "github.com/hiromaily/go-crypto-wallet/internal/infrastructure/api/eth/eth"
	"github.com/hiromaily/go-crypto-wallet/pkg/serializer"
)

type signTransactionUseCase struct {
	eth        ethereum.Ethereumer
	txFileRepo file.TransactionFileRepositorier
	wtype      domainWallet.WalletType
}

// NewSignTransactionUseCase creates a new SignTransactionUseCase for sign wallet
func NewSignTransactionUseCase(
	ethAPI ethereum.Ethereumer,
	txFileRepo file.TransactionFileRepositorier,
	wtype domainWallet.WalletType,
) signusecase.SignTransactionUseCase {
	return &signTransactionUseCase{
		eth:        ethAPI,
		txFileRepo: txFileRepo,
		wtype:      wtype,
	}
}

func (u *signTransactionUseCase) Sign(
	ctx context.Context,
	input signusecase.SignTransactionInput,
) (signusecase.SignTransactionOutput, error) {
	// get tx_deposit_id from tx file name
	actionType, _, txID, signedCount, err := u.txFileRepo.ValidateFilePath(input.FilePath, domainTx.TxTypeUnsigned)
	if err != nil {
		return signusecase.SignTransactionOutput{}, err
	}

	// get hex tx from file
	data, err := u.txFileRepo.ReadFileSlice(input.FilePath)
	if err != nil {
		return signusecase.SignTransactionOutput{}, fmt.Errorf("fail to call txFileRepo.ReadFileSlice(): %w", err)
	}
	if len(data) <= 1 {
		return signusecase.SignTransactionOutput{}, errors.New("file is invalid")
	}
	serializedTxs := data[1:]

	txHexs := make([]string, 0, len(serializedTxs))
	for _, serializedTx := range serializedTxs {
		var rawTx domainEthereum.RawTx
		if err = serializer.GetDefaultSerializer().DecodeFromString(serializedTx, &rawTx); err != nil {
			return signusecase.SignTransactionOutput{}, fmt.Errorf("fail to call serial.DecodeFromString(): %w", err)
		}
		// sign
		var signedRawTx *domainEthereum.RawTx
		signedRawTx, err = u.eth.SignOnRawTransaction(&rawTx, apiethimpl.Password)
		if err != nil {
			return signusecase.SignTransactionOutput{}, fmt.Errorf("fail to call eth.SignOnRawTransaction(): %w", err)
		}
		txHexs = append(txHexs, fmt.Sprintf("%s,%s", rawTx.UUID, signedRawTx.TxHex))
	}

	// write file
	path := u.txFileRepo.CreateFilePath(actionType, domainTx.TxTypeSigned, txID, signedCount+1)
	generatedFileName, err := u.txFileRepo.WriteFileSlice(path, txHexs)
	if err != nil {
		return signusecase.SignTransactionOutput{}, fmt.Errorf("fail to call txFileRepo.WriteFileSlice(): %w", err)
	}

	// return hexTx, isSigned, generatedFileName, nil
	return signusecase.SignTransactionOutput{
		SignedData:   "",
		IsComplete:   true,
		NextFilePath: generatedFileName,
	}, nil
}
