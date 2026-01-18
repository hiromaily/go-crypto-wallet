package bch

import (
	"context"
	"errors"
	"fmt"
	"strings"

	apibtc "github.com/hiromaily/go-crypto-wallet/internal/application/ports/api/btc"
	file "github.com/hiromaily/go-crypto-wallet/internal/application/ports/file"
	repowatch "github.com/hiromaily/go-crypto-wallet/internal/application/ports/repository/watch"
	watchusecase "github.com/hiromaily/go-crypto-wallet/internal/application/usecase/watch"
	domainTx "github.com/hiromaily/go-crypto-wallet/internal/domain/transaction"
	"github.com/hiromaily/go-crypto-wallet/pkg/logger"
)

// sendTxBCHClient defines the minimal interface needed for BCH transaction sending.
// This follows the Interface Segregation Principle - depend only on methods actually used.
// BCH does not use PSBT, so this interface only needs TransactionSender.
type sendTxBCHClient interface {
	apibtc.TransactionSender
}

// sendTransactionUseCase implements BCH-specific transaction broadcasting.
// BCH uses Raw Transaction Hex format instead of PSBT.
type sendTransactionUseCase struct {
	bchClient    sendTxBCHClient
	addrRepo     repowatch.AddressRepositorier
	txRepo       repowatch.BTCTxRepositorier
	txOutputRepo repowatch.TxOutputRepositorier
	txFileRepo   file.TransactionFileRepositorier
}

// NewBCHSendTransactionUseCase creates a new BCH SendTransactionUseCase.
func NewBCHSendTransactionUseCase(
	bchClient apibtc.BCHer,
	addrRepo repowatch.AddressRepositorier,
	txRepo repowatch.BTCTxRepositorier,
	txOutputRepo repowatch.TxOutputRepositorier,
	txFileRepo file.TransactionFileRepositorier,
) watchusecase.SendTransactionUseCase {
	return &sendTransactionUseCase{
		bchClient:    bchClient,
		addrRepo:     addrRepo,
		txRepo:       txRepo,
		txOutputRepo: txOutputRepo,
		txFileRepo:   txFileRepo,
	}
}

func (u *sendTransactionUseCase) Execute(
	ctx context.Context,
	input watchusecase.SendTransactionInput,
) (watchusecase.SendTransactionOutput, error) {
	// Extract transaction metadata from file path
	fileNameType, err := u.txFileRepo.GetFileNameType(input.FilePath)
	if err != nil {
		return watchusecase.SendTransactionOutput{}, fmt.Errorf("invalid file path: %w", err)
	}
	actionType := fileNameType.ActionType
	txID := fileNameType.TxID

	logger.Debug("sending BCH transaction",
		"action_type", actionType.String(),
		"tx_id", txID,
		"tx_type", fileNameType.TxType.String())

	// Read signed transaction hex from file (BCH format)
	signedHex, err := u.processRawTxFile(input.FilePath)
	if err != nil {
		return watchusecase.SendTransactionOutput{}, fmt.Errorf("failed to process tx file: %w", err)
	}

	// Broadcast transaction to BCH network
	hash, err := u.bchClient.SendTransactionByHex(signedHex)
	if err != nil {
		return watchusecase.SendTransactionOutput{}, fmt.Errorf("failed to broadcast transaction: %w", err)
	}

	if hash == nil {
		logger.Info("transaction already sent", "tx_id", txID)
		return watchusecase.SendTransactionOutput{TxID: ""}, nil
	}

	// Update transaction status in database
	affectedNum, err := u.txRepo.UpdateAfterTxSent(txID, domainTx.TxTypeSent, signedHex, hash.String())
	if err != nil {
		logger.Warn(
			"transaction sent but database update failed - manual correction required",
			"tx_id", txID,
			"tx_type", domainTx.TxTypeSent.String(),
			"signed_hex", signedHex,
			"tx_hash", hash.String(),
			"error", err)
		return watchusecase.SendTransactionOutput{},
			fmt.Errorf("transaction sent but database update failed (txID: %d): %w", txID, err)
	}

	if affectedNum == 0 {
		logger.Info("no records updated",
			"tx_id", txID,
			"tx_hash", hash.String())
		return watchusecase.SendTransactionOutput{TxID: hash.String()}, nil
	}

	// Update address allocation status
	if actionType != domainTx.ActionTypePayment {
		if err := u.updateAddressAllocation(txID); err != nil {
			logger.Error(
				"transaction sent but address allocation update failed - manual correction required",
				"tx_id", txID,
				"error", err)
			return watchusecase.SendTransactionOutput{}, err
		}
	}

	return watchusecase.SendTransactionOutput{
		TxID: hash.String(),
	}, nil
}

// processRawTxFile reads and extracts the transaction hex from BCH raw tx file.
// BCH files contain the hex on the first line, followed by optional prevtx metadata.
func (u *sendTransactionUseCase) processRawTxFile(filePath string) (string, error) {
	// Read file content
	content, err := u.txFileRepo.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to read tx file: %w", err)
	}

	// Extract hex (first line)
	lines := strings.Split(strings.TrimSpace(content), "\n")
	if len(lines) == 0 || lines[0] == "" {
		return "", errors.New("empty transaction file")
	}

	txHex := strings.TrimSpace(lines[0])

	logger.Debug("BCH transaction extracted from file", "hex_length", len(txHex))

	return txHex, nil
}

// updateAddressAllocation marks the receiver address as allocated
func (u *sendTransactionUseCase) updateAddressAllocation(txID int64) error {
	txOutputs, err := u.txOutputRepo.GetAllByTxID(txID)
	if err != nil {
		return fmt.Errorf("failed to get transaction outputs: %w", err)
	}

	if len(txOutputs) == 0 {
		return errors.New("no transaction outputs found")
	}

	_, err = u.addrRepo.UpdateIsAllocated(true, txOutputs[0].OutputAddress)
	if err != nil {
		return fmt.Errorf("failed to update address allocation status: %w", err)
	}

	return nil
}
