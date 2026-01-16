package shared

import (
	"context"
	"fmt"
	"strings"

	appdto "github.com/hiromaily/go-crypto-wallet/internal/application/dto"
	file "github.com/hiromaily/go-crypto-wallet/internal/application/ports/file"
	repository "github.com/hiromaily/go-crypto-wallet/internal/application/ports/repository"
	watchusecase "github.com/hiromaily/go-crypto-wallet/internal/application/usecase/watch"
	domainAddress "github.com/hiromaily/go-crypto-wallet/internal/domain/address"
	domainCoin "github.com/hiromaily/go-crypto-wallet/internal/domain/coin"
	domainWallet "github.com/hiromaily/go-crypto-wallet/internal/domain/wallet"
	"github.com/hiromaily/go-crypto-wallet/internal/infrastructure/storage/file/address"
)

type importAddressUseCase struct {
	addrRepo     repository.AddressRepositorier
	addrFileRepo file.AddressFileRepositorier
	coinTypeCode domainCoin.CoinTypeCode
	addrType     domainAddress.AddrType
	wtype        domainWallet.WalletType
}

// NewImportAddressUseCase creates a new ImportAddressUseCase for watch wallet
func NewImportAddressUseCase(
	addrRepo repository.AddressRepositorier,
	addrFileRepo file.AddressFileRepositorier,
	coinTypeCode domainCoin.CoinTypeCode,
	addrType domainAddress.AddrType,
	wtype domainWallet.WalletType,
) watchusecase.ImportAddressUseCase {
	return &importAddressUseCase{
		addrRepo:     addrRepo,
		addrFileRepo: addrFileRepo,
		coinTypeCode: coinTypeCode,
		addrType:     addrType,
		wtype:        wtype,
	}
}

func (u *importAddressUseCase) Execute(ctx context.Context, input watchusecase.ImportAddressInput) error {
	// read file for public key
	pubKeys, err := u.addrFileRepo.ImportAddress(input.FileName)
	if err != nil {
		return fmt.Errorf("fail to call addrFileRepo.ImportAddress(): %w", err)
	}

	pubKeyData := make([]*domainAddress.Address, 0, len(pubKeys))
	for _, key := range pubKeys {
		// coin, account, ...
		inner := strings.Split(key, ",")

		var addrFmt *appdto.AddressFormat
		addrFmt, err = address.ConvertLine(u.coinTypeCode, inner)
		if err != nil {
			return err
		}

		pubKeyData = append(pubKeyData, &domainAddress.Address{
			CoinTypeCode:  u.coinTypeCode,
			AccountType:   addrFmt.AccountType,
			WalletAddress: addrFmt.P2PKHAddress,
			IsAllocated:   false,
		})
	}

	// insert imported pubKey
	err = u.addrRepo.InsertBulk(ctx, pubKeyData)
	if err != nil {
		return fmt.Errorf("fail to call addrRepo.InsertBulk(): %w", err)
		// TODO:What if this inserting is failed, how it can be recovered to keep consistancy
		// pubkey is added in wallet, but database doesn't have records
		// try to run this func again
	}

	return nil
}
