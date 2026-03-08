package di

import (
	"context"
	"errors"

	watchusecase "github.com/hiromaily/go-crypto-wallet/internal/application/usecase/watch"
)

// notImplementedAddMultisigSignatureUseCase is a no-op implementation of
// AddMultisigSignatureUseCase used by the DI container until the full
// implementation is wired. It always returns a clear "not yet implemented"
// error instead of panicking at startup.
//
// A nil interface return is explicitly forbidden — a nil interface value
// causes a silent nil pointer panic at the first call site, which is worse
// than the current named startup panic.
type notImplementedAddMultisigSignatureUseCase struct{}

// NewNotImplementedAddMultisigSignatureUseCase returns a no-op
// AddMultisigSignatureUseCase that satisfies the interface but returns
// an error on every call to Execute.
func NewNotImplementedAddMultisigSignatureUseCase() watchusecase.AddMultisigSignatureUseCase {
	return &notImplementedAddMultisigSignatureUseCase{}
}

func (*notImplementedAddMultisigSignatureUseCase) Execute(
	_ context.Context,
	_ watchusecase.AddMultisigSignatureInput,
) (watchusecase.AddMultisigSignatureOutput, error) {
	return watchusecase.AddMultisigSignatureOutput{},
		errors.New("AddMultisigSignature is not yet implemented")
}
