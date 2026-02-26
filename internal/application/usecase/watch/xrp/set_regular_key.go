package xrp

import (
	"context"
	"errors"
	"fmt"

	dtoxrp "github.com/hiromaily/go-crypto-wallet/internal/application/dto/xrp"
	apixrp "github.com/hiromaily/go-crypto-wallet/internal/application/ports/api/xrp"
	file "github.com/hiromaily/go-crypto-wallet/internal/application/ports/file"
	repocold "github.com/hiromaily/go-crypto-wallet/internal/application/ports/repository/cold"
	watchusecase "github.com/hiromaily/go-crypto-wallet/internal/application/usecase/watch"
	domainXRP "github.com/hiromaily/go-crypto-wallet/internal/domain/chains/xrp"
	domainTx "github.com/hiromaily/go-crypto-wallet/internal/domain/transaction"
	"github.com/hiromaily/go-crypto-wallet/pkg/uuid"
)

// xrpRegularKeyClient defines the interface for XRP operations needed by setRegularKeyUseCase.
// This follows the Interface Segregation Principle - depend only on what you need.
type xrpRegularKeyClient interface {
	apixrp.RegularKeyPreparer
}

type setRegularKeyUseCase struct {
	xrper          xrpRegularKeyClient
	uuidHandler    uuid.UUIDHandler
	regularKeyRepo repocold.XRPRegularKeyRepositorier
	txFileRepo     file.TransactionFileRepositorier
}

// NewSetRegularKeyUseCase creates a new SetRegularKeyUseCase.
// The xrper parameter accepts any type that implements xrpRegularKeyClient (RegularKeyPreparer).
// Typically, apixrp.XRPer is passed which implements all required methods.
func NewSetRegularKeyUseCase(
	xrper xrpRegularKeyClient,
	uuidHandler uuid.UUIDHandler,
	regularKeyRepo repocold.XRPRegularKeyRepositorier,
	txFileRepo file.TransactionFileRepositorier,
) watchusecase.SetRegularKeyUseCase {
	return &setRegularKeyUseCase{
		xrper:          xrper,
		uuidHandler:    uuidHandler,
		regularKeyRepo: regularKeyRepo,
		txFileRepo:     txFileRepo,
	}
}

func (u *setRegularKeyUseCase) Execute(
	ctx context.Context,
	input watchusecase.SetRegularKeyInput,
) (watchusecase.SetRegularKeyOutput, error) {
	// Validate input
	if input.AccountAddress == "" {
		return watchusecase.SetRegularKeyOutput{}, errors.New("account address is required")
	}

	// Prepare the SetRegularKey transaction
	instructions := &dtoxrp.Instructions{
		MaxLedgerVersionOffset: domainXRP.MaxLedgerVersionOffset,
	}

	txInput, txJSON, err := u.xrper.PrepareSetRegularKeyTransaction(
		ctx, input.AccountAddress, input.RegularKeyAddress, instructions)
	if err != nil {
		return watchusecase.SetRegularKeyOutput{},
			fmt.Errorf("failed to prepare SetRegularKey transaction: %w", err)
	}

	// Generate UUID for tracking
	uid, err := u.uuidHandler.GenerateV7()
	if err != nil {
		return watchusecase.SetRegularKeyOutput{},
			fmt.Errorf("failed to generate UUID: %w", err)
	}

	// Store regular key info in database if setting a new key
	if input.RegularKeyAddress != "" {
		// Deactivate any existing regular keys for this account
		if err = u.regularKeyRepo.DeactivateByAccountID(ctx, input.AccountAddress); err != nil {
			return watchusecase.SetRegularKeyOutput{},
				fmt.Errorf("failed to deactivate existing regular keys: %w", err)
		}

		// Create new regular key record (will be activated after tx is confirmed)
		regularKey, err := domainXRP.NewXRPRegularKey(
			input.AccountAddress,
			input.RegularKeyAddress,
			"", // Public key will be set later
			"", // Public key hex will be set later
		)
		if err != nil {
			return watchusecase.SetRegularKeyOutput{},
				fmt.Errorf("failed to create regular key entity: %w", err)
		}
		regularKey.IsActive = false // Will be activated after confirmation

		if _, err = u.regularKeyRepo.Insert(ctx, regularKey); err != nil {
			return watchusecase.SetRegularKeyOutput{},
				fmt.Errorf("failed to insert regular key: %w", err)
		}
	}

	// Serialize transaction for file storage
	serializedTx := fmt.Sprintf("%s,%s", uid, txJSON)

	// Write transaction to file
	path := u.txFileRepo.CreateFilePath(
		domainTx.ActionTypeTransfer, // Using transfer action type for admin operations
		domainTx.TxTypeUnsigned,
		0, // No txID yet
		0,
	)

	generatedFileName, err := u.txFileRepo.WriteFileSlice(path, []string{
		input.AccountAddress, // First line is the account
		serializedTx,
	})
	if err != nil {
		return watchusecase.SetRegularKeyOutput{},
			fmt.Errorf("failed to write transaction file: %w", err)
	}

	return watchusecase.SetRegularKeyOutput{
		FileName:  generatedFileName,
		TxJSON:    txJSON,
		AccountID: txInput.Account,
	}, nil
}
