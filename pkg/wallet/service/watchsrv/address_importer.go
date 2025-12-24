package watchsrv

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/hiromaily/go-crypto-wallet/pkg/address"
	domainCoin "github.com/hiromaily/go-crypto-wallet/pkg/domain/coin"
	domainWallet "github.com/hiromaily/go-crypto-wallet/pkg/domain/wallet"
	models "github.com/hiromaily/go-crypto-wallet/pkg/models/rdb"
	"github.com/hiromaily/go-crypto-wallet/pkg/repository/watchrepo"
)

// AddressImporter is AddressImporter interface
type AddressImporter interface {
	ImportAddress(fileName string) error
}

// AddressImport type
type AddressImport struct {
	dbConn       *sql.DB
	addrRepo     watchrepo.AddressRepositorier
	addrFileRepo address.FileRepositorier
	coinTypeCode domainCoin.CoinTypeCode
	addrType     address.AddrType
	wtype        domainWallet.WalletType
}

// NewAddressImport returns AddressImport object
func NewAddressImport(
	dbConn *sql.DB,
	addrRepo watchrepo.AddressRepositorier,
	addrFileRepo address.FileRepositorier,
	coinTypeCode domainCoin.CoinTypeCode,
	addrType address.AddrType,
	wtype domainWallet.WalletType,
) *AddressImport {
	return &AddressImport{
		dbConn:       dbConn,
		addrRepo:     addrRepo,
		addrFileRepo: addrFileRepo,
		coinTypeCode: coinTypeCode,
		addrType:     addrType,
		wtype:        wtype,
	}
}

// ImportAddress import PubKey from csv filecsv into database,
//   - if account is client, which doesn't have account ??
func (a *AddressImport) ImportAddress(fileName string) error {
	// read file for public key
	pubKeys, err := a.addrFileRepo.ImportAddress(fileName)
	if err != nil {
		return fmt.Errorf("fail to call key.ImportPubKey(): %w", err)
	}

	pubKeyData := make([]*models.Address, 0, len(pubKeys))
	for _, key := range pubKeys {
		// coin, account, ...
		inner := strings.Split(key, ",")

		var addrFmt *address.AddressFormat
		addrFmt, err = address.ConvertLine(a.coinTypeCode, inner)
		if err != nil {
			return err
		}

		pubKeyData = append(pubKeyData, &models.Address{
			Coin:          a.coinTypeCode.String(),
			Account:       addrFmt.AccountType.String(),
			WalletAddress: addrFmt.P2PKHAddress,
		})
	}

	// insert imported pubKey
	err = a.addrRepo.InsertBulk(pubKeyData)
	if err != nil {
		return fmt.Errorf("fail to call repo.Pubkey().InsertBulk(): %w", err)
		// TODO:What if this inserting is failed, how it can be recovered to keep consistancy
		// pubkey is added in wallet, but database doesn't have records
		// try to run this func again
	}

	return nil
}
