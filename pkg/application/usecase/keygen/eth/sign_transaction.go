package eth

import (
	"context"
	"fmt"

	keygenusecase "github.com/hiromaily/go-crypto-wallet/pkg/application/usecase/keygen"
	signsrv "github.com/hiromaily/go-crypto-wallet/pkg/wallet/service/sign"
)

type signTransactionUseCase struct {
	signer signsrv.Signer
}

// NewSignTransactionUseCase creates a new SignTransactionUseCase for ETH keygen
func NewSignTransactionUseCase(signer signsrv.Signer) keygenusecase.SignTransactionUseCase {
	return &signTransactionUseCase{
		signer: signer,
	}
}

func (u *signTransactionUseCase) Sign(
	ctx context.Context, input keygenusecase.SignTransactionInput,
) (keygenusecase.SignTransactionOutput, error) {
	_, isSigned, generatedFileName, err := u.signer.SignTx(input.FilePath)
	if err != nil {
		return keygenusecase.SignTransactionOutput{}, fmt.Errorf("failed to sign transaction: %w", err)
	}

	return keygenusecase.SignTransactionOutput{
		FilePath:      generatedFileName,
		IsDone:        isSigned,
		SignedCount:   1, // ETH signs one transaction at a time
		UnsignedCount: 0,
	}, nil
}
