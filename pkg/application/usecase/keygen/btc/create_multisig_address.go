package btc

import (
	"context"
	"fmt"

	"github.com/hiromaily/go-crypto-wallet/pkg/application/usecase/keygen"
	btckeygensrv "github.com/hiromaily/go-crypto-wallet/pkg/wallet/service/keygen/btc"
)

type createMultisigAddressUseCase struct {
	multisig *btckeygensrv.Multisig
}

// NewCreateMultisigAddressUseCase creates a new CreateMultisigAddressUseCase
func NewCreateMultisigAddressUseCase(multisig *btckeygensrv.Multisig) keygen.CreateMultisigAddressUseCase {
	return &createMultisigAddressUseCase{
		multisig: multisig,
	}
}

func (u *createMultisigAddressUseCase) Create(ctx context.Context, input keygen.CreateMultisigAddressInput) error {
	if err := u.multisig.AddMultisigAddress(input.AccountType, input.AddressType); err != nil {
		return fmt.Errorf("failed to create multisig address: %w", err)
	}
	return nil
}
