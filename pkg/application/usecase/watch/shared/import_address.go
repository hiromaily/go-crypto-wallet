package shared

import (
	"context"
	"fmt"
	"strings"

	"github.com/hiromaily/go-crypto-wallet/pkg/address"
	watchusecase "github.com/hiromaily/go-crypto-wallet/pkg/application/usecase/watch"
	domainCoin "github.com/hiromaily/go-crypto-wallet/pkg/domain/coin"
	domainWallet "github.com/hiromaily/go-crypto-wallet/pkg/domain/wallet"
	"github.com/hiromaily/go-crypto-wallet/pkg/infrastructure/repository/watch"
	"github.com/hiromaily/go-crypto-wallet/pkg/infrastructure/storage/file"
	models "github.com/hiromaily/go-crypto-wallet/pkg/models/rdb"
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

	pubKeyData := make([]*models.Address, 0, len(pubKeys))
	for _, key := range pubKeys {
		// coin, account, ...
		inner := strings.Split(key, ",")

		var addrFmt *address.AddressFormat
		addrFmt, err = address.ConvertLine(u.coinTypeCode, inner)
		if err != nil {
			return err
		}

		pubKeyData = append(pubKeyData, &models.Address{
			Coin:          u.coinTypeCode.String(),
			Account:       addrFmt.AccountType.String(),
			WalletAddress: addrFmt.P2PKHAddress,
		})
	}

	// insert imported pubKey
	err = u.addrRepo.InsertBulk(pubKeyData)
	if err != nil {
		return fmt.Errorf("fail to call addrRepo.InsertBulk(): %w", err)
		// TODO:What if this inserting is failed, how it can be recovered to keep consistancy
		// pubkey is added in wallet, but database doesn't have records
		// try to run this func again
	}

	return nil
}
