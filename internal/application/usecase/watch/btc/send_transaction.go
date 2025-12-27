package btc

import (
	"context"
	"errors"
	"fmt"
	"strings"

	watchusecase "github.com/hiromaily/go-crypto-wallet/internal/application/usecase/watch"
	domainTx "github.com/hiromaily/go-crypto-wallet/internal/domain/transaction"
	"github.com/hiromaily/go-crypto-wallet/internal/infrastructure/api/bitcoin"
	watchrepo "github.com/hiromaily/go-crypto-wallet/internal/infrastructure/repository/watch"
	"github.com/hiromaily/go-crypto-wallet/internal/infrastructure/storage/file"
	"github.com/hiromaily/go-crypto-wallet/pkg/logger"
)

type sendTransactionUseCase struct {
	btcClient    bitcoin.Bitcoiner
	addrRepo     watchrepo.AddressRepositorier
	txRepo       watchrepo.BTCTxRepositorier
	txOutputRepo watchrepo.TxOutputRepositorier
	txFileRepo   file.TransactionFileRepositorier
}

// NewSendTransactionUseCase creates a new SendTransactionUseCase
func NewSendTransactionUseCase(
	btcClient bitcoin.Bitcoiner,
	addrRepo watchrepo.AddressRepositorier,
	txRepo watchrepo.BTCTxRepositorier,
	txOutputRepo watchrepo.TxOutputRepositorier,
	txFileRepo file.TransactionFileRepositorier,
) watchusecase.SendTransactionUseCase {
	return &sendTransactionUseCase{
		btcClient:    btcClient,
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
	// Validate file path and extract transaction metadata
	actionType, _, txID, _, err := u.txFileRepo.ValidateFilePath(input.FilePath, domainTx.TxTypeSigned)
	if err != nil {
		return watchusecase.SendTransactionOutput{}, fmt.Errorf("invalid file path: %w", err)
	}

	logger.Debug("sending transaction", "action_type", actionType.String(), "tx_id", txID)

	// Determine file format based on extension and read accordingly
	var signedHex string
	if isPSBTFile(input.FilePath) {
		// PSBT flow: finalize → extract → convert to hex
		signedHex, err = u.processPSBTFile(input.FilePath)
		if err != nil {
			return watchusecase.SendTransactionOutput{}, fmt.Errorf("failed to process PSBT file: %w", err)
		}
	} else {
		// Legacy flow: read hex directly from file
		signedHex, err = u.txFileRepo.ReadFile(input.FilePath)
		if err != nil {
			return watchusecase.SendTransactionOutput{}, fmt.Errorf("failed to read transaction file: %w", err)
		}
	}

	// Broadcast transaction to Bitcoin network
	hash, err := u.btcClient.SendTransactionByHex(signedHex)
	if err != nil {
		return watchusecase.SendTransactionOutput{}, fmt.Errorf("failed to broadcast transaction: %w", err)
	}

	// Check if transaction was already sent
	if hash == nil {
		logger.Info("transaction already sent", "tx_id", txID)
		return watchusecase.SendTransactionOutput{TxID: ""}, nil
	}

	// Update transaction status in database
	affectedNum, err := u.txRepo.UpdateAfterTxSent(txID, domainTx.TxTypeSent, signedHex, hash.String())
	if err != nil {
		// Critical: transaction is broadcasted but database update failed
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

	// Update address allocation status (skip for payment transactions with anonymous receivers)
	if actionType != domainTx.ActionTypePayment {
		if err := u.updateAddressAllocation(txID); err != nil {
			// Critical: transaction sent and DB updated, but address allocation failed
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

// processPSBTFile processes a PSBT file: validates, finalizes, extracts, and converts to hex
func (u *sendTransactionUseCase) processPSBTFile(filePath string) (string, error) {
	// Read PSBT from file
	psbtBase64, err := u.txFileRepo.ReadPSBTFile(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to read PSBT file: %w", err)
	}

	// Validate PSBT is fully signed
	isComplete, err := u.btcClient.IsPSBTComplete(psbtBase64)
	if err != nil {
		return "", fmt.Errorf("failed to validate PSBT: %w", err)
	}
	if !isComplete {
		return "", errors.New("PSBT is not fully signed - cannot finalize incomplete PSBT")
	}

	logger.Debug("PSBT validation passed", "is_complete", isComplete)

	// Finalize PSBT (combine signatures into final scriptSig/witness)
	finalizedPSBT, err := u.btcClient.FinalizePSBT(psbtBase64)
	if err != nil {
		return "", fmt.Errorf("failed to finalize PSBT: %w", err)
	}

	logger.Debug("PSBT finalized successfully")

	// Extract final transaction from PSBT
	msgTx, err := u.btcClient.ExtractTransaction(finalizedPSBT)
	if err != nil {
		return "", fmt.Errorf("failed to extract transaction from PSBT: %w", err)
	}

	// Convert transaction to hex for broadcasting
	hexTx, err := u.btcClient.ToHex(msgTx)
	if err != nil {
		return "", fmt.Errorf("failed to convert transaction to hex: %w", err)
	}

	logger.Debug("transaction extracted from PSBT", "hex_length", len(hexTx))

	return hexTx, nil
}

// isPSBTFile checks if the file path indicates a PSBT file (by extension)
func isPSBTFile(filePath string) bool {
	return strings.HasSuffix(strings.ToLower(filePath), ".psbt")
}

// updateAddressAllocation marks the receiver address as allocated
func (u *sendTransactionUseCase) updateAddressAllocation(txID int64) error {
	// Get transaction outputs
	txOutputs, err := u.txOutputRepo.GetAllByTxID(txID)
	if err != nil {
		return fmt.Errorf("failed to get transaction outputs: %w", err)
	}

	if len(txOutputs) == 0 {
		return errors.New("no transaction outputs found")
	}

	// Mark first output address as allocated
	_, err = u.addrRepo.UpdateIsAllocated(true, txOutputs[0].OutputAddress)
	if err != nil {
		return fmt.Errorf("failed to update address allocation status: %w", err)
	}

	return nil
}
