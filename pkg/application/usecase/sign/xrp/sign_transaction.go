package xrp

import (
	"context"
	"fmt"

	"github.com/hiromaily/go-crypto-wallet/pkg/application/usecase/sign"
	xrpsignsrv "github.com/hiromaily/go-crypto-wallet/pkg/wallet/service/sign/xrp"
)

type signTransactionUseCase struct {
	signer *xrpsignsrv.Sign
}

// NewSignTransactionUseCase creates a new SignTransactionUseCase
func NewSignTransactionUseCase(signer *xrpsignsrv.Sign) sign.SignTransactionUseCase {
	return &signTransactionUseCase{
		signer: signer,
	}
}

func (u *signTransactionUseCase) Sign(
	ctx context.Context,
	input sign.SignTransactionInput,
) (sign.SignTransactionOutput, error) {
	signedHex, isComplete, nextFilePath, err := u.signer.SignTx(input.FilePath)
	if err != nil {
		return sign.SignTransactionOutput{}, fmt.Errorf("failed to sign transaction: %w", err)
	}

	return sign.SignTransactionOutput{
		SignedHex:    signedHex,
		IsComplete:   isComplete,
		NextFilePath: nextFilePath,
	}, nil
}
