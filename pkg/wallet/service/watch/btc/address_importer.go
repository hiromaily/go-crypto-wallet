package btc

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/hiromaily/go-crypto-wallet/pkg/address"
	domainAccount "github.com/hiromaily/go-crypto-wallet/pkg/domain/account"
	domainCoin "github.com/hiromaily/go-crypto-wallet/pkg/domain/coin"
	domainWallet "github.com/hiromaily/go-crypto-wallet/pkg/domain/wallet"
	"github.com/hiromaily/go-crypto-wallet/pkg/infrastructure/api/bitcoin"
	"github.com/hiromaily/go-crypto-wallet/pkg/infrastructure/repository/watch"
	"github.com/hiromaily/go-crypto-wallet/pkg/infrastructure/storage/file"
	"github.com/hiromaily/go-crypto-wallet/pkg/logger"
	models "github.com/hiromaily/go-crypto-wallet/pkg/models/rdb"
)

// AddressImport type
type AddressImport struct {
	btc          bitcoin.Bitcoiner
	dbConn       *sql.DB
	addrRepo     watch.AddressRepositorier
	addrFileRepo file.AddressFileRepositorier
	coinTypeCode domainCoin.CoinTypeCode
	addrType     address.AddrType
	wtype        domainWallet.WalletType
}

// NewAddressImport returns AddressImport object
func NewAddressImport(
	btc bitcoin.Bitcoiner,
	dbConn *sql.DB,
	addrRepo watch.AddressRepositorier,
	addrFileRepo file.AddressFileRepositorier,
	coinTypeCode domainCoin.CoinTypeCode,
	addrType address.AddrType,
	wtype domainWallet.WalletType,
) *AddressImport {
	return &AddressImport{
		btc:          btc,
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
func (a *AddressImport) ImportAddress(fileName string, isRescan bool) error {
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
		addrFmt, err = address.ConvertLine(a.btc.CoinTypeCode(), inner)
		if err != nil {
			return err
		}

		var targetAddr string
		if addrFmt.AccountType == domainAccount.AccountTypeClient {
			switch a.btc.CoinTypeCode() {
			case domainCoin.BTC:
				switch a.addrType {
				case address.AddrTypeBech32:
					targetAddr = addrFmt.Bech32Address
				case address.AddrTypeLegacy, address.AddrTypeP2shSegwit,
					address.AddrTypeBCHCashAddr, address.AddrTypeETH:
					targetAddr = addrFmt.P2SHSegwitAddress // p2sh_segwit_address
				default:
					targetAddr = addrFmt.P2SHSegwitAddress // p2sh_segwit_address
				}
			case domainCoin.BCH:
				targetAddr = addrFmt.P2PKHAddress // p2pkh_address
			case domainCoin.LTC, domainCoin.ETH, domainCoin.XRP, domainCoin.ERC20, domainCoin.HYT:
				return fmt.Errorf("coinTypeCode is out of range: %s", a.btc.CoinTypeCode().String())
			default:
				return fmt.Errorf("coinTypeCode is out of range: %s", a.btc.CoinTypeCode().String())
			}
		} else {
			targetAddr = addrFmt.MultisigAddress // multisig_address
		}

		// call bitcoin API `importaddress` with account(label)
		// Note: Error would occur when using only 1 bitcoin core server under development
		// because address is already imported
		// isRescan would be `false` usually
		err = a.btc.ImportAddressWithLabel(targetAddr, addrFmt.AccountType.String(), isRescan)
		if err != nil {
			//-4: The wallet already contains the private key for this address or script
			logger.Warn(
				"fail to call btc.ImportAddressWithLabel() but continue following addresses",
				"address", targetAddr,
				"account_type", addrFmt.AccountType.String(),
				"error", err)
			continue
		}

		pubKeyData = append(pubKeyData, &models.Address{
			Coin:          a.coinTypeCode.String(),
			Account:       addrFmt.AccountType.String(),
			WalletAddress: targetAddr,
		})

		// confirm pubkey is added as watch only wallet
		a.checkImportedPubKey(targetAddr)
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

// checkImportedPubKey confirm pubkey is added as watch only wallet
func (a *AddressImport) checkImportedPubKey(addr string) {
	addrInfo, err := a.btc.GetAddressInfo(addr)
	if err != nil {
		logger.Error(
			"fail to call btc.GetAddressInfo()",
			"address", addr,
			"error", err)
		return
	}
	logger.Debug("account is found",
		"account", addrInfo.GetLabelName(),
		"address", addr)

	// `watch only wallet` is expected
	// TODO: if wallet,keygen,sign is working on only one bitcoin core server,
	// result would be `iswatchonly=false`
	if !addrInfo.Iswatchonly {
		logger.Warn("this address must be watch only wallet")
		// grok.Value(addrInfo)
	}
}
