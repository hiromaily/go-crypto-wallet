package xrp

import (
	"context"
	"errors"
	"fmt"

	apixrp "github.com/hiromaily/go-crypto-wallet/internal/application/ports/api/xrp"
	repocold "github.com/hiromaily/go-crypto-wallet/internal/application/ports/repository/cold"
	repowatch "github.com/hiromaily/go-crypto-wallet/internal/application/ports/repository/watch"
	watchusecase "github.com/hiromaily/go-crypto-wallet/internal/application/usecase/watch"
	domainXRP "github.com/hiromaily/go-crypto-wallet/internal/domain/chains/xrp"
)

// xrpAddMultisigSigClient defines the interface for XRP operations needed by addMultisigSignatureUseCase.
// This follows the Interface Segregation Principle - depend only on what you need.
type xrpAddMultisigSigClient interface {
	apixrp.TransactionCombiner
}

type addMultisigSignatureUseCase struct {
	xrper               xrpAddMultisigSigClient
	signerListRepo      repocold.XRPSignerListRepositorier
	signerEntryRepo     repocold.XRPSignerEntryRepositorier
	pendingMultisigRepo repowatch.XRPPendingMultisigRepositorier
	signatureRepo       repowatch.XRPMultisigSignatureRepositorier
}

// NewAddMultisigSignatureUseCase creates a new AddMultisigSignatureUseCase.
// The xrper parameter accepts any type that implements xrpAddMultisigSigClient (TransactionCombiner).
// Typically, apixrp.XRPer is passed which implements all required methods.
func NewAddMultisigSignatureUseCase(
	xrper xrpAddMultisigSigClient,
	signerListRepo repocold.XRPSignerListRepositorier,
	signerEntryRepo repocold.XRPSignerEntryRepositorier,
	pendingMultisigRepo repowatch.XRPPendingMultisigRepositorier,
	signatureRepo repowatch.XRPMultisigSignatureRepositorier,
) watchusecase.AddMultisigSignatureUseCase {
	return &addMultisigSignatureUseCase{
		xrper:               xrper,
		signerListRepo:      signerListRepo,
		signerEntryRepo:     signerEntryRepo,
		pendingMultisigRepo: pendingMultisigRepo,
		signatureRepo:       signatureRepo,
	}
}

func (u *addMultisigSignatureUseCase) Execute(
	ctx context.Context,
	input watchusecase.AddMultisigSignatureInput,
) (watchusecase.AddMultisigSignatureOutput, error) {
	if err := u.validateInput(input); err != nil {
		return watchusecase.AddMultisigSignatureOutput{}, err
	}

	pendingTx, err := u.getPendingTransaction(ctx, input.TxUUID)
	if err != nil {
		return watchusecase.AddMultisigSignatureOutput{}, err
	}

	if err = u.checkExistingSignature(ctx, pendingTx.ID, input.SignerAccount); err != nil {
		return watchusecase.AddMultisigSignatureOutput{}, err
	}

	signerEntry, err := u.getSignerEntry(ctx, pendingTx.AccountID, input.SignerAccount)
	if err != nil {
		return watchusecase.AddMultisigSignatureOutput{}, err
	}

	if err = u.createAndInsertSignature(ctx, pendingTx.ID, input, signerEntry.SignerWeight); err != nil {
		return watchusecase.AddMultisigSignatureOutput{}, err
	}

	newWeight := pendingTx.CurrentWeight + signerEntry.SignerWeight
	quorumMet := pendingTx.AddWeight(signerEntry.SignerWeight)

	if err = u.pendingMultisigRepo.UpdateWeight(ctx, pendingTx.ID, newWeight); err != nil {
		return watchusecase.AddMultisigSignatureOutput{}, fmt.Errorf("failed to update weight: %w", err)
	}

	combinedTxBlob, err := u.handleQuorumIfMet(ctx, pendingTx.ID, quorumMet)
	if err != nil {
		return watchusecase.AddMultisigSignatureOutput{}, err
	}

	return watchusecase.AddMultisigSignatureOutput{
		CurrentWeight:  newWeight,
		RequiredQuorum: pendingTx.RequiredQuorum,
		IsReady:        quorumMet,
		CombinedTxBlob: combinedTxBlob,
	}, nil
}

