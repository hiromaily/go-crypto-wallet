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

// notImplementedCreateETHMultisigTransactionUseCase is a placeholder until Task 6.1 wires SafeClient.
type notImplementedCreateETHMultisigTransactionUseCase struct{}

// NewNotImplementedCreateETHMultisigTransactionUseCase returns a stub that errors on every call.
func NewNotImplementedCreateETHMultisigTransactionUseCase() watchusecase.CreateETHMultisigTransactionUseCase {
	return &notImplementedCreateETHMultisigTransactionUseCase{}
}

func (*notImplementedCreateETHMultisigTransactionUseCase) Execute(
	_ context.Context,
	_ watchusecase.CreateETHMultisigTransactionInput,
) (watchusecase.CreateETHMultisigTransactionOutput, error) {
	return watchusecase.CreateETHMultisigTransactionOutput{},
		errors.New("CreateETHMultisigTransaction is not yet implemented (wire SafeClient in Task 6.1)")
}

// notImplementedSendETHMultisigTransactionUseCase is a placeholder until Task 6.1 wires SafeClient.
type notImplementedSendETHMultisigTransactionUseCase struct{}

// NewNotImplementedSendETHMultisigTransactionUseCase returns a stub that errors on every call.
func NewNotImplementedSendETHMultisigTransactionUseCase() watchusecase.SendETHMultisigTransactionUseCase {
	return &notImplementedSendETHMultisigTransactionUseCase{}
}

func (*notImplementedSendETHMultisigTransactionUseCase) Execute(
	_ context.Context,
	_ watchusecase.SendETHMultisigTransactionInput,
) (watchusecase.SendETHMultisigTransactionOutput, error) {
	return watchusecase.SendETHMultisigTransactionOutput{},
		errors.New("SendETHMultisigTransaction is not yet implemented (wire SafeClient in Task 6.1)")
}

// notImplementedETHSafeInfoUseCase is a placeholder until Task 6.1 wires SafeClient.
type notImplementedETHSafeInfoUseCase struct{}

// NewNotImplementedETHSafeInfoUseCase returns a stub that errors on every call.
func NewNotImplementedETHSafeInfoUseCase() watchusecase.ETHSafeInfoUseCase {
	return &notImplementedETHSafeInfoUseCase{}
}

func (*notImplementedETHSafeInfoUseCase) Execute(
	_ context.Context,
	_ watchusecase.ETHSafeInfoInput,
) (watchusecase.ETHSafeInfoOutput, error) {
	return watchusecase.ETHSafeInfoOutput{},
		errors.New("ETHSafeInfo is not yet implemented (wire SafeClient in Task 6.1)")
}
