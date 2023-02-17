package watchsrv

import (
	"database/sql"
	"strings"

	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/hiromaily/go-crypto-wallet/pkg/address"
	models "github.com/hiromaily/go-crypto-wallet/pkg/models/rdb"
	"github.com/hiromaily/go-crypto-wallet/pkg/repository/watchrepo"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/coin"
)

// AddressImporter is AddressImporter interface
type AddressImporter interface {
	ImportAddress(fileName string) error
}

// AddressImport type
type AddressImport struct {
	logger       *zap.Logger
	dbConn       *sql.DB
	addrRepo     watchrepo.AddressRepositorier
	addrFileRepo address.FileRepositorier
	coinTypeCode coin.CoinTypeCode
	addrType     address.AddrType
	wtype        wallet.WalletType
}

// NewAddressImport returns AddressImport object
func NewAddressImport(
	logger *zap.Logger,
	dbConn *sql.DB,
	addrRepo watchrepo.AddressRepositorier,
	addrFileRepo address.FileRepositorier,
	coinTypeCode coin.CoinTypeCode,
	addrType address.AddrType,
	wtype wallet.WalletType,
) *AddressImport {
	return &AddressImport{
		logger:       logger,
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
		return errors.Wrap(err, "fail to call key.ImportPubKey()")
	}

	pubKeyData := make([]*models.Address, 0, len(pubKeys))
	for _, key := range pubKeys {
		// coin, account, ...
		inner := strings.Split(key, ",")

		addrFmt, err := address.ConvertLine(a.coinTypeCode, inner)
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
		return errors.Wrap(err, "fail to call repo.Pubkey().InsertBulk()")
		// TODO:What if this inserting is failed, how it can be recovered to keep consistancy
		// pubkey is added in wallet, but database doesn't have records
		// try to run this func again
	}

	return nil
}
