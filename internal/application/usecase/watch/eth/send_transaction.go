package eth

import (
	"context"
	"errors"
	"fmt"

	apieth "github.com/hiromaily/go-crypto-wallet/internal/application/ports/api/eth"
	file "github.com/hiromaily/go-crypto-wallet/internal/application/ports/file"
	repowatch "github.com/hiromaily/go-crypto-wallet/internal/application/ports/repository/watch"
	watchusecase "github.com/hiromaily/go-crypto-wallet/internal/application/usecase/watch"
	domainTx "github.com/hiromaily/go-crypto-wallet/internal/domain/transaction"
	"github.com/hiromaily/go-crypto-wallet/pkg/logger"
)

type sendTransactionUseCase struct {
	ethClient    apieth.ETHTransactionSender
	txDetailRepo repowatch.ETHDetailTXRepositorier
	txFileRepo   file.TransactionFileRepositorier
}

// NewSendTransactionUseCase creates a new SendTransactionUseCase
func NewSendTransactionUseCase(
	ethClient apieth.ETHTransactionSender,
	txDetailRepo repowatch.ETHDetailTXRepositorier,
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

	logger.Debug("send_tx", "action_type", actionType.String(), "tx_id", txID)

	// Read signed transaction JSON file
	txFile, err := u.txFileRepo.ReadETHJSONFile(input.FilePath)
	if err != nil {
		return watchusecase.SendTransactionOutput{}, fmt.Errorf("fail to call txFileRepo.ReadETHJSONFile(): %w", err)
	}

	// Validate required fields before submission
	if txFile.UUID == "" {
		return watchusecase.SendTransactionOutput{}, errors.New("uuid is empty in transaction file")
	}
	if txFile.SignedTxHex == "" {
		return watchusecase.SendTransactionOutput{}, errors.New("signed_tx_hex is empty in transaction file")
	}

	// Send signed transaction to Ethereum network
	sentTxHash, err := u.ethClient.SendSignedRawTransaction(ctx, txFile.SignedTxHex)
	if err != nil {
		return watchusecase.SendTransactionOutput{}, fmt.Errorf("fail to call eth.SendSignedRawTransaction(): %w", err)
	}

	logger.Info("transaction submitted",
		"uuid", txFile.UUID,
		"from", txFile.From,
		"to", txFile.To,
		"tx_hash", sentTxHash,
	)

	// Update eth_detail_tx table: persist txHash and mark as TxTypeSent
	affectedNum, err := u.txDetailRepo.UpdateAfterTxSent(
		txFile.UUID, domainTx.TxTypeSent, txFile.SignedTxHex, sentTxHash)
	if err != nil {
		// Transaction is already sent; return the hash so the caller knows, but surface the DB error
		logger.Warn(
			"fail to call txDetailRepo.UpdateAfterTxSent() but tx is already sent. "+
				"So database should be updated manually",
			"uuid", txFile.UUID,
			"tx_type", domainTx.TxTypeSent.String(),
			"signed_hex_tx", txFile.SignedTxHex,
			"sent_hash_tx", sentTxHash,
		)
		return watchusecase.SendTransactionOutput{TxID: sentTxHash},
			fmt.Errorf("fail to call txDetailRepo.UpdateAfterTxSent(): %w", err)
	}
	if affectedNum == 0 {
		logger.Info("no records updated in eth_detail_tx",
			"uuid", txFile.UUID,
			"tx_type", domainTx.TxTypeSent.String(),
			"sent_hash_tx", sentTxHash,
		)
	}

	return watchusecase.SendTransactionOutput{
		TxID: sentTxHash,
	}, nil
}
