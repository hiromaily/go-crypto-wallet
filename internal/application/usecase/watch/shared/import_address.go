package shared

import (
	"context"
	"fmt"
	"strings"

	watchusecase "github.com/hiromaily/go-crypto-wallet/internal/application/usecase/watch"
	domainCoin "github.com/hiromaily/go-crypto-wallet/internal/domain/coin"
	domainWallet "github.com/hiromaily/go-crypto-wallet/internal/domain/wallet"
	sqlc "github.com/hiromaily/go-crypto-wallet/internal/infrastructure/database/mysql/sqlcgen"
	"github.com/hiromaily/go-crypto-wallet/internal/infrastructure/repository/watch"
	"github.com/hiromaily/go-crypto-wallet/internal/infrastructure/storage/file"
	"github.com/hiromaily/go-crypto-wallet/internal/infrastructure/storage/file/address"
)

type importAddressUseCase struct {
	addrRepo     watch.AddressRepositorier
	addrFileRepo file.AddressFileRepositorier
	coinTypeCode domainCoin.CoinTypeCode
	addrType     address.AddrType
	wtype        domainWallet.WalletType
}

// NewImportAddressUseCase creates a new ImportAddressUseCase for watch wallet
func NewImportAddressUseCase(
	addrRepo watch.AddressRepositorier,
	addrFileRepo file.AddressFileRepositorier,
	coinTypeCode domainCoin.CoinTypeCode,
	addrType address.AddrType,
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

	pubKeyData := make([]*sqlc.Address, 0, len(pubKeys))
	for _, key := range pubKeys {
		// coin, account, ...
		inner := strings.Split(key, ",")

		var addrFmt *address.AddressFormat
		addrFmt, err = address.ConvertLine(u.coinTypeCode, inner)
		if err != nil {
			return err
		}

		pubKeyData = append(pubKeyData, &sqlc.Address{
			Coin:          sqlc.AddressCoin(u.coinTypeCode.String()),
			Account:       sqlc.AddressAccount(addrFmt.AccountType.String()),
			WalletAddress: addrFmt.P2PKHAddress,
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
