package btc

import (
	"context"
	"fmt"

	"github.com/hiromaily/go-crypto-wallet/pkg/account"
	"github.com/hiromaily/go-crypto-wallet/pkg/address"
	keygenusecase "github.com/hiromaily/go-crypto-wallet/pkg/application/usecase/keygen"
	"github.com/hiromaily/go-crypto-wallet/pkg/infrastructure/api/bitcoin"
	"github.com/hiromaily/go-crypto-wallet/pkg/infrastructure/api/bitcoin/btc"
	"github.com/hiromaily/go-crypto-wallet/pkg/infrastructure/repository/cold"
	"github.com/hiromaily/go-crypto-wallet/pkg/logger"
)

type createMultisigAddressUseCase struct {
	btc                bitcoin.Bitcoiner
	authFullPubKeyRepo cold.AuthFullPubkeyRepositorier
	accountKeyRepo     cold.AccountKeyRepositorier
	multisigAccount    account.MultisigAccounter
}

// NewCreateMultisigAddressUseCase creates a new CreateMultisigAddressUseCase
func NewCreateMultisigAddressUseCase(
	btc bitcoin.Bitcoiner,
	authFullPubKeyRepo cold.AuthFullPubkeyRepositorier,
	accountKeyRepo cold.AccountKeyRepositorier,
	multisigAccount account.MultisigAccounter,
) keygenusecase.CreateMultisigAddressUseCase {
	return &createMultisigAddressUseCase{
		btc:                btc,
		authFullPubKeyRepo: authFullPubKeyRepo,
		accountKeyRepo:     accountKeyRepo,
		multisigAccount:    multisigAccount,
	}
}

func (u *createMultisigAddressUseCase) Create(
	ctx context.Context,
	input keygenusecase.CreateMultisigAddressInput,
) error {
	logger.Debug("addmultisigaddress",
		"account_type", input.AccountType.String(),
	)

	// Validate accountType
	if !u.multisigAccount.IsMultisigAccount(input.AccountType) {
		logger.Info("only multisig account is allowed")
		return nil
	}

	var requiredSig int
	var authFullPubKeys []string

	// Get fullPubKey from auth_fullpubkey table
	// AccountTypeDeposit: { //2:5+1
	//	2: {AuthType1, AuthType2, AuthType3, AuthType4, AuthType5},
	// },
	for sigCount, authTypes := range u.multisigAccount.MultiAccounts()[input.AccountType] {
		requiredSig = sigCount
		for _, authType := range authTypes {
			// Get record from auth_fullpubkey table
			fullPubKeyItem, err := u.authFullPubKeyRepo.GetOne(authType)
			if err != nil {
				return fmt.Errorf("fail to call authFullPubKeyRepo.GetOne() %s: %w", authType.String(), err)
			}
			authFullPubKeys = append(authFullPubKeys, fullPubKeyItem.FullPublicKey)
		}
		logger.Debug("don't repeat again")
	}

	// Get target addresses from account_key table, addr_status=AddrStatusPrivKeyImported
	accountKeyItems, err := u.accountKeyRepo.GetAllAddrStatus(input.AccountType, address.AddrStatusPrivKeyImported)
	if err != nil {
		return fmt.Errorf("fail to call accountKeyRepo.GetAllAddrStatus(%s): %w", input.AccountType.String(), err)
	}

	// Call bitcoinAPI `addmultisigaddress`
	for _, item := range accountKeyItems {
		addrs := make([]string, len(authFullPubKeys)+1)
		copy(addrs, authFullPubKeys)
		addrs[len(authFullPubKeys)] = item.FullPublicKey

		var resAddr *btc.AddMultisigAddressResult
		resAddr, err = u.btc.AddMultisigAddress(
			requiredSig,
			addrs,
			fmt.Sprintf("multi_%s", input.AccountType), // this is not important
			input.AddressType,
		)
		if err != nil {
			// [Error] -5: no full public key for address mkPmdpo59gpU7ZioGYwwoMTQJjh7MiqUvd
			logger.Error(
				"fail to call btc.AddMultisigAddress()",
				"signature_count", requiredSig,
				"full public key for accountType", item.FullPublicKey,
				"full public key for authType", authFullPubKeys,
				"error", err,
			)
			continue
		}

		// Update generated multisig address, redeemScript, addrStatus
		item.MultisigAddress = resAddr.Address
		item.RedeemScript = resAddr.RedeemScript
		item.AddrStatus = address.AddrStatusMultisigAddressGenerated.Int8()

		_, err = u.accountKeyRepo.UpdateMultisigAddr(input.AccountType, item)
		if err != nil {
			return fmt.Errorf("fail to call accountKeyRepo.UpdateMultisigAddr(%s): %w", input.AccountType.String(), err)
		}
	}

	return nil
}