func (*addMultisigSignatureUseCase) validateInput(
	input watchusecase.AddMultisigSignatureInput,
) error {
	if input.TxUUID == "" {
		return errors.New("transaction UUID is required")
	}
	if input.SignerAccount == "" {
		return errors.New("signer account is required")
	}
	if input.SignedTxBlob == "" {
		return errors.New("signed transaction blob is required")
	}
	return nil
}

func (u *addMultisigSignatureUseCase) getPendingTransaction(
	ctx context.Context,
	txUUID string,
) (*domainXRP.XRPPendingMultisig, error) {
	pendingTx, err := u.pendingMultisigRepo.GetByUUID(ctx, txUUID)
	if err != nil {
		return nil, fmt.Errorf("failed to get pending transaction: %w", err)
	}
	if pendingTx == nil {
		return nil, fmt.Errorf("pending transaction not found: %s", txUUID)
	}
	if !pendingTx.CanAddSignature() {
		return nil, fmt.Errorf("transaction cannot accept signatures (status: %s)", pendingTx.Status)
	}
	return pendingTx, nil
}

func (u *addMultisigSignatureUseCase) checkExistingSignature(
	ctx context.Context,
	pendingID int64,
	signerAccount string,
) error {
	existingSig, err := u.signatureRepo.GetByPendingAndSigner(ctx, pendingID, signerAccount)
	if err != nil {
		return fmt.Errorf("failed to check existing signature: %w", err)
	}
	if existingSig != nil {
		return fmt.Errorf("signer %s has already signed this transaction", signerAccount)
	}
	return nil
}

func (u *addMultisigSignatureUseCase) getSignerEntry(
	ctx context.Context,
	accountID, signerAccount string,
) (*domainXRP.XRPSignerEntry, error) {
	signerList, err := u.signerListRepo.GetByAccountID(ctx, accountID)
	if err != nil {
		return nil, fmt.Errorf("failed to get signer list: %w", err)
	}
	if signerList == nil {
		return nil, fmt.Errorf("no active signer list found for account %s", accountID)
	}

	signerEntry, err := u.signerEntryRepo.GetByListAndAccount(ctx, signerList.ID, signerAccount)
	if err != nil {
		return nil, fmt.Errorf("failed to get signer entry: %w", err)
	}
	if signerEntry == nil {
		return nil, fmt.Errorf("signer %s is not authorized for account %s", signerAccount, accountID)
	}
	return signerEntry, nil
}

func (u *addMultisigSignatureUseCase) createAndInsertSignature(
	ctx context.Context,
	pendingID int64,
	input watchusecase.AddMultisigSignatureInput,
	weight uint32,
) error {
	signature, err := domainXRP.NewXRPMultisigSignature(
		pendingID,
		input.SignerAccount,
		input.SignedTxBlob,
		weight,
	)
	if err != nil {
		return fmt.Errorf("failed to create signature entity: %w", err)
	}

	if _, err = u.signatureRepo.Insert(ctx, signature); err != nil {
		return fmt.Errorf("failed to insert signature: %w", err)
	}
	return nil
}

func (u *addMultisigSignatureUseCase) handleQuorumIfMet(
	ctx context.Context,
	pendingID int64,
	quorumMet bool,
) (string, error) {
	if !quorumMet {
		return "", nil
	}

	signatures, err := u.signatureRepo.GetByPendingID(ctx, pendingID)
	if err != nil {
		return "", fmt.Errorf("failed to get signatures: %w", err)
	}

	signedTxBlobs := make([]string, len(signatures))
	for i, sig := range signatures {
		signedTxBlobs[i] = sig.SignedTxBlob
	}

	_, combinedTxBlob, err := u.xrper.CombineTransaction(ctx, signedTxBlobs)
	if err != nil {
		return "", fmt.Errorf("failed to combine signatures: %w", err)
	}

	if err = u.pendingMultisigRepo.SetReady(ctx, pendingID, combinedTxBlob); err != nil {
		return "", fmt.Errorf("failed to set transaction ready: %w", err)
	}

	return combinedTxBlob, nil
}
