package xrp

import (
	"context"
	"errors"
	"fmt"

	dtoRipple "github.com/hiromaily/go-crypto-wallet/internal/application/dto/ripple"
	apixrp "github.com/hiromaily/go-crypto-wallet/internal/application/ports/api/xrp"
	file "github.com/hiromaily/go-crypto-wallet/internal/application/ports/file"
	repocold "github.com/hiromaily/go-crypto-wallet/internal/application/ports/repository/cold"
	watchusecase "github.com/hiromaily/go-crypto-wallet/internal/application/usecase/watch"
	domainTx "github.com/hiromaily/go-crypto-wallet/internal/domain/transaction"
	domainXrp "github.com/hiromaily/go-crypto-wallet/internal/domain/xrp"
	"github.com/hiromaily/go-crypto-wallet/pkg/uuid"
)

type setRegularKeyUseCase struct {
	rippler        apixrp.Rippler
	uuidHandler    uuid.UUIDHandler
	regularKeyRepo repocold.XRPRegularKeyRepositorier
	txFileRepo     file.TransactionFileRepositorier
}

// NewSetRegularKeyUseCase creates a new SetRegularKeyUseCase
func NewSetRegularKeyUseCase(
	rippler apixrp.Rippler,
	uuidHandler uuid.UUIDHandler,
	regularKeyRepo repocold.XRPRegularKeyRepositorier,
	txFileRepo file.TransactionFileRepositorier,
) watchusecase.SetRegularKeyUseCase {
	return &setRegularKeyUseCase{
		rippler:        rippler,
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
	instructions := &dtoRipple.Instructions{
		MaxLedgerVersionOffset: domainXrp.MaxLedgerVersionOffset,
	}

	txInput, txJSON, err := u.rippler.PrepareSetRegularKeyTransaction(
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
		regularKey, err := domainXrp.NewXRPRegularKey(
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
