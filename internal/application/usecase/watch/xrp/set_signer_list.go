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

// xrpSignerListClient defines the interface for XRP operations needed by setSignerListUseCase.
// This follows the Interface Segregation Principle - depend only on what you need.
type xrpSignerListClient interface {
	apixrp.SignerListPreparer
}

type setSignerListUseCase struct {
	xrper           xrpSignerListClient
	uuidHandler     uuid.UUIDHandler
	signerListRepo  repocold.XRPSignerListRepositorier
	signerEntryRepo repocold.XRPSignerEntryRepositorier
	txFileRepo      file.TransactionFileRepositorier
}

// NewSetSignerListUseCase creates a new SetSignerListUseCase.
// The xrper parameter accepts any type that implements xrpSignerListClient (SignerListPreparer).
// Typically, apixrp.XRPer is passed which implements all required methods.
func NewSetSignerListUseCase(
	xrper xrpSignerListClient,
	uuidHandler uuid.UUIDHandler,
	signerListRepo repocold.XRPSignerListRepositorier,
	signerEntryRepo repocold.XRPSignerEntryRepositorier,
	txFileRepo file.TransactionFileRepositorier,
) watchusecase.SetSignerListUseCase {
	return &setSignerListUseCase{
		xrper:           xrper,
		uuidHandler:     uuidHandler,
		signerListRepo:  signerListRepo,
		signerEntryRepo: signerEntryRepo,
		txFileRepo:      txFileRepo,
	}
}

func (u *setSignerListUseCase) Execute(
	ctx context.Context,
	input watchusecase.SetSignerListInput,
) (watchusecase.SetSignerListOutput, error) {
	// Validate input
	if input.AccountAddress == "" {
		return watchusecase.SetSignerListOutput{}, errors.New("account address is required")
	}

	// Validate signer entries using domain validator
	signerInputs := make([]domainXRP.SignerEntryInput, len(input.SignerEntries))
	for i, entry := range input.SignerEntries {
		signerInputs[i] = domainXRP.SignerEntryInput{
			Account: entry.Account,
			Weight:  entry.Weight,
		}
	}

	if err := domainXRP.ValidateSignerEntries(signerInputs, input.SignerQuorum); err != nil {
		return watchusecase.SetSignerListOutput{},
			fmt.Errorf("invalid signer entries: %w", err)
	}

	// Convert to API input format
	apiSignerEntries := make([]apixrp.SignerEntryInput, len(input.SignerEntries))
	for i, entry := range input.SignerEntries {
		apiSignerEntries[i] = apixrp.SignerEntryInput{
			Account: entry.Account,
			Weight:  entry.Weight,
		}
	}

	// Prepare the SignerListSet transaction
	instructions := &dtoxrp.Instructions{
		MaxLedgerVersionOffset: domainXRP.MaxLedgerVersionOffset,
	}

	_, txJSON, err := u.xrper.PrepareSignerListSetTransaction(
		ctx, input.AccountAddress, input.SignerQuorum, apiSignerEntries, instructions)
	if err != nil {
		return watchusecase.SetSignerListOutput{},
			fmt.Errorf("failed to prepare SignerListSet transaction: %w", err)
	}

	// Generate UUID for tracking
	uid, err := u.uuidHandler.GenerateV7()
	if err != nil {
		return watchusecase.SetSignerListOutput{},
			fmt.Errorf("failed to generate UUID: %w", err)
	}

	// Deactivate any existing signer lists for this account
	if err = u.signerListRepo.DeactivateByAccountID(ctx, input.AccountAddress); err != nil {
		return watchusecase.SetSignerListOutput{},
			fmt.Errorf("failed to deactivate existing signer lists: %w", err)
	}

	// Create new signer list record (will be activated after tx is confirmed)
	signerList, err := domainXRP.NewXRPSignerList(input.AccountAddress, input.SignerQuorum)
	if err != nil {
		return watchusecase.SetSignerListOutput{},
			fmt.Errorf("failed to create signer list entity: %w", err)
	}
	signerList.IsActive = false // Will be activated after confirmation

	signerListID, err := u.signerListRepo.Insert(ctx, signerList)
	if err != nil {
		return watchusecase.SetSignerListOutput{},
			fmt.Errorf("failed to insert signer list: %w", err)
	}

	// Create signer entries
	for _, entry := range input.SignerEntries {
		signerEntry, err := domainXRP.NewXRPSignerEntry(signerListID, entry.Account, entry.Weight)
		if err != nil {
			return watchusecase.SetSignerListOutput{},
				fmt.Errorf("failed to create signer entry entity: %w", err)
		}

		if _, err = u.signerEntryRepo.Insert(ctx, signerEntry); err != nil {
			return watchusecase.SetSignerListOutput{},
				fmt.Errorf("failed to insert signer entry for %s: %w", entry.Account, err)
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
		return watchusecase.SetSignerListOutput{},
			fmt.Errorf("failed to write transaction file: %w", err)
	}

	return watchusecase.SetSignerListOutput{
		FileName:     generatedFileName,
		TxJSON:       txJSON,
		SignerListID: signerListID,
	}, nil
}
