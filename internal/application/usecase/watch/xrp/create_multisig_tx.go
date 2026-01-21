package xrp

import (
	"context"
	"errors"
	"fmt"
	"time"

	dtoxrp "github.com/hiromaily/go-crypto-wallet/internal/application/dto/xrp"
	apixrp "github.com/hiromaily/go-crypto-wallet/internal/application/ports/api/xrp"
	repocold "github.com/hiromaily/go-crypto-wallet/internal/application/ports/repository/cold"
	repowatch "github.com/hiromaily/go-crypto-wallet/internal/application/ports/repository/watch"
	watchusecase "github.com/hiromaily/go-crypto-wallet/internal/application/usecase/watch"
	domainXrp "github.com/hiromaily/go-crypto-wallet/internal/domain/xrp"
	"github.com/hiromaily/go-crypto-wallet/pkg/uuid"
)

// DefaultMultisigExpiration is the default expiration time for pending multisig transactions
const DefaultMultisigExpiration = 24 * time.Hour

type createMultisigTxUseCase struct {
	rippler             apixrp.XRPer
	uuidHandler         uuid.UUIDHandler
	signerListRepo      repocold.XRPSignerListRepositorier
	pendingMultisigRepo repowatch.XRPPendingMultisigRepositorier
}

// NewCreateMultisigTxUseCase creates a new CreateMultisigTxUseCase
func NewCreateMultisigTxUseCase(
	rippler apixrp.XRPer,
	uuidHandler uuid.UUIDHandler,
	signerListRepo repocold.XRPSignerListRepositorier,
	pendingMultisigRepo repowatch.XRPPendingMultisigRepositorier,
) watchusecase.CreateMultisigTxUseCase {
	return &createMultisigTxUseCase{
		rippler:             rippler,
		uuidHandler:         uuidHandler,
		signerListRepo:      signerListRepo,
		pendingMultisigRepo: pendingMultisigRepo,
	}
}

func (u *createMultisigTxUseCase) Execute(
	ctx context.Context,
	input watchusecase.CreateMultisigTxInput,
) (watchusecase.CreateMultisigTxOutput, error) {
	// Validate input
	if input.AccountAddress == "" {
		return watchusecase.CreateMultisigTxOutput{}, errors.New("account address is required")
	}
	if input.ReceiverAddress == "" {
		return watchusecase.CreateMultisigTxOutput{}, errors.New("receiver address is required")
	}
	if input.Amount <= 0 {
		return watchusecase.CreateMultisigTxOutput{}, errors.New("amount must be positive")
	}

	// Get the active signer list for this account to determine required quorum
	signerList, err := u.signerListRepo.GetByAccountID(ctx, input.AccountAddress)
	if err != nil {
		return watchusecase.CreateMultisigTxOutput{},
			fmt.Errorf("failed to get signer list: %w", err)
	}
	if signerList == nil {
		return watchusecase.CreateMultisigTxOutput{},
			fmt.Errorf("no active signer list found for account %s", input.AccountAddress)
	}

	// Prepare the transaction (Payment for now)
	instructions := &dtoxrp.Instructions{
		MaxLedgerVersionOffset: domainXrp.MaxLedgerVersionOffset,
	}

	_, txJSON, err := u.rippler.CreateRawTransaction(
		ctx, input.AccountAddress, input.ReceiverAddress, input.Amount, instructions)
	if err != nil {
		return watchusecase.CreateMultisigTxOutput{},
			fmt.Errorf("failed to create raw transaction: %w", err)
	}

	// Generate UUID for tracking
	uid, err := u.uuidHandler.GenerateV7()
	if err != nil {
		return watchusecase.CreateMultisigTxOutput{},
			fmt.Errorf("failed to generate UUID: %w", err)
	}

	// Set expiration time
	expiresAt := time.Now().Add(DefaultMultisigExpiration)

	// Determine transaction type
	txType := input.TxType
	if txType == "" {
		txType = "Payment"
	}

	// Create pending multisig record
	pendingMultisig, err := domainXrp.NewXRPPendingMultisig(
		uid.String(),
		input.AccountAddress,
		txJSON,
		txType,
		signerList.SignerQuorum,
		&expiresAt,
	)
	if err != nil {
		return watchusecase.CreateMultisigTxOutput{},
			fmt.Errorf("failed to create pending multisig entity: %w", err)
	}

	pendingID, err := u.pendingMultisigRepo.Insert(ctx, pendingMultisig)
	if err != nil {
		return watchusecase.CreateMultisigTxOutput{},
			fmt.Errorf("failed to insert pending multisig: %w", err)
	}

	return watchusecase.CreateMultisigTxOutput{
		TxUUID:         uid.String(),
		PendingID:      pendingID,
		UnsignedTxJSON: txJSON,
		RequiredQuorum: signerList.SignerQuorum,
	}, nil
}
