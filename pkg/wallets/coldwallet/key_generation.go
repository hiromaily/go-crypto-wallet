package coldwallet

//cold wallet (keygen, sing)

import (
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/hiromaily/go-bitcoin/pkg/account"
	"github.com/hiromaily/go-bitcoin/pkg/model/rdb/coldrepo"
	ctype "github.com/hiromaily/go-bitcoin/pkg/wallets/api/types"
	"github.com/hiromaily/go-bitcoin/pkg/wallets/types"
	"github.com/hiromaily/go-bitcoin/pkg/wallets/wkey"
)

//1. generate seed and store it in database
//2. generate multisig key and store it in database
//    auth device should be separated into each auth accounts (auth1, auth2 ...)
//3. generate pubkey for client and store it in database
//4.Receipt Keyの生成 + Multisig対応 + DBに登録 (1日1Key消費するイメージ)
//5.Payment Keyの生成+ Multisig + DBに登録 (1日1Key消費するイメージ)

// GenerateSeed generate seed and store it in database
func (w *ColdWallet) GenerateSeed() ([]byte, error) {

	// retrieve seed from database
	bSeed, err := w.retrieveSeed()
	if err == nil {
		return bSeed, nil
	}

	// generate seed
	bSeed, err = wkey.GenerateSeed()
	if err != nil {
		return nil, errors.Wrap(err, "fail to call key.GenerateSeed()")
	}
	strSeed := wkey.SeedToString(bSeed)

	// insert seed in database
	_, err = w.storager.InsertSeed(strSeed, nil, true)
	if err != nil {
		return nil, errors.Wrap(err, "fail to call storager.InsertSeed()")
	}

	return bSeed, nil
}

// store given seed from command line args
//  development use
func (w *ColdWallet) StoreSeed(strSeed string) ([]byte, error) {
	bSeed, err := wkey.SeedToByte(strSeed)
	if err != nil {
		return nil, errors.Wrap(err, "fail to call key.SeedToByte() ")
	}

	// insert seed in database
	_, err = w.storager.InsertSeed(strSeed, nil, true)
	if err != nil {
		return nil, errors.Wrap(err, "fail to call storager.InsertSeed()")
	}

	return bSeed, nil
}

// retrieve seed from database
func (w *ColdWallet) retrieveSeed() ([]byte, error) {
	// get seed from database, seed is expected only one record
	seed, err := w.storager.GetSeedOne()
	if err == nil && seed.Seed != "" {
		w.logger.Info("seed have already been generated")
		return wkey.SeedToByte(seed.Seed)
	}
	if err != nil {
		return nil, errors.Wrap(err, "fail to call storager.GetSeedOne()")
	}
	// in this case, though err didn't happen, but seed is blank
	return nil, errors.New("somehow seed retrieved from database is blank ")
}

// GenerateAccountKey AccountType属性のアカウントKeyを生成する
// TODO:AccountTypeAuthorizationのときは、レコードがある場合は追加できないようにしたほうがいい？？
func (w *ColdWallet) GenerateAccountKey(accountType account.AccountType, coinType ctype.CoinType, seed []byte, count uint32) ([]wkey.WalletKey, error) {
	if w.wtype == types.WalletTypeWatchOnly {
		return nil, errors.New("it's available on Coldwallet1, Coldwallet2")
	}

	//現在のindexを取得
	idx, err := w.storager.GetMaxIndexOnAccountKeyTable(accountType)
	if err != nil {
		idx = 0
	} else {
		idx++
	}
	w.logger.Info(
		"call storager.GetMaxIndexOnAccountKeyTable() current index",
		zap.Int64("idx", idx))

	return w.generateAccountKey(accountType, coinType, seed, uint32(idx), count)
}

// generateKey AccountType属性のアカウントKeyを生成する
func (w *ColdWallet) generateAccountKey(accountType account.AccountType, coinType ctype.CoinType, seed []byte, idxFrom, count uint32) ([]wkey.WalletKey, error) {
	// HDウォレットのkeyを生成する
	walletKeys, err := w.generateAccountKeyData(accountType, coinType, seed, idxFrom, count)
	if err != nil {
		return nil, errors.Errorf("key.generateAccountKeyData(AccountTypeClient) error: %s", err)
	}

	// Account
	//var account string
	//if accountType != ctype.AccountTypeClient {
	//	account = string(accountType)
	//}
	account := string(accountType)

	// DBにClientAccountのKey情報を登録
	accountKeyClients := make([]coldrepo.AccountKeyTable, len(walletKeys))
	for idx, key := range walletKeys {
		accountKeyClients[idx] = coldrepo.AccountKeyTable{
			WalletAddress:         key.Address,
			P2shSegwitAddress:     key.P2shSegwit,
			FullPublicKey:         key.FullPubKey,
			WalletMultisigAddress: "",
			RedeemScript:          key.RedeemScript,
			WalletImportFormat:    key.WIF,
			Account:               account,
			Idx:                   idxFrom,
		}
		idxFrom++
	}
	err = w.storager.InsertAccountKeyTable(accountType, accountKeyClients, nil, true)
	if err != nil {
		return nil, errors.Errorf("DB.InsertAccountKeyTable() error: %s", err)
	}

	return walletKeys, err
}

// generateKeyData AccountType属性のアカウントKeyを生成する
func (w *ColdWallet) generateAccountKeyData(accountType account.AccountType, coinType ctype.CoinType, seed []byte, idxFrom, count uint32) ([]wkey.WalletKey, error) {
	// Keyオブジェクト
	keyData := wkey.NewKey(coinType, w.btc.GetChainConf(), w.logger)

	// key生成
	priv, _, err := keyData.CreateAccount(seed, accountType)
	if err != nil {
		return nil, errors.Errorf("key.CreateAccount() error: %s", err)
	}

	walletKeys, err := keyData.CreateKeysWithIndex(priv, idxFrom, count)
	if err != nil {
		return nil, errors.Errorf("key.CreateKeysWithIndex() error: %s", err)
	}

	return walletKeys, nil
}
