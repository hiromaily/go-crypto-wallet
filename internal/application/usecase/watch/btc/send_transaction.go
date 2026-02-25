// Package btc provides BTC-specific transaction sending use cases.
//
// WARNING: This file is for BTC (Bitcoin) ONLY.
// If cointype is BCH (Bitcoin Cash), DO NOT modify this file.
// Use the BCH-specific implementation instead:
//
//	internal/application/usecase/watch/bch/send_transaction.go
//
// Key differences:
//   - BTC: Processes PSBT files, uses PSBTFinalizer interface
//   - BCH: Processes raw hex files directly
package btc

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/btcsuite/btcd/wire"

	apibtc "github.com/hiromaily/go-crypto-wallet/internal/application/ports/api/btc"
	file "github.com/hiromaily/go-crypto-wallet/internal/application/ports/file"
	repowatch "github.com/hiromaily/go-crypto-wallet/internal/application/ports/repository/watch"
	watchusecase "github.com/hiromaily/go-crypto-wallet/internal/application/usecase/watch"
	domainTx "github.com/hiromaily/go-crypto-wallet/internal/domain/transaction"
	"github.com/hiromaily/go-crypto-wallet/pkg/logger"
)

// sendTxBTCClient defines the minimal interface needed for BTC transaction sending.
// This follows the Interface Segregation Principle - depend only on methods actually used.
type sendTxBTCClient interface {
	apibtc.TransactionSender
	apibtc.PSBTFinalizer
	// ToHex is needed to convert extracted transaction to hex for broadcast
	ToHex(tx *wire.MsgTx) (string, error)
}

type sendTransactionUseCase struct {
	btcClient    sendTxBTCClient
	addrRepo     repowatch.AddressRepositorier
	txRepo       repowatch.BTCTxRepositorier
	txOutputRepo repowatch.TxOutputRepositorier
	txFileRepo   file.TransactionFileRepositorier
}

// NewSendTransactionUseCase creates a new SendTransactionUseCase
func NewSendTransactionUseCase(
	btcClient sendTxBTCClient,
	addrRepo repowatch.AddressRepositorier,
	txRepo repowatch.BTCTxRepositorier,
	txOutputRepo repowatch.TxOutputRepositorier,
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
	// Extract transaction metadata from file path
	// Note: We don't strictly validate txType here because:
	// 1. For PSBT files, the actual signature validation happens in processPSBTFile via IsPSBTComplete()
	// 2. The btcd PSBT library's IsComplete() detection may not work correctly for all multisig configurations
	// 3. The real validation is whether the PSBT has enough signatures to be finalized
	fileNameType, err := u.txFileRepo.GetFileNameType(input.FilePath)
	if err != nil {
		return watchusecase.SendTransactionOutput{}, fmt.Errorf("invalid file path: %w", err)
	}
	actionType := fileNameType.ActionType
	txID := fileNameType.TxID

	logger.Debug("sending transaction",
		"action_type", actionType.String(),
		"tx_id", txID,
		"tx_type", fileNameType.TxType.String())

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
