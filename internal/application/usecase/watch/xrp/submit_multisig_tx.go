package xrp

import (
	"context"
	"errors"
	"fmt"

	apixrp "github.com/hiromaily/go-crypto-wallet/internal/application/ports/api/xrp"
	repowatch "github.com/hiromaily/go-crypto-wallet/internal/application/ports/repository/watch"
	watchusecase "github.com/hiromaily/go-crypto-wallet/internal/application/usecase/watch"
	domainXRP "github.com/hiromaily/go-crypto-wallet/internal/domain/chains/xrp"
)

type submitMultisigTxUseCase struct {
	xrper               apixrp.TransactionSubmitter
	pendingMultisigRepo repowatch.XRPPendingMultisigRepositorier
}

// NewSubmitMultisigTxUseCase creates a new SubmitMultisigTxUseCase.
// The xrper parameter accepts any type that implements TransactionSubmitter.
// Typically, apixrp.XRPer is passed which implements all required methods.
func NewSubmitMultisigTxUseCase(
	xrper apixrp.TransactionSubmitter,
	pendingMultisigRepo repowatch.XRPPendingMultisigRepositorier,
) watchusecase.SubmitMultisigTxUseCase {
	return &submitMultisigTxUseCase{
		xrper:               xrper,
		pendingMultisigRepo: pendingMultisigRepo,
	}
}

func (u *submitMultisigTxUseCase) Execute(
	ctx context.Context,
	input watchusecase.SubmitMultisigTxInput,
) (watchusecase.SubmitMultisigTxOutput, error) {
	// Validate input
	if input.TxUUID == "" {
		return watchusecase.SubmitMultisigTxOutput{}, errors.New("transaction UUID is required")
	}

	// Get the pending multisig transaction
	pendingTx, err := u.pendingMultisigRepo.GetByUUID(ctx, input.TxUUID)
	if err != nil {
		return watchusecase.SubmitMultisigTxOutput{},
			fmt.Errorf("failed to get pending transaction: %w", err)
	}
	if pendingTx == nil {
		return watchusecase.SubmitMultisigTxOutput{},
			fmt.Errorf("pending transaction not found: %s", input.TxUUID)
	}

	// Check if transaction is ready for submission
	if pendingTx.Status != domainXRP.MultisigStatusReady {
		return watchusecase.SubmitMultisigTxOutput{},
			fmt.Errorf("transaction is not ready for submission (status: %s)", pendingTx.Status)
	}

	if pendingTx.CombinedTxBlob == nil || *pendingTx.CombinedTxBlob == "" {
		return watchusecase.SubmitMultisigTxOutput{},
			errors.New("combined transaction blob is missing")
	}

	// Submit the combined transaction
	sentTx, _, err := u.xrper.SubmitTransaction(ctx, *pendingTx.CombinedTxBlob)
	if err != nil {
		// Mark as failed
		if updateErr := u.pendingMultisigRepo.SetFailed(ctx, pendingTx.ID); updateErr != nil {
			return watchusecase.SubmitMultisigTxOutput{},
				fmt.Errorf("submit and status update failed: %w (update: %v)", err, updateErr)
		}
		return watchusecase.SubmitMultisigTxOutput{},
			fmt.Errorf("failed to submit transaction: %w", err)
	}

	// Get the transaction hash from the response
	txHash := sentTx.TxJSON.Hash
	if txHash == "" {
		// Try to get hash from the result
		txHash = "queued" // Placeholder if hash not immediately available
	}

	// Update pending transaction as submitted
	if err = u.pendingMultisigRepo.SetSubmitted(ctx, pendingTx.ID, txHash); err != nil {
		return watchusecase.SubmitMultisigTxOutput{},
			fmt.Errorf("failed to update transaction status: %w", err)
	}

	return watchusecase.SubmitMultisigTxOutput{
		TxHash:   txHash,
		IsQueued: sentTx.ResultCode == "tesSUCCESS",
	}, nil
}
