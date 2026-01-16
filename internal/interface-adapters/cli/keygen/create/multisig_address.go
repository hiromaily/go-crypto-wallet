package create

import (
	"context"
	"errors"
	"fmt"

	keygenusecase "github.com/hiromaily/go-crypto-wallet/internal/application/usecase/keygen"
	"github.com/hiromaily/go-crypto-wallet/internal/di"
	domainAccount "github.com/hiromaily/go-crypto-wallet/internal/domain/account"
	apibtcimpl "github.com/hiromaily/go-crypto-wallet/internal/infrastructure/api/btc/btc"
)

// runMultisigWithFlags is the actual implementation that accepts parsed flags
func runMultisigWithFlags(container di.Container, acnt, multisigType string) error {
	// Validate account type
	if !domainAccount.ValidateAccountType(acnt) {
		return errors.New("account option [--account] is invalid")
	}

	// Validate multisig type
	switch multisigType {
	case "traditional":
		return runTraditionalMultisig(container, acnt)
	case "musig2":
		return runMuSig2Address(container, acnt)
	default:
		return fmt.Errorf("invalid multisig-type: %s (must be 'traditional' or 'musig2')", multisigType)
	}
}

// runTraditionalMultisig creates traditional P2SH/P2WSH multisig addresses
func runTraditionalMultisig(container di.Container, acnt string) error {
	fmt.Println("create traditional multisig address")

	// create multisig address
	useCase := container.NewKeygenCreateMultisigAddressUseCase()
	err := useCase.Create(context.Background(), keygenusecase.CreateMultisigAddressInput{
		AccountType: domainAccount.AccountType(acnt),
		AddressType: apibtcimpl.ToAddressType(container.AddressType()),
	})
	if err != nil {
		return fmt.Errorf("fail to create traditional multisig address: %w", err)
	}

	fmt.Println("✓ Traditional multisig addresses created successfully")
	return nil
}

// runMuSig2Address creates MuSig2 Taproot addresses
func runMuSig2Address(container di.Container, acnt string) error {
	fmt.Println("create MuSig2 Taproot address")

	// create MuSig2 address
	useCase := container.NewKeygenCreateMuSig2AddressUseCase()
	err := useCase.Create(context.Background(), keygenusecase.CreateMuSig2AddressInput{
		AccountType: domainAccount.AccountType(acnt),
	})
	if err != nil {
		return fmt.Errorf("fail to create MuSig2 address: %w", err)
	}

	fmt.Println("✓ MuSig2 Taproot addresses created successfully")
	return nil
}
