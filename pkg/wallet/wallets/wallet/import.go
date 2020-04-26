package wallet

import (
	"strings"

	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/hiromaily/go-bitcoin/pkg/account"
	models "github.com/hiromaily/go-bitcoin/pkg/models/rdb"
)

// ImportPubKey import PubKey from csv filecsv into database,
//  - if account is client, which doesn't have account ??
func (w *Wallet) ImportPubKey(fileName string, accountType account.AccountType, isRescan bool) error {
	// read file for public key
	pubKeys, err := w.addrFileRepo.ImportPubKey(fileName)
	if err != nil {
		return errors.Wrap(err, "fail to call key.ImportPubKey()")
	}

	//var pubKeyData []walletrepo.AccountPublicKeyTable
	pubKeyData := make([]*models.Pubkey, 0, len(pubKeys))
	for _, key := range pubKeys {
		inner := strings.Split(key, ",")
		//kind of required address is different according to account
		var addr string
		if accountType == account.AccountTypeClient {
			addr = inner[1] //p2sh_segwit_address
		} else {
			addr = inner[3] //wallet_import_format
		}

		//call bitcoin API `importaddress` with account(label)
		//Note: Error would occur when using only 1 bitcoin core server under development
		// because address is already imported
		//isRescan would be `false` usually
		err := w.btc.ImportAddressWithLabel(addr, accountType.String(), isRescan)
		if err != nil {
			//-4: The wallet already contains the private key for this address or script
			w.logger.Warn(
				"fail to call btc.ImportAddressWithLabel() but continue following addresses",
				zap.String("address", addr),
				zap.String("account_type", accountType.String()),
				zap.Error(err))
			continue
		}

		pubKeyData = append(pubKeyData, &models.Pubkey{
			Coin:          w.GetBTC().CoinTypeCode().String(),
			Account:       accountType.String(),
			WalletAddress: addr,
		})

		//confirm pubkey is added as watch only wallet
		w.checkImportedPubKey(addr)
	}

	//insert imported pubKey
	err = w.repo.Pubkey().InsertBulk(pubKeyData)
	if err != nil {
		return errors.Wrap(err, "fail to call pubkeyRepo.InsertBulk()")
		//TODO:What if this inserting is failed, how it can be recovered to keep consistancy
		// pubkey is added in wallet, but database doesn't have records
		// try to run this func again
	}

	return nil
}

// checkImportedPubKey confirm pubkey is added as watch only wallet
func (w *Wallet) checkImportedPubKey(addr string) {
	addrInfo, err := w.btc.GetAddressInfo(addr)
	if err != nil {
		w.logger.Error(
			"fail to call btc.GetAddressInfo()",
			zap.String("address", addr),
			zap.Error(err))
		return
	}
	w.logger.Debug("account is found",
		zap.String("account", addrInfo.GetLabelName()),
		zap.String("address", addr))

	//`watch only wallet` is expected
	//TODO: if wallet,keygen,sign is working on only one bitcoin core server,
	// result would be `iswatchonly=false`
	if !addrInfo.Iswatchonly {
		w.logger.Warn("this address must be watch only wallet")
		//grok.Value(addrInfo)
	}
}
