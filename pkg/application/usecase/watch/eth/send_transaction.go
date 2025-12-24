package eth

import (
	"context"
	"errors"
	"fmt"
	"strings"

	watchusecase "github.com/hiromaily/go-crypto-wallet/pkg/application/usecase/watch"
	domainTx "github.com/hiromaily/go-crypto-wallet/pkg/domain/transaction"
	"github.com/hiromaily/go-crypto-wallet/pkg/infrastructure/api/ethereum"
	watchrepo "github.com/hiromaily/go-crypto-wallet/pkg/infrastructure/repository/watch"
	"github.com/hiromaily/go-crypto-wallet/pkg/infrastructure/storage/file"
	"github.com/hiromaily/go-crypto-wallet/pkg/logger"
)

type sendTransactionUseCase struct {
	ethClient    ethereum.Ethereumer
	txDetailRepo watchrepo.EthDetailTxRepositorier
	txFileRepo   file.TransactionFileRepositorier
}

// NewSendTransactionUseCase creates a new SendTransactionUseCase
func NewSendTransactionUseCase(
	ethClient ethereum.Ethereumer,
	txDetailRepo watchrepo.EthDetailTxRepositorier,
	txFileRepo file.TransactionFileRepositorier,
) watchusecase.SendTransactionUseCase {
	return &sendTransactionUseCase{
		ethClient:    ethClient,
		txDetailRepo: txDetailRepo,
		txFileRepo:   txFileRepo,
	}
}

func (u *sendTransactionUseCase) Execute(
	ctx context.Context,
	input watchusecase.SendTransactionInput,
) (watchusecase.SendTransactionOutput, error) {
	// Validate file path and extract transaction metadata
	actionType, _, txID, _, err := u.txFileRepo.ValidateFilePath(input.FilePath, domainTx.TxTypeSigned)
	if err != nil {
		return watchusecase.SendTransactionOutput{}, fmt.Errorf("fail to call txFileRepo.ValidateFilePath(): %w", err)
	}

	logger.Debug("send_tx", "action_type", actionType.String())

	// Read hex from file
	data, err := u.txFileRepo.ReadFileSlice(input.FilePath)
	if err != nil {
		return watchusecase.SendTransactionOutput{}, fmt.Errorf("fail to call txFileRepo.ReadFile(): %w", err)
	}

	// Process each signed transaction from the file
	for _, txHex := range data {
		// data is csv [rawTx.TxHex, signedRawTx.TxHex]
		// rawTx.TxHex is used to record status by updating database
		tmp := strings.Split(txHex, ",")
		if len(tmp) != 2 {
			return watchusecase.SendTransactionOutput{}, errors.New("data format is invalid in file")
		}
		uuid := tmp[0]
		signedTx := tmp[1]

		// Send signed transaction to Ethereum network
		var sentTx string
		sentTx, err = u.ethClient.SendSignedRawTransaction(ctx, signedTx)
		if err != nil {
			logger.Warn("fail to call eth.SendSignedRawTransaction()",
				"error", err,
			)
			continue
		}
		if sentTx == "" {
			logger.Warn("no sentTx by calling eth.SendSignedRawTransaction()",
				"error", err,
			)
			continue
		}

		// Update eth_detail_tx table
		var affectedNum int64
		affectedNum, err = u.txDetailRepo.UpdateAfterTxSent(uuid, domainTx.TxTypeSent, signedTx, sentTx)
		if err != nil {
			// TODO: even if error occurred, tx is already sent. so db should be corrected manually
			logger.Warn(
				"fail to call repo.Tx().UpdateAfterTxSent() but tx is already sent. "+
					"So database should be updated manually",
				"tx_id", txID,
				"tx_type", domainTx.TxTypeSent.String(),
				"tx_type_value", domainTx.TxTypeSent.Int8(),
				"signed_hex_tx", signedTx,
				"sent_hash_tx", sentTx,
			)
			continue
		}
		if affectedNum == 0 {
			logger.Info("no records to update tx_table",
				"tx_id", txID,
				"tx_type", domainTx.TxTypeSent.String(),
				"tx_type_value", domainTx.TxTypeSent.Int8(),
				"signed_hex_tx", signedTx,
				"sent_hash_tx", sentTx,
			)
			continue
		}
	}

	// TODO: update is_allocated in account_pubkey_table
	// Ethereum should use same address because no utxo
	return watchusecase.SendTransactionOutput{
		TxID: "",
	}, nil
}
